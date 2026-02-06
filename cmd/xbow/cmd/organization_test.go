package cmd

import (
	"reflect"
	"testing"

	"github.com/rsclarke/xbow"
)

func TestParseMembers(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    []xbow.OrganizationMember
		wantErr bool
	}{
		{
			name:  "single member",
			input: []string{"email=alice@example.com,name=Alice"},
			want: []xbow.OrganizationMember{
				{Email: "alice@example.com", Name: "Alice"},
			},
		},
		{
			name: "multiple members",
			input: []string{
				"email=alice@example.com,name=Alice",
				"email=bob@example.com,name=Bob",
			},
			want: []xbow.OrganizationMember{
				{Email: "alice@example.com", Name: "Alice"},
				{Email: "bob@example.com", Name: "Bob"},
			},
		},
		{
			name:    "missing email",
			input:   []string{"name=Alice"},
			wantErr: true,
		},
		{
			name:    "missing name",
			input:   []string{"email=alice@example.com"},
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   []string{""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMembers(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseMembers() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseMembers() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
