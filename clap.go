/*
BSD 2-Clause License

# Copyright (c) 2025, Bozhin Zafirov

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice, this
    list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package clap

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"unicode/utf8"
)

/* argumentHelp stores command line arguments info */
type argumentHelp struct {
	HasValue  bool
	ShortName rune
	LongName  string
	HelpText  string
	Required  bool
	Value     interface{}
}

/* global variables */
var (
	argHelp    []argumentHelp
	isTerminal bool
)

/* paint wraps text in the specified ANSI escape sequence when output is a terminal */
func paint(code string, text string) string {
	if !isTerminal {
		return text
	}
	return fmt.Sprintf("\033[%sm%s\033[0m", code, text)
}

/* ANSI styles */
const (
	styleBold      = "1"
	styleUnderline = "1;4"
	styleError     = "1;31"
	styleRequired  = "32"
)

/* placeholder returns the value placeholder of the argument */
func (a argumentHelp) placeholder() string {
	if len(a.LongName) == 0 {
		return "<VALUE>"
	}
	return fmt.Sprintf("<%s>", strings.ToUpper(a.LongName))
}

/* flagName returns the preferred (long if available) option name with dashes */
func (a argumentHelp) flagName() string {
	if len(a.LongName) != 0 {
		return "--" + a.LongName
	}
	return "-" + string(a.ShortName)
}

/* usage returns option name followed by value placeholder, e.g. "--name <NAME>" */
func (a argumentHelp) usage(style string) string {
	if !a.HasValue {
		return paint(style, a.flagName())
	}
	return fmt.Sprintf("%s %s", paint(style, a.flagName()), a.placeholder())
}

/* option returns the plain and colored representations of the option column */
func (a argumentHelp) option() (plain string, colored string) {
	var names []string
	var coloredNames []string
	if a.ShortName != 0 {
		names = append(names, "-"+string(a.ShortName))
	} else {
		/* align long-only options with the long names of other options */
		names = append(names, "  ")
	}
	if len(a.LongName) != 0 {
		names = append(names, "--"+a.LongName)
	}
	for _, name := range names {
		if strings.TrimSpace(name) == "" {
			coloredNames = append(coloredNames, name)
		} else {
			coloredNames = append(coloredNames, paint(styleBold, name))
		}
	}
	separator := ", "
	if a.ShortName == 0 {
		separator = "  "
	}
	plain = "  " + strings.Join(names, separator)
	colored = "  " + strings.Join(coloredNames, separator)
	if a.HasValue {
		plain = fmt.Sprintf("%s %s", plain, a.placeholder())
		colored = fmt.Sprintf("%s %s", colored, a.placeholder())
	}
	return
}

/* String representation of the argument help, padded to maxLen */
func (a argumentHelp) String(maxLen int) string {
	plain, colored := a.option()
	padding := maxLen - utf8.RuneCountInString(plain) + 2
	return fmt.Sprintf("%s%s%s\n", colored, strings.Repeat(" ", padding), a.HelpText)
}

/* isBoolFlag returns true if the flag value does not need an argument */
func isBoolFlag(value flag.Value) bool {
	boolFlag, ok := value.(interface{ IsBoolFlag() bool })
	return ok && boolFlag.IsBoolFlag()
}

/* genericAddVar is a wrapper around flag functions to add command line arguments */
func genericAddVar[T any](data *T, name string, initial T, usage string) {
	switch any(*data).(type) {
	case int:
		flag.IntVar(any(data).(*int), name, any(initial).(int), usage)
	case int64:
		flag.Int64Var(any(data).(*int64), name, any(initial).(int64), usage)
	case float64:
		flag.Float64Var(any(data).(*float64), name, any(initial).(float64), usage)
	case string:
		flag.StringVar(any(data).(*string), name, any(initial).(string), usage)
	case uint:
		flag.UintVar(any(data).(*uint), name, any(initial).(uint), usage)
	case uint64:
		flag.Uint64Var(any(data).(*uint64), name, any(initial).(uint64), usage)
	case bool:
		flag.BoolVar(any(data).(*bool), name, any(initial).(bool), usage)
	}
}

/* genericVar defines a T flag in the argument result */
func genericVar[T any](result *T, short rune, long string, value T, usage string, required bool) {
	/* add generic options */
	if short != 0 {
		genericAddVar[T](result, string(short), value, usage)
	}
	if len(long) != 0 {
		genericAddVar[T](result, long, value, usage)
	}
	_, isBool := any(value).(bool)
	/* update help data */
	argHelp = append(
		argHelp,
		argumentHelp{
			HasValue:  !isBool,
			ShortName: short,
			LongName:  long,
			HelpText:  usage,
			Required:  required,
			Value:     value,
		},
	)
}

/* generic defines a T flag and returns a pointer to it */
func generic[T any](short rune, long string, value T, usage string, required bool) *T {
	result := new(T)
	genericVar[T](result, short, long, value, usage, required)
	return result
}

/* define global clap functions */
var (
	/* define functions assigning pointer to a flag */
	StringVar  = genericVar[string]
	IntVar     = genericVar[int]
	Int64Var   = genericVar[int64]
	Float64Var = genericVar[float64]
	UintVar    = genericVar[uint]
	Uint64Var  = genericVar[uint64]
	BoolVar    = genericVar[bool]
	/* define functions returning pointer to a flag */
	String  = generic[string]
	Int     = generic[int]
	Int64   = generic[int64]
	Float64 = generic[float64]
	Uint    = generic[uint]
	Uint64  = generic[uint64]
	Bool    = generic[bool]
	/* direct flag mappings */
	NArg = flag.NArg
	Arg  = flag.Arg
)

