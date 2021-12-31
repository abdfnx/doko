package shared

import (
	"os"
	"fmt"
	"sort"
	"time"
	"bytes"
	"strings"
	"encoding/json"
)

// StructToJSON -> convert any struct to json.
func StructToJSON(i interface{}) string {
	j, err := json.Marshal(i)

	if err != nil {
		return ""
	}

	out := new(bytes.Buffer)
	json.Indent(out, j, "", "    ")
	return out.String()
}

// SortKeys -> sort keys.
func SortKeys(keys []string) []string {
	sort.Strings(keys)
	return keys
}

// GetEnv -> get os environment.
func GetEnv(env string) string {
	keyval := strings.SplitN(env, "=", 2)

	if keyval[1][:1] == "$" {
		keyval[1] = os.Getenv(keyval[1][1:])
		return strings.Join(keyval, "=")
	}

	return env
}

// ParseDateToString -> parse date to string.
func ParseDateToString(unixtime int64) string {
	t := time.Unix(unixtime, 0)
	return t.Format("2006/01/02 15:04:05")
}

// ParseSizeToString -> parse size to string.
func ParseSizeToString(size int64) string {
	mb := float64(size) / 1024 / 1024
	return fmt.Sprintf("%.1fMB", mb)
}
