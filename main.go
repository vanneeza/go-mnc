package main

import "github.com/vanneeza/go-mnc/utils/server"

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
