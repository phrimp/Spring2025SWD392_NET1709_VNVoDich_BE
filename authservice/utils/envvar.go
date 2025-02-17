package utils

import (
	"os"
)

type ServicesRoute struct {
	UserService string
}

var (
	API_KEY         string
	SERVICES_ROUTES ServicesRoute
)

func init() {
	API_KEY = os.Getenv("API_KEY")
	SERVICES_ROUTES.UserService = "http://user-service:" + os.Getenv("USER_SERVICE_PORT")
}
