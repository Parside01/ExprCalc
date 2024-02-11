package models

import (
	"encoding/json"
	"time"

	"github.com/twharmon/gouid"
)

type Expression struct {
	Result      int           `json:"result"`
	Expression  string        `json:"expression"`
	GUID        string        `json:"guid"`
	ExecuteTime time.Duration `json:"execute-time"`
	Err         error         `json:"err"`
}

func NewExpression(exp string) *Expression {
	guid := gouid.Bytes(16)
	return &Expression{Expression: exp, GUID: guid.String()}
}

func (e *Expression) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Expression) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &e)
}
