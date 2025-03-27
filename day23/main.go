package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	links := readInput(fileName)

	fmt.Println(links)
}

func readInput(filePath string) []Link {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	links := []Link{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		nodes := strings.Split(line, "-")
		links = append(links, Link{a: nodes[0], b: nodes[1]})
	}

	return links
}

func Part1(links []Link) int {
	return 0
}

type Link struct {
	a string
	b string
}
