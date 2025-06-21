package clap

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

/* module constants */
const (
	argString = iota
	argInt
	argBool
)

/* argumentHelp stores command line arguments info */
type argumentHelp struct {
	ArgType   int
	ArgLen    int
	ShortName string
	LongName  string
	HelpText  string
	Required  bool
	Value     interface{}
}

/* String representation of the argument help */
func (a argumentHelp) String(isTerminal bool, maxLen int) (result string) {
	/* generate color format */
	terminalFmt := func(maxLen int) (result string) {
		switch a.ArgType {
		case argString, argInt:
			result = fmt.Sprintf(
				"  \033[01;38m-%%s, --%%s\033[00m <%%s>%s%%s\n",
				strings.Repeat(" ", maxLen-a.ArgLen+4),
			)
		case argBool:
			result = fmt.Sprintf(
				"  \033[01;38m-%%s, --%%s\033[00m%s%%s\n",
				strings.Repeat(" ", maxLen-a.ArgLen+7),
			)
		}
		return
	}
	/* generate non-color format */
	nonTerminalFmt := func(maxLen int) (result string) {
		switch a.ArgType {
		case argString, argInt:
			result = fmt.Sprintf(
				"  -%%s, --%%s <%%s>%s%%s\n",
				strings.Repeat(" ", maxLen-a.ArgLen+4),
			)
		case argBool:
			result = fmt.Sprintf(
				"  -%%s, --%%s%s%%s\n",
				strings.Repeat(" ", maxLen-a.ArgLen+7),
			)
		}
		return
	}
	/* return result */
	var fmtStr string
	if isTerminal {
		fmtStr = terminalFmt(maxLen)
	} else {
		fmtStr = nonTerminalFmt(maxLen)
	}
	switch a.ArgType {
	case argString, argInt:
		result = fmt.Sprintf(
			fmtStr,
			a.ShortName,
			a.LongName,
			strings.ToUpper(a.LongName),
			a.HelpText,
		)
	case argBool:
		result = fmt.Sprintf(
			fmtStr,
			a.ShortName,
			a.LongName,
			a.HelpText,
		)
	}
	return
}

/* global variables */
var (
	argHelp    []argumentHelp
	isTerminal bool
)

/* String command line option */
func String(shortName, longName string, defaultValue string, helpText string, required bool) *string {
	result := new(string)
	flag.StringVar(result, shortName, defaultValue, helpText)
	flag.StringVar(result, longName, defaultValue, helpText)
	argHelp = append(
		argHelp,
		argumentHelp{
			ArgType:   argString,
			ArgLen:    len(shortName) + 2*len(longName),
			ShortName: shortName,
			LongName:  longName,
			HelpText:  helpText,
			Required:  required,
			Value:     result,
		},
	)
	return result
}

/* Int command line option */
func Int(shortName, longName string, defaultValue int, helpText string, required bool) *int {
	result := new(int)
	flag.IntVar(result, shortName, defaultValue, helpText)
	flag.IntVar(result, longName, defaultValue, helpText)
	argHelp = append(
		argHelp,
		argumentHelp{
			ArgType:   argInt,
			ArgLen:    len(shortName) + 2*len(longName),
			ShortName: shortName,
			LongName:  longName,
			HelpText:  helpText,
			Required:  required,
			Value:     result,
		},
	)
	return result
}

/* Bool command line option */
func Bool(shortName, longName string, defaultValue bool, helpText string, required bool) *bool {
	result := new(bool)
	flag.BoolVar(result, shortName, defaultValue, helpText)
	flag.BoolVar(result, longName, defaultValue, helpText)
	argHelp = append(
		argHelp,
		argumentHelp{
			ArgType:   argBool,
			ArgLen:    len(shortName) + len(longName),
			ShortName: shortName,
			LongName:  longName,
			HelpText:  helpText,
			Required:  required,
			Value:     result,
		},
	)
	return result
}

