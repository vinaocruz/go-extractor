/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/vinaocruz/go-extractor/cmd"
)

func main() {
	start := time.Now()

	loadEnv()
	cmd.Execute()

	elapsed := time.Since(start)
	log.Printf("Executed in %s", elapsed)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
