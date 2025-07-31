// Package ui provides functionality for displaying different colors and layouts in the terminal.
package ui

import (
	"errors"
	"fmt"
)

// TODO: ui stuff will go here, this code is here to learn about go doc
// Todo: what should PrintColor return? A string?

// Color represents a color code for terminal text.
type Color int

// Supported terminal colors.
const (
	Red   Color = iota // Red color code
	Green              // Green color code
	Blue               // Blue color code
)

// PrintColor prints a given text in the specified color.
// Parameters:
//
//	text: the text to be printed.
//	color: the color in which the text should be printed.
//
// Returns:
//
//	error: an error if the color is not supported or if printing fails.
func PrintColor(text string, color Color) error {
	switch color {
	case Red:
		fmt.Print("\033[31m" + text + "\033[0m")
	case Green:
		fmt.Print("\033[32m" + text + "\033[0m")
	case Blue:
		fmt.Print("\033[34m" + text + "\033[0m")
	default:
		return errors.New("unsupported color or printing failure")
	}
	return nil
}
