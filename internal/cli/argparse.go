package cli

import (
	"github.com/go-playground/validator/v10"
	"github.com/jessevdk/go-flags"
)

// Parses arguments for the commandline interface.
func ParseArgs(argv []string) (*Options, error) {
	//Parse args from the user
	opts := Options{}

	_, err := flags.ParseArgs(&opts, argv)
	if err != nil {
		//Do not error out on help errors
		if e, ok := err.(*flags.Error); ok {
			if e.Type == flags.ErrHelp {
				return nil, nil //Options should not be used in this state
			}
		}

		return nil, err
	}

	//Perform primary validation via go-playground/validator
	validate := validator.New(validator.WithRequiredStructEnabled())
	registerCustomHandlers(validate)
	if err := validate.Struct(&opts); err != nil {
		return nil, handleValidationErr(err)
	}

	return &opts, nil
}
