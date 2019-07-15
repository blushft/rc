package main

import (
	"log"

	"github.com/blushft/rc/rc/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Error: ", err)
	}
}
