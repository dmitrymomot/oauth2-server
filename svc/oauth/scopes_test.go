package oauth_test

import (
	"testing"

	"github.com/dmitrymomot/oauth2-server/svc/oauth"
)

func TestMatchScope(t *testing.T) {
	type args struct {
		requiredScope string
		allowedScopes string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "simple",
			args: args{
				requiredScope: "user:read",
				allowedScopes: "user:read user:write",
			},
			want: true,
		},
		{
			name: "asterisk",
			args: args{
				requiredScope: "user:read",
				allowedScopes: "user:*",
			},
			want: true,
		},
		{
			name: "asterisk2",
			args: args{
				requiredScope: "user:read:1",
				allowedScopes: "user:read:*",
			},
			want: true,
		},
		{
			name: "merchant:member:read > merchant:*",
			args: args{
				requiredScope: "merchant:member:read",
				allowedScopes: "merchant:*",
			},
			want: true,
		},
		{
			name: "merchant:member:read > merchant:member:*",
			args: args{
				requiredScope: "merchant:member:read",
				allowedScopes: "merchant:member:*",
			},
			want: true,
		},
		{
			name: "merchant:member:read > *",
			args: args{
				requiredScope: "merchant:member:read",
				allowedScopes: "*",
			},
			want: true,
		},
		{
			name: "merchant:member:read > merchant:*:*",
			args: args{
				requiredScope: "merchant:member:read",
				allowedScopes: "merchant:*:*",
			},
			want: true,
		},
		{
			name: "merchant > merchant:*",
			args: args{
				requiredScope: "merchant",
				allowedScopes: "merchant:*",
			},
			want: false,
		},
		{
			name: "merchant:member > merchant:*",
			args: args{
				requiredScope: "merchant:member",
				allowedScopes: "merchant:*",
			},
			want: true,
		},
		{
			name: "merchant:member > merchant:member",
			args: args{
				requiredScope: "merchant:member",
				allowedScopes: "merchant:member",
			},
			want: true,
		},
		{
			name: "no allowed scopes",
			args: args{
				requiredScope: "merchant:member",
				allowedScopes: "",
			},
			want: false,
		},
		{
			name: "empty allowed scopes",
			args: args{
				requiredScope: "merchant:member",
				allowedScopes: " ",
			},
			want: false,
		},
		{
			name: "empty scopes",
			args: args{
				requiredScope: "",
				allowedScopes: "",
			},
			want: true,
		},
		{
			name: "empty required scopes",
			args: args{
				requiredScope: "",
				allowedScopes: "user:*",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := oauth.MatchScope(tt.args.requiredScope, tt.args.allowedScopes); got != tt.want {
				t.Errorf("MatchScope() = %v, want %v", got, tt.want)
			}
		})
	}
}
