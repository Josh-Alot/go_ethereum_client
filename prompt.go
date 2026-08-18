package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func ReadPassword(prompt string) ([]byte, error) {
	fd := int(os.Stdin.Fd())

	if term.IsTerminal(fd) {
		fmt.Fprint(os.Stderr, prompt)

		res, err := term.ReadPassword(fd)
		fmt.Fprintf(os.Stderr, "\n")
		if err != nil {
			return nil, fmt.Errorf("read password error: %w", err)
		}

		return res, nil
	}

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("read password error: %w", err)
	}

	trimmedLine := strings.TrimSuffix(line, "\n")
	return []byte(trimmedLine), nil
}
