package cli

import "strings"

type stringSlice []string // custom struct to accept repeated arguments in `flag`

func (f *stringSlice) String() string {
	return strings.Join(*f, ",")
}

func (f *stringSlice) Set(val string) error {
	*f = append(*f, val)
	return nil
}
