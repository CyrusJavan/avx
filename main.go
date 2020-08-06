package main

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-aviatrix/aviatrix"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	cfg := aviatrix.Config{
		Username: "admin",
		Password: "Pass_word1",
		ControllerIP: "13.57.122.42",
	}

	client, err := cfg.Client()
	if err != nil {
		return fmt.Errorf("could not get client from config: %w", err)
	}

	fmt.Printf("CID: %q\n", client.CID)

	return nil
}
