package utils

import "testing"

func TestAddQueryParam(t *testing.T) {
	type args struct {
		u string
		k string
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty url",
			args: args{
				u: "",
				k: "k",
				v: "v",
			},
			want: "/?k=v",
		},
		{
			name: "empty key",
			args: args{
				u: "http://example.com",
				k: "",
				v: "v",
			},
			want: "http://example.com?=v",
		},
		{
			name: "empty value",
			args: args{
				u: "http://example.com",
				k: "k",
				v: "",
			},
			want: "http://example.com?k=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddQueryParam(tt.args.u, tt.args.k, tt.args.v); got != tt.want {
				t.Errorf("AddQueryParam() = %v, want %v", got, tt.want)
			}
		})
	}
}
