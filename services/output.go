package services

import "fmt"

type ConsoleOutput struct{}

func NewConsoleOutput() Output {
	return &ConsoleOutput{}
}

func (c *ConsoleOutput) Printf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (c *ConsoleOutput) Println(args ...any) {
	fmt.Println(args...)
}
