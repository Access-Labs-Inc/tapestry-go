package tapestry

import (
	"encoding/json"
	"testing"
)

func TestUnixTimestamp_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    UnixTimestamp
		wantErr bool
	}{
		{
			name:  "regular integer timestamp",
			input: "1234567890",
			want:  UnixTimestamp(1234567890),
		},
		{
			name:  "low/high timestamp",
			input: `{"low":-188638304,"high":402}`,
			want:  UnixTimestamp(1730683181984),
		},
		{
			name:    "invalid format",
			input:   `"invalid"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got UnixTimestamp
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
