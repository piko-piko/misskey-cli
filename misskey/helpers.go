package misskey

import (
	"time"
)

func convert(timestamp string) (string) {
	t, _ := time.ParseInLocation("2006-01-02T15:04:05Z", timestamp, time.UTC)
	return t.In(time.Local).Format("2006/01/02 15:04:05")
}

