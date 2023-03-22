package utils_test

import (
	"testing"

	"github.com/dmitrymomot/oauth2-server/internal/utils"
)

func TestTrimStringBetween(t *testing.T) {
	type args struct {
		str   string
		start string
		end   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "with counter",
			args: args{
				str:   "[16] test string trim",
				start: "[",
				end:   "]",
			},
			want: "test string trim",
		},
		{
			name: "without counter",
			args: args{
				str:   " test string trim",
				start: "[",
				end:   "]",
			},
			want: "test string trim",
		},
		{
			name: "without counter and spaces",
			args: args{
				str:   "test string trim",
				start: "[",
				end:   "]",
			},
			want: "test string trim",
		},
		{
			name: "with empty counter",
			args: args{
				str:   "[] test string trim",
				start: "[",
				end:   "]",
			},
			want: "test string trim",
		},
		{
			name: "with 2 counters",
			args: args{
				str:   "[16] test string [12] trim",
				start: "[",
				end:   "]",
			},
			want: "test string [12] trim",
		},
		{
			name: "with counter in the middle",
			args: args{
				str:   "test string [12] trim",
				start: "[",
				end:   "]",
			},
			want: "test string  trim",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.TrimStringBetween(tt.args.str, tt.args.start, tt.args.end); got != tt.want {
				t.Errorf("TrimStringBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty string",
			args: args{
				s: "",
			},
			want: "",
		},
		{
			name: "one word",
			args: args{
				s: "test",
			},
			want: "test",
		},
		{
			name: "two words",
			args: args{
				s: "testString",
			},
			want: "test_string",
		},
		{
			name: "two words with spaces",
			args: args{
				s: "test String",
			},
			want: "test_string",
		},
		{
			name: "two words with spaces and numbers",
			args: args{
				s: "test String 123",
			},
			want: "test_string_123",
		},
		{
			name: "two words with spaces and numbers and special chars",
			args: args{
				s: "test String 123 !@#$%^&*()",
			},
			want: "test_string_123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.ToSnakeCase(tt.args.s); got != tt.want {
				t.Errorf("ToSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
