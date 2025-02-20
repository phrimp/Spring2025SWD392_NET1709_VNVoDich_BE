package utils

import (
	"os"
)

var API_KEY string

func init() {
	API_KEY = os.Getenv("API_KEY")
}
