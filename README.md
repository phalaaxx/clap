# clap

CLAP (or Command Line Arguments Parser) is a wrapper around flag package. Its purpose is to make it a bit easier to
parse arguments with short and long options.

Currently clap supports only String, Int and Bool options.

## Usage

Options are provided with a single function call with the following parameters:

* Short name    (rune)
* Long name     (string)
* Default value (varies)
* Help string   (string)
* Required flag (bool)

    nameOpt := clap.String('n', "name", "", "Name of someone", false)
    ageOpt := clap.Int('a', "age", 20, "Age of someone [default: 20]", false)
    realOpt := clap.Bool('r', "real", false, "Real person", false)
    longOpt := clap.Int(0, "long-option", 15, "Long-option only [default: 15], false)

If an option is required, but is not provided on the command line, an error message is displayed and the program
will exit with code 255. Currently only string options work like this.

Default values are not automatically provided in the help text.
