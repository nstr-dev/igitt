package main

import (
	"os"

	"github.com/noahstreller/igitt/internal/utilities/logger"
)

func main() {
	arugments := os.Args

	for i := 0; i < len(arugments); i++ {
		print(arugments[i] + " ")
	}
	logger.InfoLogger.Println("Hello")
}