# clap

CLAP (or Command Line Arguments Parser) is a wrapper around flag package. Its purpose is to make it easier to parse
arguments with short and long options.

Currently clap supports most of the types available in flag package: string, int, int64, uint, uint64, float64, bool.

## Installation

```
go get github.com/phalaaxx/clap
```

## Example

```go
package main

import (
	"fmt"

	"github.com/phalaaxx/clap"
)

func main() {
	name := clap.String('n', "name", "", "Name of someone", true)
	age := clap.Int('a', "age", 20, "Age of someone [default: 20]", false)
	verbose := clap.Bool('v', "verbose", false, "Verbose output", false)
	clap.Parse(false)

	fmt.Println(*name, *age, *verbose)
	for idx := 0; idx < clap.NArg(); idx++ {
		fmt.Println("positional:", clap.Arg(idx))
	}
}
```

Running the program with `--help` shows:

```
Usage: example [OPTIONS] --name <NAME>

Options:
  -n, --name <NAME>  Name of someone
  -a, --age <AGE>    Age of someone [default: 20]
  -v, --verbose      Verbose output
```

## Usage

Options are provided with a single function call with the following parameters:

* Short name    (rune)
* Long name     (string)
* Default value (varies)
* Help string   (string)
* Required flag (bool)

```go
nameOpt := clap.String('n', "name", "", "Name of someone", false)
ageOpt := clap.Int('a', "age", 20, "Age of someone [default: 20]", false)
realOpt := clap.Bool('r', "real", false, "Real person", false)
longOpt := clap.Int(0, "long-option", 15, "Long-option only [default: 15]", false)
```

Since short options are runes, an empty short option can be specified with a 0 (see the example). Similarly, an empty
string can be used for an option without a long name.

The functions above return a pointer to the option value. Similarly to flag package, there is also a `*Var` variant of
each function, which stores the value in an existing variable instead:

```go
var name string
clap.StringVar(&name, 'n', "name", "", "Name of someone", false)
```

Available functions: `String`, `Int`, `Int64`, `Uint`, `Uint64`, `Float64`, `Bool` and their `StringVar`, `IntVar`,
`Int64Var`, `UintVar`, `Uint64Var`, `Float64Var`, `BoolVar` counterparts.

Long options must be prefixed with a double dash (`--name`) and short options with a single dash (`-n`), otherwise
an error message is displayed. Values can be provided either as a separate argument or after `=` (`--name=value`).

If an option is marked as required, but is not provided on the command line, an error message listing the missing
options is displayed and the program will exit with status code 255. The required flag has no effect on bool options.

Default values are not automatically provided in the help text.

### Var arguments

Similarly to flag package, a Var argument can be specified for more complex data types (like string list arguments).
For this to work, clap.Var needs a flag.Value interface as its first function argument. Since the value is already
provided by the flag.Value, clap.Var has no default value parameter:

```go
type stringList []string

func (s *stringList) String() string     { return strings.Join(*s, ",") }
func (s *stringList) Set(v string) error { *s = append(*s, v); return nil }

var tags stringList
clap.Var(&tags, 't', "tag", "Tag (can be repeated)", false)
```

If the value implements `IsBoolFlag() bool` returning true (as in flag package), the option is treated as a bool option
and does not require a value.

### Duration arguments

CLI arguments with Duration type are not yet supported.

## Parsing

The Parse function is used to parse command line options. The difference between clap.Parse and flag.Parse is that
clap.Parse accepts one argument - required (bool). When set to true, at least one command line argument needs to be
provided, otherwise an error message will be displayed and the program will exit with status code 255.

```go
clap.Parse(true)
```

As in flag package, parsing stops at the first non-option argument or after the `--` terminator. The remaining
arguments are available with `clap.NArg()` and `clap.Arg(i)`.

## Output

Help and error messages are printed to stderr, same as in flag package. When stderr is a terminal, the output is
colored, unless the `NO_COLOR` environment variable is set or `TERM` is set to `dumb`.

Exit status codes:

* `0` - help was requested with `-h` or `--help`
* `2` - unknown option or invalid option value (reported by flag package)
* `255` - no arguments, missing required options or an option used with the wrong number of dashes

## Limitations

clap uses the global flag.CommandLine flag set, so subcommands (separate flag sets) are not supported.
