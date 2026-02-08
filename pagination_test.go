package xbow

import (
	"context"
	"errors"
	"testing"
)

func ptr(s string) *string { return &s }

func TestPaginate(t *testing.T) {
	t.Run("iterates through multiple pages", func(t *testing.T) {
		pages := []*Page[string]{
			{Items: []string{"a", "b"}, PageInfo: PageInfo{NextCursor: ptr("cursor1"), HasMore: true}},
			{Items: []string{"c", "d"}, PageInfo: PageInfo{NextCursor: ptr("cursor2"), HasMore: true}},
			{Items: []string{"e"}, PageInfo: PageInfo{HasMore: false}},
		}
		callCount := 0

		fetch := func(ctx context.Context, opts *ListOptions) (*Page[string], error) {
			idx := callCount
			callCount++
			return pages[idx], nil
		}

		got, err := Collect(paginate(context.Background(), nil, fetch))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := []string{"a", "b", "c", "d", "e"}
		if len(got) != len(want) {
			t.Fatalf("got %d items, want %d", len(got), len(want))
		}
		for i, v := range got {
			if v != want[i] {
				t.Errorf("got[%d] = %q, want %q", i, v, want[i])
			}
		}

		if callCount != 3 {
			t.Errorf("fetch called %d times, want 3", callCount)
		}
	})

	t.Run("respects initial cursor from opts", func(t *testing.T) {
		var receivedCursor string
		fetch := func(ctx context.Context, opts *ListOptions) (*Page[string], error) {
			receivedCursor = opts.After
			return &Page[string]{Items: []string{"x"}, PageInfo: PageInfo{HasMore: false}}, nil
		}

		opts := &ListOptions{After: "start-here"}
		_, _ = Collect(paginate(context.Background(), opts, fetch))

		if receivedCursor != "start-here" {
			t.Errorf("cursor = %q, want 'start-here'", receivedCursor)
		}
	})

	t.Run("passes cursor between pages", func(t *testing.T) {
		cursors := []string{}
		callCount := 0

		fetch := func(ctx context.Context, opts *ListOptions) (*Page[string], error) {
			cursors = append(cursors, opts.After)
			callCount++
			if callCount == 1 {
				return &Page[string]{Items: []string{"a"}, PageInfo: PageInfo{NextCursor: ptr("next"), HasMore: true}}, nil
			}
			return &Page[string]{Items: []string{"b"}, PageInfo: PageInfo{HasMore: false}}, nil
		}

		_, _ = Collect(paginate(context.Background(), nil, fetch))

		if len(cursors) != 2 {
			t.Fatalf("expected 2 calls, got %d", len(cursors))
		}
		if cursors[0] != "" {
			t.Errorf("first cursor = %q, want empty", cursors[0])
		}
		if cursors[1] != "next" {
			t.Errorf("second cursor = %q, want 'next'", cursors[1])
		}
	})

	t.Run("propagates fetch errors", func(t *testing.T) {
		expectedErr := errors.New("fetch failed")
		fetch := func(ctx context.Context, opts *ListOptions) (*Page[string], error) {
			return nil, expectedErr
		}

		_, err := Collect(paginate(context.Background(), nil, fetch))
		if !errors.Is(err, expectedErr) {
			t.Errorf("error = %v, want %v", err, expectedErr)
		}
	})

	t.Run("stops early when yield returns false", func(t *testing.T) {
		callCount := 0
		fetch := func(ctx context.Context, opts *ListOptions) (*Page[string], error) {
			callCount++
			return &Page[string]{
				Items:    []string{"a", "b", "c"},
				PageInfo: PageInfo{NextCursor: ptr("next"), HasMore: true},
			}, nil
		}

		iter := paginate(context.Background(), nil, fetch)
		count := 0
		for _, err := range iter {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			count++
			if count == 2 {
				break
			}
		}

		if count != 2 {
			t.Errorf("iterated %d times, want 2", count)
		}
		if callCount != 1 {
			t.Errorf("fetch called %d times, want 1", callCount)
		}
	})

	t.Run("errors when HasMore is true but cursor is nil", func(t *testing.T) {
		fetch := func(ctx context.Context, opts *ListOptions) (*Page[string], error) {
			return &Page[string]{
				Items:    []string{"a"},
				PageInfo: PageInfo{HasMore: true, NextCursor: nil},
			}, nil
		}

		got, err := Collect(paginate(context.Background(), nil, fetch))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "xbow: server indicated more pages but returned no cursor" {
			t.Errorf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Errorf("got %d items before error, want 1", len(got))
		}
	})

	t.Run("errors when HasMore is true but cursor is empty", func(t *testing.T) {
		fetch := func(ctx context.Context, opts *ListOptions) (*Page[string], error) {
			return &Page[string]{
				Items:    []string{"a"},
				PageInfo: PageInfo{HasMore: true, NextCursor: ptr("")},
			}, nil
		}

		got, err := Collect(paginate(context.Background(), nil, fetch))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "xbow: server indicated more pages but returned no cursor" {
			t.Errorf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Errorf("got %d items before error, want 1", len(got))
		}
	})

	t.Run("errors when cursor does not advance", func(t *testing.T) {
		callCount := 0
		fetch := func(ctx context.Context, opts *ListOptions) (*Page[string], error) {
			callCount++
			return &Page[string]{
				Items:    []string{"a"},
				PageInfo: PageInfo{HasMore: true, NextCursor: ptr("same-cursor")},
			}, nil
		}

		got, err := Collect(paginate(context.Background(), &ListOptions{After: "same-cursor"}, fetch))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "xbow: server returned same cursor, stopping to prevent infinite loop" {
			t.Errorf("unexpected error: %v", err)
		}
		if callCount != 1 {
			t.Errorf("fetch called %d times, want 1", callCount)
		}
		if len(got) != 1 {
			t.Errorf("got %d items before error, want 1", len(got))
		}
	})

	t.Run("handles empty page", func(t *testing.T) {
		fetch := func(ctx context.Context, opts *ListOptions) (*Page[string], error) {
			return &Page[string]{Items: []string{}, PageInfo: PageInfo{HasMore: false}}, nil
		}

		got, err := Collect(paginate(context.Background(), nil, fetch))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("got %d items, want 0", len(got))
		}
	})
}

func TestCollect(t *testing.T) {
	t.Run("collects all items", func(t *testing.T) {
		seq := func(yield func(int, error) bool) {
			for i := 1; i <= 3; i++ {
				if !yield(i, nil) {
					return
				}
			}
		}

		got, err := Collect(seq)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := []int{1, 2, 3}
		if len(got) != len(want) {
			t.Fatalf("got %d items, want %d", len(got), len(want))
		}
	})

	t.Run("returns partial results on error", func(t *testing.T) {
		expectedErr := errors.New("mid-stream error")
		seq := func(yield func(int, error) bool) {
			yield(1, nil)
			yield(2, nil)
			yield(0, expectedErr)
		}

		got, err := Collect(seq)
		if !errors.Is(err, expectedErr) {
			t.Errorf("error = %v, want %v", err, expectedErr)
		}
		if len(got) != 2 {
			t.Errorf("got %d items before error, want 2", len(got))
		}
	})
}
