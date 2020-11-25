package hrflow

import (
	"strings"
	"time"
)

type JSONDate time.Time

func (d *JSONDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	*d = JSONDate(t)
	return nil
}

func (d JSONDate) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(d).Format(hrFlowDateFormat) + "\""), nil
}

func (d JSONDate) Format(s string) string {
	t := time.Time(d)
	return t.Format(s)
}
