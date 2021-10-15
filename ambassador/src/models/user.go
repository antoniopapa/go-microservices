package models

import (
	"gorm.io/gorm"
)

type User struct {
	Id           uint    `json:"id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Email        string  `json:"email"`
	Password     []byte  `json:"-"`
	IsAmbassador bool    `json:"is_ambassador"`
	Revenue      float64 `json:"revenue"`
}

func (user *User) Name() string {
	return user.FirstName + " " + user.LastName
}

type Ambassador User

func (ambassador *Ambassador) CalculateRevenue(db *gorm.DB) {
	var orders []Order

	db.Find(&orders, &Order{
		UserId: ambassador.Id,
	})

	var revenue = 0.0

	for _, order := range orders {
		revenue += order.Total
	}

	ambassador.Revenue = revenue
}
