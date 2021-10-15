package services

import (
	"fmt"
	"github.com/antoniopapa/go-user-service"
	"os"
)

var UserService services.Service

func Setup() {
	UserService = services.CreateService(fmt.Sprintf("%s/api/", os.Getenv("USERS_MS")))
}
