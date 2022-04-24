package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello!")

	// Read CLI flags.
	configPath := flag.String("config-path", "", "path to RSS feeds")
	sendTo := flag.String("send-to", "", "email to send to")
	// Read environment variables.
	sendgridKey := os.Getenv("SENDGRID_KEY")
}
