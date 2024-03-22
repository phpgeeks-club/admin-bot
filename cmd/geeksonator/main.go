package main

import (
	"fmt"

	"geeksonator/internal/app/geeksonator"
)

func main() {
	if err := geeksonator.Start(); err != nil {
		panic(fmt.Errorf("geeksonator.Start: %v", err))
	}
}
