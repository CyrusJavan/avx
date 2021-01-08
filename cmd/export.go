package cmd

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/CyrusJavan/avx/color"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export TF config.",
	Args:  cobra.ExactArgs(2),
	RunE:  exportFunc,
}

func exportFunc(cmd *cobra.Command, args []string) error {
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

	var dataBuffer bytes.Buffer
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf(color.Sprint("marshalling json data: %v", color.Red), err)
	}

	err = json.Indent(&dataBuffer, jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf(color.Sprint("indenting json data: %v", color.Red), err)
	}
	if !JsonOnly {
		fmt.Printf("controller IP: %s\n", client.ControllerIP)
		fmt.Printf("request body:\n"+color.Sprint("%s\n", color.Green), dataBuffer.String())
	}

	start := time.Now()
	_, b, err := client.Do("POST", data)
	end := time.Now()
	if !JsonOnly {
		fmt.Printf("latency: %dms\n", end.Sub(start).Milliseconds())
	}
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

	for _, zipFile := range zipReader.File {
		unzippedFileBytes, err := readZipFile(zipFile)
		if err != nil {
			return fmt.Errorf("reading zipped file: %v", err)
		}
		tfFilePath := args[1]
		if !strings.HasSuffix(tfFilePath, "/") {
			tfFilePath += "/"
		}
		err = ioutil.WriteFile(tfFilePath+zipFile.Name, unzippedFileBytes, 0644)
		if err != nil {
			return fmt.Errorf("writing file %s: %v", zipFile.Name, err)
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
