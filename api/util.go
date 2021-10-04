package api

import (
	"encoding/json"
	"time"
)

type PlusServerTime time.Time

func (m *PlusServerTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	return json.Unmarshal(data, (*time.Time)(m))
}