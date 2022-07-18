package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/nelsonjchen/google-cloud-nuke/v1/resources"
)

var (
	ReasonSkip        = *color.New(color.FgYellow)
	ReasonError       = *color.New(color.FgRed)
	ReasonWaitPending = *color.New(color.FgBlue)
	ReasonSuccess     = *color.New(color.FgGreen)
)

var (
	ColorResourceType       = *color.New()
	ColorResourceID         = *color.New(color.Bold)
	ColorResourceProperties = *color.New(color.Italic)
)

// Format the resource properties in sorted order ready for printing.
// This ensures that multiple runs of google-cloud-nuke produce stable output so
// that they can be compared with each other.
func Sorted(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sorted := make([]string, 0, len(m))
	for k := range keys {
		sorted = append(sorted, fmt.Sprintf("%s: \"%s\"", keys[k], m[keys[k]]))
	}
	return fmt.Sprintf("[%s]", strings.Join(sorted, ", "))
}

func Log(resourceType string, r resources.Resource, c color.Color, msg string) {
	_, _ = ColorResourceType.Print(resourceType)
	fmt.Printf(" - ")

	rString, ok := r.(resources.LegacyStringer)
	if ok {
		_, _ = ColorResourceID.Print(rString.String())
		fmt.Printf(" - ")
	}

	rProp, ok := r.(resources.ResourcePropertyGetter)
	if ok {
		_, _ = ColorResourceProperties.Print(Sorted(rProp.Properties()))
		fmt.Printf(" - ")
	}

	_, _ = c.Printf("%s\n", msg)
}
