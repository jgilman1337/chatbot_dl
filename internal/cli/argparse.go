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
	}

	//Perform primary validation via go-playground/validator
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&opts); err != nil {
		return nil, err
	}

	//Perform secondary validation against URLs; maybe
	/*
		okUrls := make([]string, len(opts.Positional.URLs))
		for i, url := range opts.Positional.URLs {

		}
	*/

	return &opts, nil
}