/* errorHelp renders error message when required arguments are missing */
func errorHelp(isTerminal bool) string {
	/* prepare error message header */
	fmtStr := "\033[01;31merror:\033[00m the following arguments are not provided:\n"
	if !isTerminal {
		fmtStr = "error: the following arguments are not provided:\n"
	}
	/* add arguments to error message */
	for _, arg := range argHelp {
		if arg.Required {
			argStr := "  \033[02;32m--%s <%s>\033[00m\n"
			if !isTerminal {
				argStr = "  --%s <%s>\n"
			}
			fmtStr = fmt.Sprintf(
				"%s%s",
				fmtStr,
				fmt.Sprintf(argStr, arg.LongName, strings.ToUpper(arg.LongName)),
			)
		}
	}
	/* add usage information */
	usageHelp := "\033[04m\033[01;38mUsage:\033[00m \033[01;38m%s\033[00m"
	if !isTerminal {
		usageHelp = "Usage: %s"
	}
	for _, arg := range argHelp {
		if !arg.Required {
			continue
		}
		argHelp := "\033[01;38m--%s\033[00m <%s>"
		if !isTerminal {
			argHelp = "--%s <%s>"
		}
		usageHelp = fmt.Sprintf(
			"%s %s",
			fmt.Sprintf(
				usageHelp,
				os.Args[0],
			),
			fmt.Sprintf(
				argHelp,
				arg.LongName,
				strings.ToUpper(arg.LongName),
			),
		)
	}
	fmtStr = fmt.Sprintf("%s\n%s\n", fmtStr, usageHelp)
	/* add help tip */
	helpTip := "For more information, try '\033[01;38m--help\033[00m'."
	if !isTerminal {
		helpTip = "For more information, try '--help'."
	}

	return fmt.Sprintf("%s\n%s\n", fmtStr, helpTip)
}

/* usageHeader prints flags usage header */
func usageHeader(isTerminal bool) string {
	var header string
	var fmtStr string
	header = "\033[04m\033[01;38mUsage:\033[00m \033[01;38m%s\033[00m [OPTIONS]"
	if !isTerminal {
		header = "Usage: %s [OPTIONS]"
	}
	for _, arg := range argHelp {
		if arg.Required {
			fmtStr = "%s \033[01;38m--%s\033[00m <%s>"
			if !isTerminal {
				fmtStr = "%s --%s <%s>"
			}
			header = fmt.Sprintf(
				fmtStr,
				header,
				arg.LongName,
				strings.ToUpper(arg.LongName),
			)
		}
	}
	optHdr := "\033[04m\033[01;38mOptions:\033[00m"
	if !isTerminal {
		optHdr = "Options:"
	}
	return fmt.Sprintf(
		"%s\n\n%s",
		fmt.Sprintf(header, os.Args[0]),
		optHdr,
	)
}

/* Parse is a wrapper around flag.Parse function */
func Parse() {
	flag.Parse()
	/* check if required arguments are provided */
	for _, arg := range argHelp {
		/* do not check non-mandatory arguments */
		if !arg.Required {
			continue
		}
		/* do not check bool arguments */
		if arg.ArgType == argBool {
			continue
		}
		/* make sure value is provided */
		if len(*arg.Value.(*string)) == 0 {
			fmt.Printf(errorHelp(isTerminal))
			os.Exit(-1)
		}
	}
}

/* initialize clap parser */
func init() {
	/* determine if program is running inside a terminal */
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		isTerminal = true
	}
	/* replace flag usage */
	flag.Usage = func() {
		maxLength := int(0)
		/* get the length of the longest argument */
		for idx := range argHelp {
			if argHelp[idx].ArgLen > maxLength {
				maxLength = argHelp[idx].ArgLen
			}
		}
		/* print header */
		fmt.Printf("%s\n", usageHeader(isTerminal))
		/* print options */
		for _, arg := range argHelp {
			fmt.Printf(arg.String(isTerminal, maxLength))
		}
	}
}
