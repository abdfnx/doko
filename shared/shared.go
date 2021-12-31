package shared

import (
	"bytes"
	"encoding/json"
	"sort"
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
