package models

type MessageRabbit struct {
	Exchange string `json:"exchange"`
	Queue    string `json:"queue"`
	Payload  []byte `json:"payload"`
}
