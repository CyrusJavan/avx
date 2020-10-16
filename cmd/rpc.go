package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/CyrusJavan/avx/color"
	"github.com/spf13/cobra"
)

var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "Make an API call to the controller.",
	Args:  cobra.MinimumNArgs(1),
	RunE:  rpcFunc,
}

func rpcFunc(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("could not get client: %w", err)
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
		return fmt.Errorf(color.Sprint("marshalling json data: %v", color.Red), err)
	}

	err = json.Indent(&dataBuffer, jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf(color.Sprint("indenting json data: %v", color.Red), err)
	}

	fmt.Printf("controller IP: %s\n", client.ControllerIP)
	fmt.Printf("request body:\n"+color.Sprint("%s\n", color.Green), dataBuffer.String())

	start := time.Now()
	_, b, err := client.Do("POST", data)
	end := time.Now()
	fmt.Printf("latency: %dms\n", end.Sub(start).Milliseconds())
	if err != nil {
		return fmt.Errorf(color.Sprint("non-nil error from API: %v", color.Red), err)
	}

	var pp bytes.Buffer
	err = json.Indent(&pp, b, "", "  ")

	fmt.Printf("response body:\n%s\n", color.Sprint(pp.String(), color.Green))

	return nil
}