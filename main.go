package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type formatStyle string

const (
	formatCisco  formatStyle = "cisco"
	formatColon  formatStyle = "colon"
	formatHyphen formatStyle = "hyphen"
	formatNone   formatStyle = "none"
)

var errInvalidMAC = errors.New("invalid Ethernet address")

func normalizeMAC(input string) (string, error) {
	if parsed, err := net.ParseMAC(input); err == nil && len(parsed) == 6 {
		return strings.ToLower(hex.EncodeToString(parsed)), nil
	}

	decoded, err := hex.DecodeString(input)
	if err != nil || len(decoded) != 6 {
		return "", errInvalidMAC
	}

	return strings.ToLower(input), nil
}

func parseFormat(input string) (formatStyle, error) {
	switch strings.ToLower(input) {
	case string(formatCisco):
		return formatCisco, nil
	case string(formatColon):
		return formatColon, nil
	case string(formatHyphen):
		return formatHyphen, nil
	case string(formatNone):
		return formatNone, nil
	default:
		return "", fmt.Errorf("invalid format %q (must be one of: cisco, colon, hyphen, none)", input)
	}
}

func formatMAC(cleaned string, style formatStyle, uppercase bool) string {
	var result string

	switch style {
	case formatCisco:
		result = joinGroups(cleaned, 4, ".")
	case formatColon:
		result = joinGroups(cleaned, 2, ":")
	case formatHyphen:
		result = joinGroups(cleaned, 2, "-")
	case formatNone:
		result = cleaned
	}

	if uppercase {
		return strings.ToUpper(result)
	}

	return result
}

func joinGroups(input string, width int, delimiter string) string {
	var result strings.Builder
	result.Grow(len(input) + (len(input)/width - 1))

	for i := 0; i < len(input); i += width {
		if i > 0 {
			result.WriteString(delimiter)
		}
		result.WriteString(input[i : i+width])
	}

	return result.String()
}

func allFormats(cleaned string) []string {
	return []string{
		formatMAC(cleaned, formatCisco, false),
		formatMAC(cleaned, formatNone, true),
		formatMAC(cleaned, formatColon, false),
		formatMAC(cleaned, formatColon, true),
		formatMAC(cleaned, formatHyphen, true),
	}
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "Usage: macf [-f format] [-u] <ethernet_address>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Formats:")
	fmt.Fprintln(w, "  cisco   a1b2.c3d4.e5f6")
	fmt.Fprintln(w, "  colon   a1:b2:c3:d4:e5:f6")
	fmt.Fprintln(w, "  hyphen  a1-b2-c3-d4-e5-f6")
	fmt.Fprintln(w, "  none    a1b2c3d4e5f6")
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("macf", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		usage(stderr)
	}

	var format string
	var uppercase bool

	fs.StringVar(&format, "format", "", "output format: cisco, colon, hyphen, none")
	fs.StringVar(&format, "f", "", "output format: cisco, colon, hyphen, none")
	fs.BoolVar(&uppercase, "uppercase", false, "uppercase output")
	fs.BoolVar(&uppercase, "u", false, "uppercase output")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	if fs.NArg() != 1 {
		usage(stderr)
		return 2
	}

	cleaned, err := normalizeMAC(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", errInvalidMAC)
		return 1
	}

	if format == "" {
		for _, line := range allFormats(cleaned) {
			fmt.Fprintln(stdout, line)
		}
		return 0
	}

	style, err := parseFormat(format)
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return 1
	}

	fmt.Fprintln(stdout, formatMAC(cleaned, style, uppercase))
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}
