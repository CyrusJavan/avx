package cmd

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/CyrusJavan/avx/color"
	"github.com/spf13/cobra"
)

var WriteToStdOut bool
var IncludeShellFile bool

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export TF config.",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  exportFunc,
}

func exportFunc(cmd *cobra.Command, args []string) error {
	if len(args) == 1 && !WriteToStdOut {
		return fmt.Errorf("if flag (--use-stdout, -s) is not set then a path must be provided as the second argument")
	}

	client, err := getClient()
	if err != nil {
		return jsonErr("could not get client", err)
	}

	resourceName := args[0]

	data := map[string]interface{}{
		"action":                        "run_export_terraform",
		"CID":                           client.CID,
		"operation":                     "export_tf",
		"manage_attachments_internally": "false",
		"resource":                      resourceName,
	}

	_, b, err := client.Do("POST", data)
	if err != nil {
		return fmt.Errorf(color.Sprint("non-nil error from API: %v", color.Red), err)
	}

	type Resp struct {
		Results string
	}
	var r Resp
	err = json.Unmarshal(b, &r)
	if err != nil {
		return fmt.Errorf("decoding response: %v", err)
	}

	path := "https://" + client.ControllerIP + "/v1/download?" +
		fmt.Sprintf("filename=%s", r.Results) +
		"&" + fmt.Sprintf("CID=%s", client.CID)
	resp, err := client.Request("GET", path, nil)
	if err != nil {
		return fmt.Errorf("making download request: %v", err)
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
		if WriteToStdOut {
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
