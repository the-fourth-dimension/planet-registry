package lib

import (
	"bytes"
	"encoding/json"
)

func MakeAuthHeader(token string) string {
	return "Bearer " + token
}

func SerializeBody(body interface{}) *bytes.Buffer {
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(body)
	return &b
}
