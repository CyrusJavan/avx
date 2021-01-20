package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/terraform-providers/terraform-provider-aviatrix/aviatrix"
	"github.com/terraform-providers/terraform-provider-aviatrix/goaviatrix"
)

var (
	JsonOnly bool

	rootCmd = &cobra.Command{
		Use:   "avx",
		Short: "Aviatrix CLI",
		Long:  "Avx is an Aviatrix API CLI tool.",
		Args:  cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := checkEnvVars()
			if err != nil {
				return err
			}
			return nil
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rpcCmd.PersistentFlags().BoolVarP(&JsonOnly, "json-only", "j", false, "json response only output")
	rootCmd.AddCommand(rpcCmd)
	rootCmd.AddCommand(loginCmd)
	exportCmd.PersistentFlags().BoolVarP(&WriteToFile, "file", "f", false, "write output to file instead of stdout")
	exportCmd.PersistentFlags().BoolVarP(&IncludeShellFile, "include-shell-script", "i", false, "also output the import shell script")
	exportCmd.PersistentFlags().BoolVarP(&ManageInternally, "manage-attm-internally", "m", false, "export with attachments managed internally")
	rootCmd.AddCommand(exportCmd)
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

func getClient() (*goaviatrix.Client, error) {
	cfg := aviatrix.Config{
		Username:     os.Getenv("AVIATRIX_USERNAME"),
		Password:     os.Getenv("AVIATRIX_PASSWORD"),
		ControllerIP: os.Getenv("AVIATRIX_CONTROLLER_IP"),
	}

	client, err := cfg.Client()
	if err != nil {
		return nil, fmt.Errorf("could not get client from config: %w", err)
	}
	return client, nil
}
