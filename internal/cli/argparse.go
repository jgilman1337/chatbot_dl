package cli

// A list of valid arguments for the argument parser.
type args struct {
	Positional posargs `positional-args:"true" required:"true"`
}

// Holds positional arguments.
type posargs struct {
	URLs []string `positional-arg-name:"URLs"`
}

// Parses arguments for the commandline interface.
func ParseArgs(argv []string) (*Options, error) {

	return nil, nil
}
