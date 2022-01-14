package utils

import "encoding/json"

func StringAddr(s string) *string {
	return &s
}

func BoolAddr(b bool) *bool {
	return &b
}

func Int64Addr(i int64) *int64 {
	return &i
}

func Stringify(v interface{}) string {
	bytes, _ := json.MarshalIndent(v, "", " ")
	return string(bytes)
}
