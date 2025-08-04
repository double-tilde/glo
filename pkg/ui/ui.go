package ui

import (
	"errors"
	"fmt"
)

// TODO: ui stuff will go here, this code is here to learn about go doc

type Color int

const (
	Red   Color = iota // Red color code
	Green              // Green color code
	Blue               // Blue color code
)

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

func Printing() {
	// Green backgrounds
	fmt.Println("\033[48;5;22m \033[0m")
	fmt.Println("\033[48;5;34m \033[0m")
	fmt.Println("\033[48;5;46m \033[0m")

	// Blue backgrounds
	fmt.Println("\033[48;5;18m \033[0m")
	fmt.Println("\033[48;5;19m \033[0m")
	fmt.Println("\033[48;5;21m \033[0m")

	// Red backgrounds
	fmt.Println("\033[48;5;52m \033[0m")
	fmt.Println("\033[48;5;88m \033[0m")
	fmt.Println("\033[48;5;124m \033[0m")

	// var grid [][]string

	for i := 0; i <= 7; i++ {
		for j := 0; j <= 53; j++ {
			fmt.Print("x")
		}
		fmt.Println()
	}
}
