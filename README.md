# clap

CLAP (or Command Line Arguments Parser) is a wrapper around flag package. Its purpose is to make it easier to parse
arguments with short and long options.

Currently clap supports most of the types available in flag package: string, int, int64, uint, uint64, float64, bool.


## Usage

Options are provided with a single function call with the following parameters:

* Short name    (rune)
* Long name     (string)
* Default value (varies)
* Help string   (string)
* Required flag (bool)

```
nameOpt := clap.String('n', "name", "", "Name of someone", false)
ageOpt := clap.Int('a', "age", 20, "Age of someone [default: 20]", false)
realOpt := clap.Bool('r', "real", false, "Real person", false)
longOpt := clap.Int(0, "long-option", 15, "Long-option only [default: 15], false)
```

Since short options are runes, an empty short option can be specified with a 0 (see the example).

If an option is required, but is not provided on the command line, an error message is displayed and the program
will exit with status code 255. Currently only string options work like this.

Default values are not automatically provided in the help text.

### Var arguments

Similarly to flag package, a Var argument can be specified for more complex data types (like string list arguments).
For this to work, clap.Var needs a flag.Value interface as its first function argument.

### Duration arguments

CLI arguments with Duration type are not yet supported.

## Parsing

The Parse function is used to parse command line options. Difference between clap.Parse and flag.Parse is that in
clap.Parse accepts one argument - required (bool). When set to true, at least one command line option needs to be
provided, otherwise an error message will be displayed and the program will exit with status code 255.
```
clap.Parse(true)
```
