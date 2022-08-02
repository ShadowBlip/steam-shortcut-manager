package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
)

type FatalError struct {
	err error
}

func (e *FatalError) Error() string {
	return e.err.Error()
}

type MultiError struct {
	Errors []error
	foo    *multierror.Error
}

func (m *MultiError) Error() string {
	report := make([]string, 0, len(m.Errors)+1)
	report = append(report, fmt.Sprintf("%d errors occurred", len(m.Errors)))
	for _, err := range m.Errors {
		report = append(report, err.Error())
	}
	return strings.Join(report, "; ")
}

// ExitError will print an error and exit depending on the output format
func ExitError(err error, format string) {
	switch format {
	case "json":
		out, _ := json.Marshal(map[string]string{"errors": err.Error()})
		fmt.Println(string(out))
		os.Exit(1)
	default:
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// Print debug messages if debug is enabled
func DebugPrintln(s ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Println(s...)
	}
}
