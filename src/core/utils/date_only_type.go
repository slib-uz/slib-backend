package utils

import (
	"encoding/json"
	"strings"
	"time"
)

type DateOnlyType struct {
	time.Time
}

func (this *DateOnlyType) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	this.Time = t
	return nil
}

func (this *DateOnlyType) MarshalJSON() ([]byte, error) {
	s := this.Format("2006-01-02")
	return json.Marshal(s)
}
