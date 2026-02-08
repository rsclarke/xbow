package xbow

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/rsclarke/xbow/internal/api"
)

func TestWebhookFromGetResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	events := makeGetEventsResponse([]string{"asset.changed", "assessment.changed"})

	resp := &api.GetAPIV1WebhooksWebhookIDResponse{
		ID:         "webhook-123",
		APIVersion: api.GetAPIV1WebhooksWebhookIDResponseAPIVersionNext,
		TargetURL:  "https://example.com/webhook",
		Events:     events,
		CreatedAt:  now,
		UpdatedAt:  now.Add(time.Hour),
	}

	got := webhookFromGetResponse(resp)

	if got.ID != "webhook-123" {
		t.Errorf("ID = %q, want 'webhook-123'", got.ID)
	}
	if got.APIVersion != WebhookAPIVersionNext {
		t.Errorf("APIVersion = %q, want %q", got.APIVersion, WebhookAPIVersionNext)
	}
	if got.TargetURL != "https://example.com/webhook" {
		t.Errorf("TargetURL = %q, want 'https://example.com/webhook'", got.TargetURL)
	}
	if len(got.Events) != 2 {
		t.Errorf("Events length = %d, want 2", len(got.Events))
	}
	if !got.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", got.CreatedAt, now)
	}
	if !got.UpdatedAt.Equal(now.Add(time.Hour)) {
		t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, now.Add(time.Hour))
	}
}

func TestWebhooksPageFromResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	t.Run("converts items and pagination info", func(t *testing.T) {
		nextCursor := "cursor-abc"
		resp := &api.GetAPIV1OrganizationsOrganizationIDWebhooksResponse{
			Items: api.GetAPIV1OrganizationsOrganizationIDWebhooks_Response_Items{
				{
					ID:         "w1",
					APIVersion: api.Next,
					TargetURL:  "https://example.com/hook1",
					Events:     makeListEventsResponse([]string{"asset.changed"}),
					CreatedAt:  now,
					UpdatedAt:  now,
				},
				{
					ID:         "w2",
					APIVersion: api.GetAPIV1OrganizationsOrganizationIDWebhooksResponseItemsAPIVersionN20260201,
					TargetURL:  "https://example.com/hook2",
					Events:     makeListEventsResponse([]string{"*"}),
					CreatedAt:  now,
					UpdatedAt:  now,
				},
			},
			NextCursor: &nextCursor,
		}

		got := webhooksPageFromResponse(resp)

		if len(got.Items) != 2 {
			t.Fatalf("got %d items, want 2", len(got.Items))
		}
		if got.Items[0].ID != "w1" {
			t.Errorf("Items[0].ID = %q, want 'w1'", got.Items[0].ID)
		}
		if got.Items[0].APIVersion != WebhookAPIVersionNext {
			t.Errorf("Items[0].APIVersion = %q, want %q", got.Items[0].APIVersion, WebhookAPIVersionNext)
		}
		if got.Items[1].APIVersion != WebhookAPIVersionN20260201 {
			t.Errorf("Items[1].APIVersion = %q, want %q", got.Items[1].APIVersion, WebhookAPIVersionN20260201)
		}
		if got.PageInfo.NextCursor == nil || *got.PageInfo.NextCursor != "cursor-abc" {
			t.Errorf("NextCursor = %v, want 'cursor-abc'", got.PageInfo.NextCursor)
		}
		if !got.PageInfo.HasMore {
			t.Error("HasMore = false, want true")
		}
	})

	t.Run("handles nil cursor", func(t *testing.T) {
		resp := &api.GetAPIV1OrganizationsOrganizationIDWebhooksResponse{
			Items:      api.GetAPIV1OrganizationsOrganizationIDWebhooks_Response_Items{},
			NextCursor: nil,
		}

		got := webhooksPageFromResponse(resp)

		if got.PageInfo.NextCursor != nil {
			t.Errorf("NextCursor = %v, want nil", got.PageInfo.NextCursor)
		}
		if got.PageInfo.HasMore {
			t.Error("HasMore = true, want false")
		}
	})

	t.Run("handles empty cursor string", func(t *testing.T) {
		empty := ""
		resp := &api.GetAPIV1OrganizationsOrganizationIDWebhooksResponse{
			Items:      api.GetAPIV1OrganizationsOrganizationIDWebhooks_Response_Items{},
			NextCursor: &empty,
		}

		got := webhooksPageFromResponse(resp)

		if got.PageInfo.HasMore {
			t.Error("HasMore = true, want false for empty cursor")
		}
	})
}

func TestDeliveriesPageFromResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	t.Run("converts delivery items", func(t *testing.T) {
		resp := &api.GetAPIV1WebhooksWebhookIDDeliveriesResponse{
			Items: api.GetAPIV1WebhooksWebhookIDDeliveries_Response_Items{
				{
					Payload: struct{}{},
					Request: api.GetAPIV1WebhooksWebhookIDDeliveries_Response_Items_Request{
						Body:    `{"type":"ping"}`,
						Headers: map[string]string{"Content-Type": "application/json"},
					},
					Response: api.GetAPIV1WebhooksWebhookIDDeliveries_Response_Items_Response{
						Body:    "",
						Headers: map[string]string{"Connection": "close"},
						Status:  204,
					},
					SentAt:  now,
					Success: true,
				},
			},
		}

		got := deliveriesPageFromResponse(resp)

		if len(got.Items) != 1 {
			t.Fatalf("got %d items, want 1", len(got.Items))
		}
		item := got.Items[0]
		if item.Request.Body != `{"type":"ping"}` {
			t.Errorf("Request.Body = %q, want JSON", item.Request.Body)
		}
		if item.Response.Status != 204 {
			t.Errorf("Response.Status = %d, want 204", item.Response.Status)
		}
		if !item.Success {
			t.Error("Success = false, want true")
		}
		if !item.SentAt.Equal(now) {
			t.Errorf("SentAt = %v, want %v", item.SentAt, now)
		}
	})

	t.Run("handles pagination", func(t *testing.T) {
		nextCursor := "cursor-xyz"
		resp := &api.GetAPIV1WebhooksWebhookIDDeliveriesResponse{
			Items:      api.GetAPIV1WebhooksWebhookIDDeliveries_Response_Items{},
			NextCursor: &nextCursor,
		}

		got := deliveriesPageFromResponse(resp)

		if !got.PageInfo.HasMore {
			t.Error("HasMore = false, want true")
		}
	})
}

func TestWebhookAPIVersionValues(t *testing.T) {
	tests := []struct {
		generated api.GetAPIV1WebhooksWebhookIDResponseAPIVersion
		want      WebhookAPIVersion
	}{
		{api.GetAPIV1WebhooksWebhookIDResponseAPIVersionN20251101, WebhookAPIVersionN20251101},
		{api.GetAPIV1WebhooksWebhookIDResponseAPIVersionN20260201, WebhookAPIVersionN20260201},
		{api.GetAPIV1WebhooksWebhookIDResponseAPIVersionNext, WebhookAPIVersionNext},
		{api.GetAPIV1WebhooksWebhookIDResponseAPIVersionUnstable, WebhookAPIVersionUnstable},
	}

	for _, tt := range tests {
		t.Run(string(tt.want), func(t *testing.T) {
			resp := &api.GetAPIV1WebhooksWebhookIDResponse{
				ID:         "w1",
				APIVersion: tt.generated,
				TargetURL:  "https://example.com",
				Events:     makeGetEventsResponse([]string{"ping"}),
			}

			got := webhookFromGetResponse(resp)
			if got.APIVersion != tt.want {
				t.Errorf("APIVersion = %q, want %q", got.APIVersion, tt.want)
			}
		})
	}
}

func TestWebhookEventTypeValues(t *testing.T) {
	eventStrings := []string{
		"ping",
		"target.changed",
		"asset.changed",
		"assessment.changed",
		"finding.changed",
		"challenge.changed",
		"*",
	}

	expectedTypes := []WebhookEventType{
		WebhookEventTypePing,
		WebhookEventTypeTargetChanged,
		WebhookEventTypeAssetChanged,
		WebhookEventTypeAssessmentChanged,
		WebhookEventTypeFindingChanged,
		WebhookEventTypeChallengeChanged,
		WebhookEventTypeAll,
	}

	events := makeGetEventsResponse(eventStrings)
	resp := &api.GetAPIV1WebhooksWebhookIDResponse{
		ID:         "w1",
		APIVersion: api.GetAPIV1WebhooksWebhookIDResponseAPIVersionNext,
		TargetURL:  "https://example.com",
		Events:     events,
	}

	got := webhookFromGetResponse(resp)

	if len(got.Events) != len(expectedTypes) {
		t.Fatalf("got %d events, want %d", len(got.Events), len(expectedTypes))
	}

	for i, want := range expectedTypes {
		if got.Events[i] != want {
			t.Errorf("Events[%d] = %q, want %q", i, got.Events[i], want)
		}
	}
}

func TestWebhookFromCreateResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	events := makeCreateEventsResponse([]string{"assessment.changed"})

	resp := &api.PostAPIV1OrganizationsOrganizationIDWebhooksResponse{
		ID:         "webhook-new",
		APIVersion: api.PostAPIV1OrganizationsOrganizationIDWebhooksResponseAPIVersionN20260201,
		TargetURL:  "https://example.com/new",
		Events:     events,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	got := webhookFromCreateResponse(resp)

	if got.ID != "webhook-new" {
		t.Errorf("ID = %q, want 'webhook-new'", got.ID)
	}
	if got.APIVersion != WebhookAPIVersionN20260201 {
		t.Errorf("APIVersion = %q, want %q", got.APIVersion, WebhookAPIVersionN20260201)
	}
	if len(got.Events) != 1 {
		t.Errorf("Events length = %d, want 1", len(got.Events))
	}
}

func TestWebhookFromPatchResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	events := makePatchEventsResponse([]string{"finding.changed", "challenge.changed"})

	resp := &api.PatchAPIV1WebhooksWebhookIDResponse{
		ID:         "webhook-updated",
		APIVersion: api.PatchAPIV1WebhooksWebhookIDResponseAPIVersionNext,
		TargetURL:  "https://example.com/updated",
		Events:     events,
		CreatedAt:  now,
		UpdatedAt:  now.Add(time.Hour),
	}

	got := webhookFromPatchResponse(resp)

	if got.ID != "webhook-updated" {
		t.Errorf("ID = %q, want 'webhook-updated'", got.ID)
	}
	if got.APIVersion != WebhookAPIVersionNext {
		t.Errorf("APIVersion = %q, want %q", got.APIVersion, WebhookAPIVersionNext)
	}
	if len(got.Events) != 2 {
		t.Errorf("Events length = %d, want 2", len(got.Events))
	}
}

// Helper functions to create test event responses

func makeGetEventsResponse(eventStrings []string) api.GetAPIV1WebhooksWebhookID_Response_Events {
	events := make(api.GetAPIV1WebhooksWebhookID_Response_Events, 0, len(eventStrings))
	for _, s := range eventStrings {
		item := api.GetAPIV1WebhooksWebhookID_Response_Events_Item{}
		item.GetAPIV1WebhooksWebhookID_Response_Events_AnyOf = &api.GetAPIV1WebhooksWebhookID_Response_Events_AnyOf{}
		_ = item.GetAPIV1WebhooksWebhookID_Response_Events_AnyOf.FromString(s)
		events = append(events, item)
	}
	return events
}

func makeListEventsResponse(eventStrings []string) api.GetAPIV1OrganizationsOrganizationIDWebhooks_Response_Items_Events {
	events := make(api.GetAPIV1OrganizationsOrganizationIDWebhooks_Response_Items_Events, 0, len(eventStrings))
	for _, s := range eventStrings {
		item := api.GetAPIV1OrganizationsOrganizationIDWebhooks_Response_Items_Events_Item{}
		item.GetAPIV1OrganizationsOrganizationIDWebhooks_Response_Items_Events_AnyOf = &api.GetAPIV1OrganizationsOrganizationIDWebhooks_Response_Items_Events_AnyOf{}
		_ = item.GetAPIV1OrganizationsOrganizationIDWebhooks_Response_Items_Events_AnyOf.FromString(s)
		events = append(events, item)
	}
	return events
}

func makeCreateEventsResponse(eventStrings []string) api.PostAPIV1OrganizationsOrganizationIDWebhooks_Response_Events {
	events := make(api.PostAPIV1OrganizationsOrganizationIDWebhooks_Response_Events, 0, len(eventStrings))
	for _, s := range eventStrings {
		item := api.PostAPIV1OrganizationsOrganizationIDWebhooks_Response_Events_Item{}
		item.PostAPIV1OrganizationsOrganizationIDWebhooks_Response_Events_AnyOf = &api.PostAPIV1OrganizationsOrganizationIDWebhooks_Response_Events_AnyOf{}
		_ = item.PostAPIV1OrganizationsOrganizationIDWebhooks_Response_Events_AnyOf.FromString(s)
		events = append(events, item)
	}
	return events
}

func makePatchEventsResponse(eventStrings []string) api.PatchAPIV1WebhooksWebhookID_Response_Events {
	events := make(api.PatchAPIV1WebhooksWebhookID_Response_Events, 0, len(eventStrings))
	for _, s := range eventStrings {
		item := api.PatchAPIV1WebhooksWebhookID_Response_Events_Item{}
		item.PatchAPIV1WebhooksWebhookID_Response_Events_AnyOf = &api.PatchAPIV1WebhooksWebhookID_Response_Events_AnyOf{}
		_ = item.PatchAPIV1WebhooksWebhookID_Response_Events_AnyOf.FromString(s)
		events = append(events, item)
	}
	return events
}

