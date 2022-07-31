package gcputil

import (
	"fmt"
	"google.golang.org/api/compute/v1"
	"path"
	"regexp"
	"time"
)

var (
	RESecretHeader = regexp.MustCompile(`(?m:^([^:]*(Auth|Security)[^:]*):.*$)`)
)

func HideSecureHeaders(dump []byte) []byte {
	return RESecretHeader.ReplaceAll(dump, []byte("$1: <hidden>"))
}

// ComputeRemoveWaiter is a waiter for compute resources
// It is used to wait for compute resource operation to be removed and will return an updated operation
func ComputeRemoveWaiter(op *compute.Operation, service *compute.Service, project string) (*compute.Operation, error) {
	if op.HTTPStatusCode == 404 {
		return op, nil
	}

	runningCount := 0
	runningCountLimit := 2
	for {
		if op.Status == "DONE" {
			break
		}

		if op.Status == "RUNNING" {
			runningCount++
		}

		if runningCount > runningCountLimit {
			return nil, fmt.Errorf("operation %s is still running. will try operation again", op.Name)
		}

		time.Sleep(1 * time.Second)

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
	}

	if op.Error != nil {
		return nil, fmt.Errorf("operation error: %s", op.Error.Errors[0].Message)
	}

	return op, nil
}

func Base(s string) string {
	if s == "" {
		return ""
	}
	return path.Base("")
}
