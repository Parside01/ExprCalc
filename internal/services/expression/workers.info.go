package expression

import "encoding/json"

type workerInfo struct {
	WorkerID   string `json:"worker-id" mapstructure:"worker-id"`
	LastTouch  string `json:"last-touch" mapstructure:"last-touch"`
	IsEmploy   bool   `json:"is-employ" mapstructure:"is-employ"`
	CurrentJob string `json:"current-job" mapstructure:"current-job"`
}

func (w *workerInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(w)
}

func (w *workerInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, w)
}
