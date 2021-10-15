package models

type KafkaError struct {
	Id    uint
	Key   []byte
	Value []byte
	Error error `gorm:"type:text"`
}
