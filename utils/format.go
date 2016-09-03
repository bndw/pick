package utils

import (
	"encoding/json"
	"os"
	"time"
)

func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	os.Stdout.Write(b)
}

func FormatUnixTime(ts int64) string {
	t := time.Unix(ts, 0)
	return t.Format("Mon Jan 2 15:04:05 2006")
}
