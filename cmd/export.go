package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/CyrusJavan/avx/color"
	"github.com/spf13/cobra"
)

var WriteToFile bool
var IncludeShellFile bool
var ManageInternally bool

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export TF config.",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  exportFunc,
}

func exportFunc(cmd *cobra.Command, args []string) error {
	if len(args) == 2 && !WriteToFile {
		return fmt.Errorf("if flag (--file, -f) is not set then a path must not be provided as the second argument")
	}

	client, err := getClient()
	if err != nil {
		return jsonErr("could not get client", err)
	}

	resourceName := args[0]

	manageInternally := "false"
	if ManageInternally {
		manageInternally = "true"
	}

	data := map[string]string{
		"action":                        "export_terraform_resource",
		"CID":                           client.CID,
		"manage_attachments_internally": manageInternally,
		"resource":                      resourceName,
	}
	v := url.Values{}
	for k, s := range data {
		v.Add(k, s)
	}
	u := "https://" + client.ControllerIP + "/v1/api/"
	resp, err := client.HTTPClient.PostForm(u, v)
	if err != nil {
		return fmt.Errorf(color.Sprint("non-nil error from API: %v", color.Red), err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return fmt.Errorf("init zip reader: %v", err)
	}

	for i, zipFile := range zipReader.File {
		if !IncludeShellFile && i > 0 {
			continue
		}
		unzippedFileBytes, err := readZipFile(zipFile)
		if err != nil {
			return fmt.Errorf("reading zipped file: %v", err)
		}
		if !WriteToFile {
			_, _ = fmt.Fprint(cmd.OutOrStdout(), string(unzippedFileBytes))
		} else {
			tfFilePath := args[1]
			if !strings.HasSuffix(tfFilePath, "/") {
				tfFilePath += "/"
			}
			err = ioutil.WriteFile(tfFilePath+zipFile.Name, unzippedFileBytes, 0644)
			if err != nil {
				return fmt.Errorf("writing file %s: %v", zipFile.Name, err)
			}
		}
	}

	return nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
