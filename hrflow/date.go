package hrflow

import (
	"strings"
	"time"
)

type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(hrFlowDateFormat, s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(d).Format(hrFlowDateFormat) + "\""), nil
}

func (d Date) Format(s string) string {
	t := time.Time(d)
	return t.Format(s)
}
