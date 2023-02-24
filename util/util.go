package util

import "encoding/json"

func JsonMarshal[T interface{}](str string) T {
	var res T
	json.Unmarshal([]byte(str), &res)
	return res
}
