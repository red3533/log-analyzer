package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/red3533/log-analyzer/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_WithTempFile(t *testing.T) {
	cases := []struct {
		name     string
		content  []string
		wantLogs int
		wantErr  bool
	}{
		{
			name: "correct_content",
			content: []string{
				`192.168.1.8 - - [31/May/2025:10:00:08 +0000] "GET /login HTTP/1.1" 200 256`,
				`192.168.1.18 - - [31/May/2025:10:00:18 +0000] "GET /favicon.ico HTTP/1.1" 404 0`,
				`192.168.1.19 - - [31/May/2025:10:00:19 +0000] "GET /robots.txt HTTP/1.1" 200 0`,
			},
			wantLogs: 3,
			wantErr:  false,
		},
		{
			name: "invalid_log_but_ok",
			content: []string{
				`192.168.1.15 - - [31/May/2025:10:00:15 +0000] "PUT /api/users/2 HTTP/1.1" 200 512`,
				`warn_1`,
				`192.168.1.1 - - [31/May/2025:10:00:01 +0000] "GET /index.html HTTP/1.1" 200 1024`,
				`warn_2`,
				`warn_3`,
			},
			wantLogs: 2,
			wantErr:  false,
		},
		{
			name:     "empty_log_file",
			content:  []string{},
			wantLogs: 0,
			wantErr:  false,
		},
	}

	nginxParser := NewNginxParser(logger.Logger{}, DefaultFileReader{})

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tmpFile := filepath.Join(t.TempDir(), "test.log")
			os.WriteFile(tmpFile, []byte(strings.Join(tc.content, "\n")), 0644)

			out, err := nginxParser.Parse(tmpFile)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, out, tc.wantLogs)

		})
	}
}

func TestExtractIP_EdgeCases(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "correct_ip",
			input:   `192.168.1.8 - - [31/May/2025:10:00:08 +0000] "GET /login HTTP/1.1" 200 256`,
			want:    "192.168.1.8",
			wantErr: false,
		},
		{
			name:    "incorrect_ip_string",
			input:   `sjhdfg - - [31/May/2025:10:00:15 +0000] "PUT /api/users/2 HTTP/1.1" 200 512`,
			wantErr: true,
		},
		{
			name:    "empty_ip",
			input:   `- - [31/May/2025:10:00:09 +0000] "POST /login HTTP/1.1" 302 0`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := extractIP(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)

		})
	}
}

func TestExtractTimestamp_EdgeCases(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "correct_timestamp",
			input:   `192.168.1.8 - - [31/May/2025:10:00:08 -0400] "GET /login HTTP/1.1" 200 256`,
			want:    time.Date(2025, time.May, 31, 10, 0, 8, 0, time.FixedZone("", -4*60*60)),
			wantErr: false,
		},
		{
			name:    "incorrect_timestamp_1000_may",
			input:   `192.168.1.4 - - [1000/May/2025:10:00:04 +0000] "GET /contact.html HTTP/1.1" 200 768`,
			wantErr: true,
		},
		{
			name:    "empty_timestamp",
			input:   `192.168.1.20 - - "GET /sitemap.xml HTTP/1.1" 200 1024`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := extractTimestamp(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)

		})
	}
}

func TestExtractMethod_EdgeCases(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "correct_method",
			input:   `192.168.1.10 - - [31/May/2025:10:00:10 +0000] "GET /dashboard HTTP/1.1" 200 2048`,
			want:    "GET",
			wantErr: false,
		},
		{
			name:    "incorrect_method_stop",
			input:   `192.168.1.4 - - [31/May/2025:10:00:04 +0000] "STOP /contact.html HTTP/1.1" 200 768`,
			wantErr: true,
		},
		{
			name:    "empty_method",
			input:   `192.168.1.18 - - [31/May/2025:10:00:18 +0000] "/favicon.ico HTTP/1.1" 404 0`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := extractMethod(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)

		})
	}
}
func TestExtractURL_EdgeCases(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "correct_url",
			input:   `192.168.1.7 - - [31/May/2025:10:00:07 +0000] "GET /products/2 HTTP/1.1" 404 512`,
			want:    "/products/2",
			wantErr: false,
		},
		{
			name:    "incorrect_url_with_space",
			input:   `192.168.1.13 - - [31/May/2025:10:00:13 +0000] "GET /api/use rs/1 HTTP/1.1" 200 1024`,
			wantErr: true,
		},
		{
			name:    "empty_url",
			input:   `192.168.1.1 - - [31/May/2025:10:00:01 +0000] "GET HTTP/1.1" 200 1024`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := extractURL(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)

		})
	}
}

func TestExtractStatus_EdgeCases(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{
			name:    "correct_status",
			input:   `192.168.1.8 - - [31/May/2025:10:00:08 +0000] "GET /login HTTP/1.1" 200 256`,
			want:    200,
			wantErr: false,
		},
		{
			name:    "incorrect_status_more_than_599",
			input:   `192.168.1.4 - - [31/May/2025:10:00:04 +0000] "GET /contact.html HTTP/1.1" 3000 768`,
			wantErr: true,
		},
		{
			name:    "empty_status",
			input:   `192.168.1.20 - - "GET /sitemap.xml HTTP/1.1" 1024`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := extractStatus(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)

		})
	}
}

func TestExtractSizeByte_EdgeCases(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{
			name:    "correct_size_byte",
			input:   `192.168.1.12 - - [31/May/2025:10:00:12 +0000] "GET /api/users HTTP/1.1" 200 512`,
			want:    512,
			wantErr: false,
		},
		{
			name:    "incorrect_size_byte_less_0",
			input:   `192.168.1.6 - - [31/May/2025:10:00:06 +0000] "GET /products/1 HTTP/1.1" 200 -1`,
			wantErr: true,
		},
		{
			name:    "empty_size_byte",
			input:   `192.168.1.2 - - [31/May/2025:10:00:02 +0000] "POST /api/data HTTP/1.1" 201`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := extractSizeByte(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)

		})
	}
}
