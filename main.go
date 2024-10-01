package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	// https://github.com/alexflint/go-arg
	"github.com/alexflint/go-arg"
)

func cleanEthernetAddress(input string) (string, error) {
	// Create a regular expression to match characters that are not [a-fA-F0-9]
	re := regexp.MustCompile("[^a-fA-F0-9]")

	// Remove any characters that don't match [a-fA-F0-9]
	stripped := re.ReplaceAllString(input, "")

	// Check if the resulting string has 12 characters
	if len(stripped) != 12 {
		return "", errors.New("Invalid Ethernet address")
	}

	// Convert the resulting string to lowercase
	cleaned := strings.ToLower(stripped)

	return cleaned, nil
}

func getCiscoDelimitedOUI(input string) (string, error) {
	cleaned, err := cleanEthernetAddress(input)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Insert colons after every two characters
	var result strings.Builder
	for i, char := range cleaned {
		if i > 0 && i%4 == 0 {
			result.WriteRune('.')
		}
		result.WriteRune(char)
	}

	return result.String(), nil
}

func getNonDelimitedUppercaseOUI(input string) (string, error) {
	cleaned, err := cleanEthernetAddress(input)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return strings.ToUpper(cleaned), nil
}

func getColonDelimitedOUI(input string) (string, error) {
	cleaned, err := cleanEthernetAddress(input)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Insert colons after every two characters
	var result strings.Builder
	for i, char := range cleaned {
		if i > 0 && i%2 == 0 {
			result.WriteRune(':')
		}
		result.WriteRune(char)
	}

	return result.String(), nil
}

func getColonDelimitedUppercaseOUI(input string) (string, error) {
	cleaned, err := cleanEthernetAddress(input)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Insert colons after every two characters
	var result strings.Builder
	for i, char := range cleaned {
		if i > 0 && i % 2 == 0 {
			result.WriteRune(':')
		}
		result.WriteRune(char)
	}

	return strings.ToUpper(result.String()), nil
}

func getHyphenDelimitedUppercaseOUI(input string) (string, error) {
	cleaned, err := cleanEthernetAddress(input)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Insert hyphens after every two characters
	var result strings.Builder
	for i, char := range cleaned {
		if i > 0 && i%2 == 0 {
			result.WriteRune('-')
		}
		result.WriteRune(char)
	}

	return strings.ToUpper(result.String()), nil
}

// The be-all function.
func formatOUI(input string, format string) (string, error) {
	spacer := 2
	delimiter := "-" // Can only be one char.

	switch format {
	case "cisco":
		spacer = 4
		delimiter = "."
	case "colon":
		spacer = 2
		delimiter = ":"
	default:
		spacer = 2
		delimiter = "-"
	}

	fmt.Println("Format: ", format, " Spacer: ", spacer, " Delim: ", delimiter)

	cleaned, err := cleanEthernetAddress(input)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Insert delimiter after every spacer characters
	var result strings.Builder
	for i, char := range cleaned {
		if i > 0 && i % spacer == 0 {
			r := rune(delimiter[0])
			result.WriteRune(r)
		}
		result.WriteRune(char)
	}

	return result.String(), nil
}

func printAllFormats(input string) (error) {
	formatFuncs := []func(string) (string, error){
		getCiscoDelimitedOUI,
		getNonDelimitedUppercaseOUI,
		getColonDelimitedOUI,
		getColonDelimitedUppercaseOUI,
		getHyphenDelimitedUppercaseOUI,
	}

	for _, fn := range formatFuncs {
		result, err := fn(input)
		if err != nil {
			return err
		}
		fmt.Println(result)
	}

	return nil
}

func printUsage() {
	fmt.Println("Usage: macf [-f/--format <format>] [-u/--uppercase] <ethernet_address>")
	fmt.Println("  Options:")
	fmt.Println("")
	fmt.Println("    --format options:")
	fmt.Println("        cisco        => a1b2.c3d4.e5f6")
	fmt.Println("        colon        => a1:b2:c3:d4:e5:f6")
	fmt.Println("        hyphen       => a1-b2-c3-d4-e5-f6")
	fmt.Println("        none         => a1b2c3d4e5f6")

	fmt.Println("")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	var args struct {
		Format string `arg:"-f,--format"`
		Address string `arg:"positional"`
		Upper bool `arg:"-u,--uppercase" default:"false"`
	}
	arg.MustParse(&args)

	input := args.Address

	if args.Format != "" {
		// Format it!
		result, err := formatOUI( input, args.Format)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	// Check upper flag and set input accordingly.
	if args.Upper == true {
		fmt.Println(strings.ToUpper(result))
	} else {
		fmt.Println(result)
	}
		
	} else {
		err := printAllFormats(input)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}