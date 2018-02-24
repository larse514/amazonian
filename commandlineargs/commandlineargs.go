package commandlineargs

import (
	"errors"
)

//ValidateArguments method to validate all required command line args are specified
func ValidateArguments(args ...string) error {
	if args == nil {
		return errors.New("No command line args were specified")
	}
	for _, arg := range args {
		if arg == "" {
			return errors.New("Unspecified required command line args")
		}
	}
	return nil
}
