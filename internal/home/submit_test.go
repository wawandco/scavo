package home

import (
	"strings"
	"testing"
)

func TestSanitizeUploadFilename(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "normal png",
			input: "photo.png",
			want:  "photo.png",
		},
		{
			name:  "path traversal",
			input: "../../../../../etc/passwd",
			want:  "passwd.bin", // no original ext → we append .bin for safety
		},
		{
			name:  "windows path",
			input: "C:\\Users\\pwned\\evil.jpg",
			want:  "C--Users-pwned-evil.jpg", // : and \ become -
		},
		{
			name:  "dangerous characters",
			input: "my file (1) [final]!.jpeg",
			want:  "my-file--1---final--.jpeg",
		},
		{
			name:  "no extension",
			input: "screenshot",
			want:  "screenshot.bin",
		},
		{
			name:  "very long name",
			input: "a" + strings.Repeat("b", 100) + ".png",
			want:  "a" + strings.Repeat("b", 75) + ".png", // 80 chars total incl. ext
		},
		{
			name:  "weird extension",
			input: "file.VERYLONGEXTENSION",
			want:  "file.bin",
		},
		{
			name:  "empty",
			input: "",
			want:  "upload.bin",
		},
		{
			name:  "dot only",
			input: ".",
			want:  "upload.bin",
		},
		{
			name:  "double dot traversal",
			input: "..",
			want:  "upload.bin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeUploadFilename(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeUploadFilename(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
