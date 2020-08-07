package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-aviatrix/aviatrix"
)

func main() {
	f, _ := os.Open("/dev/null")
	log.SetOutput(f)

	if err := run(os.Args[1:]); err != nil {
		fmt.Printf(color("%v\n", Red), err)
		os.Exit(1)
	}
}

func checkEnvVars() error {
	mustBeSet := []string{
		"AVIATRIX_CONTROLLER_IP",
		"AVIATRIX_USERNAME",
		"AVIATRIX_PASSWORD",
	}

	var notFound []string

	for _, v := range mustBeSet {
		if os.Getenv(v) == "" {
			notFound = append(notFound, v)
		}
	}

	if len(notFound) > 0 {
		return fmt.Errorf("environment variables %v must be set", notFound)
	}

	return nil
}

func run(args []string) error {
	if err := checkEnvVars(); err != nil {
		return err
	}
	cfg := aviatrix.Config{
		Username:     os.Getenv("AVIATRIX_USERNAME"),
		Password:     os.Getenv("AVIATRIX_PASSWORD"),
		ControllerIP: os.Getenv("AVIATRIX_CONTROLLER_IP"),
	}

	client, err := cfg.Client()
	if err != nil {
		return fmt.Errorf("could not get client from config: %w", err)
	}

	// If no arg given, assuming user just wanted the CID. Exit now.
	if len(args) == 0 {
		fmt.Printf("CID: "+color("%q\n", Green), client.CID)
		return nil
	}

	action := args[0]

	data := map[string]interface{}{
		"action": action,
		"CID":    client.CID,
	}

	for _, v := range args[1:] {
		parts := strings.Split(v, "=")
		if len(parts) != 2 {
			return fmt.Errorf("invalid format for API params, expected 'key=value', got %q", v)
		}
		data[parts[0]] = parts[1]
	}

	var dataBuffer bytes.Buffer
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf(color("marshalling json data: %v", Red), err)
	}

	err = json.Indent(&dataBuffer, jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf(color("indenting json data: %v", Red), err)
	}

	fmt.Printf("controller IP: %s\n", client.ControllerIP)
	fmt.Printf("request body:\n"+color("%s\n", Green), dataBuffer.String())

	start := time.Now()
	_, b, err := client.Do("POST", data)
	end := time.Now()
	fmt.Printf("latency: %dms\n", end.Sub(start).Milliseconds())
	if err != nil {
		return fmt.Errorf(color("non-nil error from API: %v", Red), err)
	}

	var pp bytes.Buffer
	err = json.Indent(&pp, b, "", "  ")

	fmt.Printf("response body:\n%s\n", color(pp.String(), Green))

	return nil
}

type Color string

const (
	Reset Color = "\033[0m"
	Red   Color = "\033[31m"
	Green Color = "\033[32m"
)

func color(s string, c Color) string {
	return string(c) + s + string(Reset)
}
