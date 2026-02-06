package cmd

import (
	"reflect"
	"testing"

	"github.com/rsclarke/xbow"
)

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    map[string][]string
		wantErr bool
	}{
		{
			name:  "nil input",
			input: nil,
			want:  nil,
		},
		{
			name:  "empty input",
			input: []string{},
			want:  nil,
		},
		{
			name:  "single header",
			input: []string{"X-Custom: value"},
			want:  map[string][]string{"X-Custom": {"value"}},
		},
		{
			name:  "multiple different headers",
			input: []string{"X-One: a", "X-Two: b"},
			want:  map[string][]string{"X-One": {"a"}, "X-Two": {"b"}},
		},
		{
			name:  "multiple values same key",
			input: []string{"X-Custom: val1", "X-Custom: val2"},
			want:  map[string][]string{"X-Custom": {"val1", "val2"}},
		},
		{
			name:  "trims whitespace",
			input: []string{"  X-Custom  :  value  "},
			want:  map[string][]string{"X-Custom": {"value"}},
		},
		{
			name:  "value with colon",
			input: []string{"Authorization: Bearer token:123"},
			want:  map[string][]string{"Authorization": {"Bearer token:123"}},
		},
		{
			name:    "missing colon",
			input:   []string{"InvalidHeader"},
			wantErr: true,
		},
		{
			name:    "empty key",
			input:   []string{": value"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHeaders(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseHeaders() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseKV(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]string
	}{
		{
			name:  "simple pair",
			input: "key=value",
			want:  map[string]string{"key": "value"},
		},
		{
			name:  "multiple pairs",
			input: "a=1,b=2,c=3",
			want:  map[string]string{"a": "1", "b": "2", "c": "3"},
		},
		{
			name:  "trims whitespace",
			input: " key = value , foo = bar ",
			want:  map[string]string{"key": "value", "foo": "bar"},
		},
		{
			name:  "empty string",
			input: "",
			want:  map[string]string{},
		},
		{
			name:  "no equals sign skipped",
			input: "a=1,noequals,b=2",
			want:  map[string]string{"a": "1", "b": "2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseKV(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseKV(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseCredentials(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    []xbow.Credential
		wantErr bool
	}{
		{
			name:  "single basic credential",
			input: []string{"name=admin,type=basic,username=user,password=pass"},
			want: []xbow.Credential{
				{Name: "admin", Type: "basic", Username: "user", Password: "pass"},
			},
		},
		{
			name: "multiple credentials",
			input: []string{
				"name=admin,type=basic,username=u1,password=p1",
				"name=viewer,type=basic,username=u2,password=p2",
			},
			want: []xbow.Credential{
				{Name: "admin", Type: "basic", Username: "u1", Password: "p1"},
				{Name: "viewer", Type: "basic", Username: "u2", Password: "p2"},
			},
		},
		{
			name:  "with optional fields",
			input: []string{"name=admin,type=basic,username=u,password=p,email-address=a@b.com,authenticator-uri=otpauth://totp/test"},
			want: []xbow.Credential{
				{
					Name:             "admin",
					Type:             "basic",
					Username:         "u",
					Password:         "p",
					EmailAddress:     strPtr("a@b.com"),
					AuthenticatorURI: strPtr("otpauth://totp/test"),
				},
			},
		},
		{
			name:  "with id field",
			input: []string{"id=cred-1,name=admin,type=basic,username=u,password=p"},
			want: []xbow.Credential{
				{ID: "cred-1", Name: "admin", Type: "basic", Username: "u", Password: "p"},
			},
		},
		{
			name:    "missing name",
			input:   []string{"type=basic,username=u,password=p"},
			wantErr: true,
		},
		{
			name:    "missing type",
			input:   []string{"name=admin,username=u,password=p"},
			wantErr: true,
		},
		{
			name:    "missing username",
			input:   []string{"name=admin,type=basic,password=p"},
			wantErr: true,
		},
		{
			name:    "missing password",
			input:   []string{"name=admin,type=basic,username=u"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCredentials(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCredentials() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseDNSRules(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    []xbow.DNSBoundaryRule
		wantErr bool
	}{
		{
			name:  "basic rule",
			input: []string{"action=allow-attack,type=hostname,filter=example.com"},
			want: []xbow.DNSBoundaryRule{
				{Action: xbow.DNSBoundaryRuleActionAllowAttack, Type: "hostname", Filter: "example.com"},
			},
		},
		{
			name:  "with include-subdomains true",
			input: []string{"action=allow-visit,type=hostname,filter=example.com,include-subdomains=true"},
			want: []xbow.DNSBoundaryRule{
				{Action: xbow.DNSBoundaryRuleActionAllowVisit, Type: "hostname", Filter: "example.com", IncludeSubdomains: boolPtr(true)},
			},
		},
		{
			name:  "with include-subdomains false",
			input: []string{"action=deny,type=hostname,filter=evil.com,include-subdomains=false"},
			want: []xbow.DNSBoundaryRule{
				{Action: xbow.DNSBoundaryRuleActionDeny, Type: "hostname", Filter: "evil.com", IncludeSubdomains: boolPtr(false)},
			},
		},
		{
			name:  "with id",
			input: []string{"id=rule-1,action=allow-attack,type=hostname,filter=example.com"},
			want: []xbow.DNSBoundaryRule{
				{ID: "rule-1", Action: xbow.DNSBoundaryRuleActionAllowAttack, Type: "hostname", Filter: "example.com"},
			},
		},
		{
			name: "multiple rules",
			input: []string{
				"action=allow-attack,type=hostname,filter=a.com",
				"action=deny,type=hostname,filter=b.com",
			},
			want: []xbow.DNSBoundaryRule{
				{Action: xbow.DNSBoundaryRuleActionAllowAttack, Type: "hostname", Filter: "a.com"},
				{Action: xbow.DNSBoundaryRuleActionDeny, Type: "hostname", Filter: "b.com"},
			},
		},
		{
			name:    "missing action",
			input:   []string{"type=hostname,filter=example.com"},
			wantErr: true,
		},
		{
			name:    "missing type",
			input:   []string{"action=deny,filter=example.com"},
			wantErr: true,
		},
		{
			name:    "missing filter",
			input:   []string{"action=deny,type=hostname"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDNSRules(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseDNSRules() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDNSRules() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseHTTPRules(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    []xbow.HTTPBoundaryRule
		wantErr bool
	}{
		{
			name:  "basic rule",
			input: []string{"action=deny,type=url,filter=https://evil.com"},
			want: []xbow.HTTPBoundaryRule{
				{Action: xbow.HTTPBoundaryRuleActionDeny, Type: "url", Filter: "https://evil.com"},
			},
		},
		{
			name:  "allow-auth action",
			input: []string{"action=allow-auth,type=url,filter=https://login.example.com"},
			want: []xbow.HTTPBoundaryRule{
				{Action: xbow.HTTPBoundaryRuleActionAllowAuth, Type: "url", Filter: "https://login.example.com"},
			},
		},
		{
			name:  "with include-subdomains",
			input: []string{"action=allow-attack,type=url,filter=https://example.com,include-subdomains=true"},
			want: []xbow.HTTPBoundaryRule{
				{Action: xbow.HTTPBoundaryRuleActionAllowAttack, Type: "url", Filter: "https://example.com", IncludeSubdomains: boolPtr(true)},
			},
		},
		{
			name:  "with id",
			input: []string{"id=rule-1,action=allow-visit,type=url,filter=https://example.com"},
			want: []xbow.HTTPBoundaryRule{
				{ID: "rule-1", Action: xbow.HTTPBoundaryRuleActionAllowVisit, Type: "url", Filter: "https://example.com"},
			},
		},
		{
			name:    "missing action",
			input:   []string{"type=url,filter=https://example.com"},
			wantErr: true,
		},
		{
			name:    "missing type",
			input:   []string{"action=deny,filter=https://example.com"},
			wantErr: true,
		},
		{
			name:    "missing filter",
			input:   []string{"action=deny,type=url"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHTTPRules(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseHTTPRules() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHTTPRules() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }
