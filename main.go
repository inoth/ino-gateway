package main

import "encoding/json"

type UserInfo struct {
	Name string
}

func main() {
	user := UserInfo{Name: "123"}
	buf, _ := json.Marshal(user)
	println(string(buf))

	user2 := JsonMarshal[UserInfo](string(buf))
	println(user2.Name)
}

func JsonMarshal[T interface{}](str string) T {
	var res T
	json.Unmarshal([]byte(str), &res)
	return res
}
