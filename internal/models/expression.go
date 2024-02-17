package models

import (
	"encoding/json"

	"github.com/twharmon/gouid"
)

type Expression struct {
	Result            int    `json:"result"`
	Expression        string `json:"expression"`
	GUID              string `json:"guid"`
	ExecuteTime       int64  `json:"execute-time"`
	ExpectExucuteTime int64  `json:"expect-execute-time"`
	Err               error  `json:"err"`
	IsDone            bool   `json:"is-done"`
	WorkerID          string `json:"worker-id"`
}

func NewExpression(exp string) *Expression {
	guid := gouid.Bytes(16)
	return &Expression{Expression: exp, GUID: guid.String(), IsDone: false}
}

func (e *Expression) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Expression) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &e)
}
