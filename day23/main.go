package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
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

	fmt.Printf("There are %s sets of 3 computers with at least 1 't' member\n", Part1(links))
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

func Part1(links []Link) string {
	parties := []string{}

	for i := range links {
		for j := range links {
			if i == j {
				continue
			}

			if !links[i].IsSpecial() {
				continue
			}

			if !links[i].Matches(links[j]) {
				continue
			}

			var missing Link

			if links[i].a == links[j].a {
				missing = Link{a: links[i].b, b: links[j].b}
			} else if links[i].b == links[j].b {
				missing = Link{a: links[i].a, b: links[j].a}
			} else if links[i].a == links[j].b {
				missing = Link{a: links[i].b, b: links[j].a}
			} else {
				missing = Link{a: links[i].a, b: links[j].b}
			}

			if slices.Contains(links, missing) {
				party := partyString([]Link{links[i], links[j], missing})

				if !slices.Contains(parties, party) {
					parties = append(parties, party)
				}
			}
		}
	}

	return strconv.Itoa(len(parties))
}

type Link struct {
	a string
	b string
}

func (link Link) IsEqual(other Link) bool {
	return (link.a == other.a && link.b == other.b) ||
		(link.a == other.b && link.b == other.a)
}

func (link Link) IsSpecial() bool {
	return link.a[0] == "t"[0] || link.b[0] == "t"[0]
}

func (link Link) Matches(other Link) bool {
	return (link.a == other.a ||
		link.b == other.b ||
		link.a == other.b ||
		link.b == other.a) &&
		link != other
}

func partyString(links []Link) string {
	nodes := []string{}

	for _, link := range links {
		if !slices.Contains(nodes, link.a) {
			nodes = append(nodes, link.a)
		}

		if !slices.Contains(nodes, link.b) {
			nodes = append(nodes, link.b)
		}
	}

	sort.Strings(nodes)

	return strings.Join(nodes, ",")
}
