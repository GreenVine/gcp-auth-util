package cli

import "strings"

type stringSlice []string

func (f *stringSlice) String() string {
	return strings.Join(*f, ",")
}

func (f *stringSlice) Set(val string) error {
	*f = append(*f, val)
	return nil
}
