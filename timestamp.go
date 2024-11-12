package tapestry

import (
	"encoding/json"
	"fmt"
)

// UnixTimestamp represents a unix timestamp that can be unmarshaled from either
// an integer or a 64 bit integer representation in object format
type UnixTimestamp int64

func (t *UnixTimestamp) UnmarshalJSON(data []byte) error {
	// Try regular int64 first
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err == nil {
		*t = UnixTimestamp(timestamp)
		return nil
	}

	// Try 64 bit integer representation
	var tsObj struct {
		Low  int64 `json:"low"`
		High int64 `json:"high"`
	}
	if err := json.Unmarshal(data, &tsObj); err != nil {
		return fmt.Errorf("failed to unmarshal timestamp: %w", err)
	}

	*t = UnixTimestamp((tsObj.High << 32) | (tsObj.Low & 0xFFFFFFFF))
	return nil
}
