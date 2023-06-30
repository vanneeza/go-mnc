package pkg

import (
	"math/rand"
	"time"
)

func GenerateRandomNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10000000000)
}
