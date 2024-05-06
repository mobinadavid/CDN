package main

import (
	"cdn/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
