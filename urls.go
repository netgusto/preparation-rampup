package main

import "strings"

func getURLs(data []byte) ([]string, error) {
	lines := strings.Split(string(data), "\n")
	return removeEmptyLines(lines), nil
}

func removeEmptyLines(lines []string) []string {
	nonEmptyLines := []string{}
	for _, line := range lines {
		if len(strings.TrimSpace(line)) > 0 {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	return nonEmptyLines
}
