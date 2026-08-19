package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

type CustomDate time.Time

func (this *CustomDate) ParseDate(b []byte, date string) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := time.Parse(date, s)
	if err != nil {
		return fmt.Errorf("error parsing date: %w", err)
	}
	*this = CustomDate(parsed)
	return nil
}
