package models

import (
	"encoding/json"
	"time"
)

type NullableTime struct {
	time.Time
}

func newNullableTime(t time.Time) NullableTime {
	return NullableTime{
		Time: t,
	}
}

func (nt *NullableTime) MarshalJSON() ([]byte, error) {
	if nt.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Time = time.Time{}
		return nil
	}
	return json.Unmarshal(data, &nt.Time)
}
