package models

type Model interface {
	MarshalBinary() ([]byte, error)
	UnmarshalBinary(data []byte) error
}
