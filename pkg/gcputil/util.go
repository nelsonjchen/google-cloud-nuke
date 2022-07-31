package gcputil

import (
	"fmt"
	"google.golang.org/api/compute/v1"
	"path"
	"regexp"
)

var (
	RESecretHeader = regexp.MustCompile(`(?m:^([^:]*(Auth|Security)[^:]*):.*$)`)
)

func HideSecureHeaders(dump []byte) []byte {
	return RESecretHeader.ReplaceAll(dump, []byte("$1: <hidden>"))
}

// ComputeRemoveWaiter is a waiter for compute resources
// It will return an operation if the operation is
func ComputeRemoveWaiter(op *compute.Operation, service *compute.Service, project string) (*compute.Operation, error) {
	if op.Status == "RUNNING" {
		return op, fmt.Errorf("operation still running")
	}
	// Refresh the operation
	if op.Zone != "" {
		call := service.ZoneOperations.Get(project, path.Base(op.Zone), op.Name)
		resp, err := call.Do()
		if err != nil {
			return nil, err
		}
		op = resp
	} else if op.Region != "" {
		call := service.RegionOperations.Get(project, path.Base(op.Region), op.Name)
		resp, err := call.Do()
		if err != nil {
			return nil, err
		}
		op = resp
	} else {
		call := service.GlobalOperations.Get(project, op.Name)
		resp, err := call.Do()
		if err != nil {
			return nil, err
		}
		op = resp
	}

	if op.Error != nil {
		if op.HTTPStatusCode == 404 {
			// It's gone, that's OK.
			return op, nil
		}
		return nil, fmt.Errorf("operation error: %s", op.Error.Errors[0].Message)
	}

	return op, nil
}
