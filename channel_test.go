package rc

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestHistoryQueryOptions_Q(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		h    *HistoryQuery
		want url.Values
	}{
		{
			name: "testquery",
			h: &HistoryQuery{
				RoomID:         "GENERAL",
				Latest:         &now,
				Inclusive:      true,
				Offset:         20,
				Count:          100,
				IncludeUnreads: true,
			},
			want: url.Values{
				"count":     []string{"100"},
				"unreads":   []string{"true"},
				"roomId":    []string{"GENERAL"},
				"latest":    []string{now.Format("2006-01-02T15:04:05.999Z")},
				"inclusive": []string{"true"},
				"offset":    []string{"20"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Q(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HistoryQueryOptions.Q() = %#v, want %v", got, tt.want)
			}
		})
	}
}