/* Var is a special-case function / wrapper around flag.Var for custom data type flags */
func Var(value flag.Value, shortName rune, longName string, helpText string, required bool) {
	if shortName != 0 {
		flag.Var(value, string(shortName), helpText)
	}
	if len(longName) != 0 {
		flag.Var(value, longName, helpText)
	}
	argHelp = append(
		argHelp,
		argumentHelp{
			HasValue:  !isBoolFlag(value),
			ShortName: shortName,
			LongName:  longName,
			HelpText:  helpText,
			Required:  required,
			Value:     value,
		},
	)
}

/* helpTip returns the hint printed at the end of error messages */
func helpTip() string {
	return fmt.Sprintf("For more information, try '%s'.\n", paint(styleBold, "--help"))
}

/* usageLine returns program usage with all required arguments */
func usageLine(options bool) string {
	result := fmt.Sprintf("%s %s", paint(styleUnderline, "Usage:"), paint(styleBold, path.Base(os.Args[0])))
	if options {
		result += " [OPTIONS]"
	}
	for _, arg := range argHelp {
		if arg.Required {
			result = fmt.Sprintf("%s %s", result, arg.usage(styleBold))
		}
	}
	return result
}

/* exitError prints an error message followed by help tip and exits */
func exitError(message string) {
	fmt.Fprintf(flag.CommandLine.Output(), "%s %s\n\n%s", paint(styleError, "error:"), message, helpTip())
	os.Exit(-1)
}

/* errorHelp renders error message for the missing required arguments */
func errorHelp(missing []argumentHelp) string {
	var result strings.Builder
	fmt.Fprintf(&result, "%s the following arguments are not provided:\n", paint(styleError, "error:"))
	for _, arg := range missing {
		fmt.Fprintf(&result, "  %s\n", arg.usage(styleRequired))
	}
	fmt.Fprintf(&result, "\n%s\n\n%s", usageLine(false), helpTip())
	return result.String()
}

/* usageHeader prints flags usage header */
func usageHeader() string {
	return fmt.Sprintf("%s\n\n%s", usageLine(true), paint(styleUnderline, "Options:"))
}

/* ErrNoArg prints error and exits when no arguments are provided while at least one is required */
func ErrNoArg() {
	exitError("no arguments are provided")
}

/* lookupArg returns the argument with the specified short or long name */
func lookupArg(name string) *argumentHelp {
	for idx := range argHelp {
		arg := &argHelp[idx]
		if name == arg.LongName || (arg.ShortName != 0 && name == string(arg.ShortName)) {
			return arg
		}
	}
	return nil
}

/* checkDashes makes sure long options are used with "--" and short options with "-" */
func checkDashes(args []string) {
	for idx := 0; idx < len(args); idx++ {
		/* flag parsing stops at the first non-flag argument or at the "--" terminator */
		arg := args[idx]
		if len(arg) < 2 || arg[0] != '-' || arg == "--" {
			return
		}
		name, dashes := arg[1:], 1
		if name[0] == '-' {
			name, dashes = name[1:], 2
		}
		name, _, hasValue := strings.Cut(name, "=")
		option := lookupArg(name)
		if option == nil {
			/* unknown options are reported by flag.Parse */
			continue
		}
		isShort := option.ShortName != 0 && name == string(option.ShortName)
		if isShort && dashes != 1 {
			exitError(fmt.Sprintf("unexpected argument '%s', did you mean '-%s'?", arg, name))
		}
		if !isShort && dashes != 2 {
			exitError(fmt.Sprintf("unexpected argument '%s', did you mean '--%s'?", arg, name))
		}
		/* skip the option value */
		if option.HasValue && !hasValue {
			idx++
		}
	}
}

/* Parse is a wrapper around flag.Parse function */
func Parse(required bool) {
	/* print error on empty arguments list when at least one argument is required */
	if required && len(os.Args) == 1 {
		ErrNoArg()
	}
	/* parse and check if required arguments are provided */
	checkDashes(os.Args[1:])
	flag.Parse()
	provided := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		provided[f.Name] = true
	})
	var missing []argumentHelp
	for _, arg := range argHelp {
		/* do not check non-mandatory and bool arguments */
		if !arg.Required || !arg.HasValue {
			continue
		}
		/* make sure value is provided */
		if !provided[arg.LongName] && (arg.ShortName == 0 || !provided[string(arg.ShortName)]) {
			missing = append(missing, arg)
		}
	}
	if len(missing) != 0 {
		fmt.Fprint(flag.CommandLine.Output(), errorHelp(missing))
		os.Exit(-1)
	}
}

/* initialize clap parser */
func init() {
	/* determine if output (stderr, same as flag package) goes to a color capable terminal */
	if fileInfo, err := os.Stderr.Stat(); err == nil && (fileInfo.Mode()&os.ModeCharDevice) != 0 {
		_, noColor := os.LookupEnv("NO_COLOR")
		isTerminal = !noColor && os.Getenv("TERM") != "dumb"
	}
	/* replace flag usage */
	flag.Usage = func() {
		maxLength := int(0)
		/* get the length of the longest argument */
		for _, arg := range argHelp {
			plain, _ := arg.option()
			maxLength = max(maxLength, utf8.RuneCountInString(plain))
		}
		/* print header */
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "%s\n", usageHeader())
		/* print options */
		for _, arg := range argHelp {
			fmt.Fprint(out, arg.String(maxLength))
		}
	}
}
