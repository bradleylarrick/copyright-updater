package handlers

import (
	"bufio"
	"fmt"
	"os"
)

var (
	scanner *bufio.Scanner
	writer  *bufio.Writer
	header  string
	footer  string
	prefix  string
	lines   []string
	ndx     int
)

/*
 * Starts the process of updating the copyright header in the source file.
 */
func startProcess(srcFile *os.File, destFile *os.File, hdr string, ftr string, pref string) error {
	header = hdr
	footer = ftr
	prefix = pref
	ndx = 0
	lines = lines[:ndx]	// clear the lines slice

	scanner = bufio.NewScanner(srcFile)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	writer = bufio.NewWriter(destFile)
	return nil
}

/*
 * findHeader finds the start of the copyright header in the source file and,
 * if found, skips to the end of the header.
 */
func findHeader() error {

	return nil
}

/*
 * Writes the copyright header to the output file.
 */
func writeCopyright(copyright *[]string) error {
	writeLine(header)
	for _, line := range *copyright {
		writeLine(prefix + " " + line)
	}
	writeLine(footer)
	return nil
}

/*
 * Copies the remaining lines from the source file to the output file.
 */
func endProcess() error {
	for i := ndx; i < len(lines); i++ {
		line := lines[i]
		// Add a blank line after the header if the next line is non-empty
		if i == ndx && len(line) > 0 {
			writeLine("")
		}
		writeLine(line)
	}

	err := writer.Flush()
	if err != nil {
		return err
	}

	// Check for reading phase errors
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(4)
	}

	return nil
}

// Write a single line to the output file.
func writeLine(line string) {
	_, err := writer.WriteString(line + "\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing line: %s\n", err)
		os.Exit(4)
	}
}
