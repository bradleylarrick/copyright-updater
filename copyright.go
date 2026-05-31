package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Copyright struct {
	copyright []string
}

func NewCopyright(template string, year string) *Copyright {
	return &Copyright{copyright: loadTemplate(template, year)}
}

// Loads the given copyright template and replaces the ${year} placeholder with the current year.
func loadTemplate(template string, year string) []string {

	var copyright []string

	file, err := os.Open(template)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed: %s\n", err.Error())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "${year}", year)
		copyright = append(copyright, line)
	}

	if scanner.Err() != nil {
		fmt.Fprintf(os.Stderr, "Failed: %s\n", scanner.Err().Error())
		os.Exit(1)
	}

	return copyright
}
