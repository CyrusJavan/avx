package cmd

import (
	"fmt"

	"github.com/CyrusJavan/avx/color"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Get a CID from the aviatrix API.",
	Args:  cobra.NoArgs,
	RunE:  loginFunc,
}

func loginFunc(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return fmt.Errorf("could not get client: %w", err)
	}

	fmt.Println("Login successful")
	fmt.Printf("CID: "+color.Sprint("%q\n", color.Green), client.CID)
	return nil
}
