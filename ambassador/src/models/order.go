package models

type Order struct {
	Id     uint    `json:"id" gorm:"primaryKey; autoIncrement:false"`
	UserId uint    `json:"user_id"`
	Code   string  `json:"code"`
	Total  float64 `json:"total"`
}
