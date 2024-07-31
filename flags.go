package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type appFlags struct {
	configFile     string
	versionAndExit bool
}

func attachFlags(flags *appFlags) {
	flag.StringVar(&flags.configFile, "config", "", "path to key configuration file")
	flag.BoolVar(&flags.versionAndExit, "version", false, "print version information")

	// Support the -file flag for backward compatibility
	var backcompatConfigFile string
	flag.StringVar(&backcompatConfigFile, "file", "", "")

	flag.Usage = flagUsage()
	flag.Parse()

	if flags.configFile == "" {
		flags.configFile = backcompatConfigFile
	}
}

// Adapted from https://cs.opensource.google/go/go/+/refs/tags/go1.22.5:src/flag/flag.go;l=607
// to suppress displaying help text for flags without usage
func flagUsage() func() {
	isZeroValue := func(f *flag.Flag, value string) (ok bool, err error) {
		// Build a zero value of the flag's Value type, and see if the
		// result of calling its String method equals the value passed in.
		// This works unless the Value type is itself an interface type.
		typ := reflect.TypeOf(f.Value)
		var z reflect.Value
		if typ.Kind() == reflect.Pointer {
			z = reflect.New(typ.Elem())
		} else {
			z = reflect.Zero(typ)
		}
		// Catch panics calling the String method, which shouldn't prevent the
		// usage message from being printed, but that we should report to the
		// user so that they know to fix their code.
		defer func() {
			if e := recover(); e != nil {
				if typ.Kind() == reflect.Pointer {
					typ = typ.Elem()
				}
				err = fmt.Errorf("panic calling String method on zero %v for flag %s: %v", typ, f.Name, e)
			}
		}()
		return value == z.Interface().(flag.Value).String(), nil
	}

	return func() {
		out := os.Stderr
		fmt.Fprintf(out, "Usage of %s:\n", appName)

		// Customized: skip flags with empty usage
		var isZeroValueErrs []error
		flag.VisitAll(func(f *flag.Flag) {
			if f.Usage == "" {
				return
			}

			var b strings.Builder
			fmt.Fprintf(&b, "  -%s", f.Name) // Two spaces before -; see next two comments.
			name, usage := flag.UnquoteUsage(f)
			if len(name) > 0 {
				b.WriteString(" ")
				b.WriteString(name)
			}
			// Boolean flags of one ASCII letter are so common we
			// treat them specially, putting their usage on the same line.
			if b.Len() <= 4 { // space, space, '-', 'x'.
				b.WriteString("\t")
			} else {
				// Four spaces before the tab triggers good alignment
				// for both 4- and 8-space tab stops.
				b.WriteString("\n    \t")
			}
			b.WriteString(strings.ReplaceAll(usage, "\n", "\n    \t"))

			// Print the default value only if it differs to the zero value
			// for this flag type.
			if isZero, err := isZeroValue(f, f.DefValue); err != nil {
				isZeroValueErrs = append(isZeroValueErrs, err)
			} else if !isZero {
				fmt.Fprintf(&b, " (default %v)", f.DefValue)
			}
			fmt.Fprint(out, b.String(), "\n")
		})
		// If calling String on any zero flag.Values triggered a panic, print
		// the messages after the full set of defaults so that the programmer
		// knows to fix the panic.
		if errs := isZeroValueErrs; len(errs) > 0 {
			fmt.Fprintln(out)
			for _, err := range errs {
				fmt.Fprintln(out, err)
			}
		}
	}
}
