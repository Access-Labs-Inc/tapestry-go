package tapestry

import (
	"encoding/json"
	"testing"
)

func TestComment_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Comment
		wantErr bool
	}{
		{
			name:  "regular integer timestamp",
			input: `{"namespace":"test","created_at":1234567890,"text":"hello","id":"123"}`,
			want: Comment{
				Namespace: "test",
				CreatedAt: 1234567890,
				Text:      "hello",
				ID:        "123",
			},
		},
		{
			name:  "low/high timestamp",
			input: `{"namespace":"test","created_at":{"low":-188638304,"high":402},"text":"hello","id":"456"}`,
			want: Comment{
				Namespace: "test",
				CreatedAt: 1730683181984,
				Text:      "hello",
				ID:        "456",
			},
		},
		{
			name:    "invalid timestamp format",
			input:   `{"namespace":"test","created_at":"invalid","text":"hello","id":"789"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Comment
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("UnmarshalJSON() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