func TestWebhookDeliveryResponseFields(t *testing.T) {
	resp := &api.GetAPIV1WebhooksWebhookIDDeliveriesResponse{
		Items: api.GetAPIV1WebhooksWebhookIDDeliveries_Response_Items{
			{
				Payload: struct{}{},
				Request: api.GetAPIV1WebhooksWebhookIDDeliveries_Response_Items_Request{
					Body: `{"eventId":"123","type":"ping"}`,
					Headers: map[string]string{
						"Content-Type":          "application/json",
						"X-Signature-Ed25519":   "abc123",
						"X-Signature-Timestamp": "1735689600",
					},
				},
				Response: api.GetAPIV1WebhooksWebhookIDDeliveries_Response_Items_Response{
					Body: "OK",
					Headers: map[string]string{
						"Connection": "close",
						"Date":       "Wed, 01 Jan 2025 00:00:00 GMT",
					},
					Status: 200,
				},
				SentAt:  time.Now(),
				Success: true,
			},
		},
	}

	got := deliveriesPageFromResponse(resp)
	item := got.Items[0]

	if item.Request.Headers["Content-Type"] != "application/json" {
		t.Errorf("Request.Headers['Content-Type'] = %q, want 'application/json'", item.Request.Headers["Content-Type"])
	}
	if item.Request.Headers["X-Signature-Ed25519"] != "abc123" {
		t.Errorf("Request.Headers['X-Signature-Ed25519'] = %q, want 'abc123'", item.Request.Headers["X-Signature-Ed25519"])
	}
	if item.Response.Body != "OK" {
		t.Errorf("Response.Body = %q, want 'OK'", item.Response.Body)
	}
	if item.Response.Status != 200 {
		t.Errorf("Response.Status = %d, want 200", item.Response.Status)
	}
}

func TestConvertWebhookEvents(t *testing.T) {
	getAnyOf := func(e api.GetAPIV1WebhooksWebhookID_Response_Events_Item) rawUnion {
		if e.GetAPIV1WebhooksWebhookID_Response_Events_AnyOf == nil {
			return nil
		}
		return e.GetAPIV1WebhooksWebhookID_Response_Events_AnyOf
	}

	t.Run("nil anyOf skipped", func(t *testing.T) {
		events := []api.GetAPIV1WebhooksWebhookID_Response_Events_Item{
			{GetAPIV1WebhooksWebhookID_Response_Events_AnyOf: nil},
		}
		got := convertWebhookEvents(events, getAnyOf)
		if len(got) != 0 {
			t.Errorf("got %d events, want 0", len(got))
		}
	})

	t.Run("converts valid events", func(t *testing.T) {
		events := makeGetEventsResponse([]string{"asset.changed", "ping"})
		got := convertWebhookEvents(events, getAnyOf)
		if len(got) != 2 {
			t.Fatalf("got %d events, want 2", len(got))
		}
		if got[0] != WebhookEventTypeAssetChanged {
			t.Errorf("got[0] = %q, want %q", got[0], WebhookEventTypeAssetChanged)
		}
		if got[1] != WebhookEventTypePing {
			t.Errorf("got[1] = %q, want %q", got[1], WebhookEventTypePing)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		got := convertWebhookEvents([]api.GetAPIV1WebhooksWebhookID_Response_Events_Item{}, getAnyOf)
		if len(got) != 0 {
			t.Errorf("got %d events, want 0", len(got))
		}
	})

	t.Run("invalid JSON skipped", func(t *testing.T) {
		item := api.GetAPIV1WebhooksWebhookID_Response_Events_Item{}
		item.GetAPIV1WebhooksWebhookID_Response_Events_AnyOf = &api.GetAPIV1WebhooksWebhookID_Response_Events_AnyOf{}
		got := convertWebhookEvents([]api.GetAPIV1WebhooksWebhookID_Response_Events_Item{item}, getAnyOf)
		if len(got) != 0 {
			t.Errorf("got %d events, want 0 (empty/invalid raw should be skipped)", len(got))
		}
	})
}

func TestWebhookDeliveryPayloadIsAny(t *testing.T) {
	delivery := WebhookDelivery{
		Payload: map[string]any{
			"eventId": "123",
			"type":    "ping",
		},
	}

	data, err := json.Marshal(delivery)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled WebhookDelivery
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	payloadMap, ok := unmarshaled.Payload.(map[string]any)
	if !ok {
		t.Fatalf("Payload is not map[string]any, got %T", unmarshaled.Payload)
	}
	if payloadMap["type"] != "ping" {
		t.Errorf("Payload['type'] = %v, want 'ping'", payloadMap["type"])
	}
}
