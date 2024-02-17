package models

import "encoding/json"

type WorkerInfo struct {
	WorkerID   string `json:"worker-id" mapstructure:"worker-id"`
	LastTouch  string `json:"last-touch" mapstructure:"last-touch"` // когда в него последний раз что-то приходило и он этого отдавал
	IsEmploy   bool   `json:"is-employ" mapstructure:"is-employ"`
	CurrentJob string `json:"current-job" mapstructure:"current-job"`
	PrevJob    string `json:"prev-job" mapstructure:"prev-job"`
}

func (w *WorkerInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(w)
}

func (w *WorkerInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, w)
}
