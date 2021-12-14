package middleware

import (
	"testing"
)

func Test_mimeType(t *testing.T) {
	tests := []struct {
		contentType string
		want        string
	}{
		{"application/json", "application/json"},
		{"APPLICATION/JSON", "application/json"},
		{"APPLICATION/json; CHARSET=UTF-8", "application/json"},
		{"application/json; charset=UTF-8", "application/json"},
		{"application/json; boundary=foo; charset=UTF-8", "application/json"},
		{"application/json;BOUNDARY=foo; charset=UTF-8", "application/json"},
		{"foo/bar; Charset=\"utf-8\"", "foo/bar"},
		{"multipart/form-data; application/json", "multipart/form-data"},
		{"multipart/form-data, application/json; charset=utf-8", "multipart/form-data, application/json"},
	}

	for _, tt := range tests {
		t.Run(tt.contentType, func(t *testing.T) {
			if got := mimeType(tt.contentType); got != tt.want {
				t.Errorf("mimeType() = %v, want %v", got, tt.want)
			}
		})
	}
}
