package pkg

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/vanneeza/go-mnc/utils/helper"
)

func GetEnv(key string) string {
	err := godotenv.Load(".env")
	helper.PanicError(err)
	return os.Getenv(key)
}
