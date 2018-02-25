package commandlineargs

import (
	"errors"
	"fmt"
)

//ValidateArguments method to validate all required command line args are specified
func ValidateArguments(args ...string) error {
	if args == nil {
		return errors.New("No command line args were specified")
	}
	for _, arg := range args {
		println("arg: ", arg)
		if arg == "" {
			fmt.Printf("argument %s missing", arg)
			return errors.New("Unspecified required command line args")
		}
	}
	return nil
}
