package sorter

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/red3533/log-analyzer/internal/logger"
	"github.com/red3533/log-analyzer/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSort_OnlySort(t *testing.T) {
	cases := []struct {
		name      string
		sortField string
		sortBy    string
		logs      []models.Log
		want      []models.Log
		wantErr   bool
	}{
		{
			name:      "correct_ip_sort",
			sortField: "ip",
			sortBy:    "desc",
			logs: []models.Log{
				{IP: "33.632.52.43"},
				{IP: "42.324.627.35.22"},
				{IP: "42.324.627.35.2"},
			},
			want: []models.Log{
				{IP: "42.324.627.35.22"},
				{IP: "42.324.627.35.2"},
				{IP: "33.632.52.43"},
			},
		},
		{
			name:      "correct_timestamp_sort",
			sortField: "timestamp",
			sortBy:    "asc",
			logs: []models.Log{
				{Timestamp: time.Date(2000, time.May, 30, 5, 20, 10, 0, time.UTC)},
				{Timestamp: time.Date(3000, time.May, 30, 5, 20, 10, 0, time.UTC)},
				{Timestamp: time.Date(3000, time.May, 30, 5, 20, 10, 1, time.UTC)},
			},
			want: []models.Log{
				{Timestamp: time.Date(2000, time.May, 30, 5, 20, 10, 0, time.UTC)},
				{Timestamp: time.Date(3000, time.May, 30, 5, 20, 10, 0, time.UTC)},
				{Timestamp: time.Date(3000, time.May, 30, 5, 20, 10, 1, time.UTC)},
			},
		},
		{
			name:      "correct_method_sort",
			sortField: "method",
			sortBy:    "asc",
			logs: []models.Log{
				{Method: "POST"},
				{Method: "GET"},
				{Method: "DELETE"},
			},
			want: []models.Log{
				{Method: "DELETE"},
				{Method: "GET"},
				{Method: "POST"},
			},
		},
		{
			name:      "correct_url_sort",
			sortField: "url",
			sortBy:    "asc",
			logs: []models.Log{
				{URL: "/favicon.ico"},
				{URL: "/dashboard"},
				{URL: "/products/2"},
			},
			want: []models.Log{
				{URL: "/dashboard"},
				{URL: "/favicon.ico"},
				{URL: "/products/2"},
			},
		},
		{
			name:      "correct_status_sort",
			sortField: "status",
			sortBy:    "desc",
			logs: []models.Log{
				{Status: 200},
				{Status: 202},
				{Status: 404},
			},
			want: []models.Log{
				{Status: 404},
				{Status: 202},
				{Status: 200},
			},
		},
		{
			name:      "correct_size_byte_sort",
			sortField: "size_byte",
			sortBy:    "desc",
			logs: []models.Log{
				{SizeByte: 500},
				{SizeByte: 0},
				{SizeByte: 503},
			},
			want: []models.Log{
				{SizeByte: 503},
				{SizeByte: 500},
				{SizeByte: 0},
			},
		},
	}

	mockLogger := logger.Logger{}
	logSorter := NewLogSorter(mockLogger)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logs := tc.logs

			err := logSorter.Sort(logs, tc.sortField, tc.sortBy)
			if tc.wantErr {
				require.Error(t, err)
			}

			require.NoError(t, err)

			if !assert.True(t, reflect.DeepEqual(tc.want, logs)) {
				fmt.Printf("Want: %v\nGot: %v\n", tc.want, logs)
			}

		})
	}

}

func TestSort_Error(t *testing.T) {
	cases := []struct {
		name      string
		sortField string
		sortBy    string
		logs      []models.Log
		wantErr   bool
	}{
		{
			name:      "incorrect_ip_upcase",
			sortField: "IP",
			logs:      []models.Log{{IP: "33.632.52.43"}},
			wantErr:   true,
		},
	}

	mockLogger := logger.Logger{}
	logSorter := NewLogSorter(mockLogger)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logs := tc.logs

			err := logSorter.Sort(logs, tc.sortField, tc.sortBy)
			require.Error(t, err)

		})
	}

}
