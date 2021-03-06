package goaviatrix

import (
	"strings"

	"github.com/pkg/errors"
)

func (c *Client) ControllerVersionValidation(supportedVersion string) error {
	suppVersion := strings.Split(supportedVersion, ".")

	currentVersion, _, err := c.GetCurrentVersion()
	if err != nil {
		return err
	}
	currVersion := strings.Split(currentVersion, ".")
	if suppVersion[0] != currVersion[0] || suppVersion[1] != currVersion[1] {
		return errors.New("current Terraform branch supports controller version: UserConnect-" + supportedVersion +
			". Please upgrade/downgrade controller or change Terraform branch.")
	}

	return nil
}
