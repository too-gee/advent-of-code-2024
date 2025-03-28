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
	fmt.Printf("The password for the largest party is %s\n", Part2(links))
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

func Part1(links Links) string {
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

			tmp := Links{links[i], links[j]}
			missing := tmp.MissingLinks()[0]

			if slices.Contains(links, missing) {
				party := Links{links[i], links[j], missing}
				password := party.GetPassword()

				if !slices.Contains(parties, password) {
					parties = append(parties, password)
				}
			}
		}
	}

	return strconv.Itoa(len(parties))
}

func Part2(links Links) string {
	password := ""
	biggestParty := 0

	queue := DumbQueue{}

	for i := range links {
		queue.push(Links{links[i]})
	}

	for len(queue) > 0 {
		current := queue.pop()

		missing := current.MissingLinks()

		foundMissingLinks := true

		for _, missingLink := range missing {
			if links.Contains(missingLink) {
				current = append(current, missingLink)
			} else {
				foundMissingLinks = false
				break
			}
		}

		if !foundMissingLinks {
			continue
		}

		if len(current) > biggestParty {
			biggestParty = len(current)
			password = current.GetPassword()
		}

		currentNodes := current.GetNodes()

		for i := range links {
			if current.Contains(links[i]) {
				continue
			}

			hasInterestingNode := false

			for _, currentNode := range currentNodes {
				if links[i].HasNode(currentNode) {
					hasInterestingNode = true
					break
				}
			}

			if hasInterestingNode {
				queue.push(current.CopyAppend(links[i]))
			}
		}
	}

	return password
}

type Link struct {
	a string
	b string
}

func (link Link) Equal(other Link) bool {
	return (link.a == other.a && link.b == other.b) ||
		(link.b == other.a && link.a == other.b)
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

func (link Link) HasNode(node string) bool {
	return link.a == node || link.b == node
}

type Links []Link

func (links Links) Contains(link Link) bool {
	for i := range links {
		if link.Equal(links[i]) {
			return true
		}
	}

	return false
}

func (links Links) CopyAppend(link Link) Links {
	newLinks := make(Links, len(links)+1)
	copy(newLinks, links)
	newLinks[len(links)] = link

	return newLinks
}

func (links Links) GetNodes() []string {
	nodes := []string{}

	for _, link := range links {
		if !slices.Contains(nodes, link.a) {
			nodes = append(nodes, link.a)
		}

		if !slices.Contains(nodes, link.b) {
			nodes = append(nodes, link.b)
		}
	}

	return nodes
}

func (links Links) GetPassword() string {
	nodes := links.GetNodes()

	sort.Strings(nodes)

	return strings.Join(nodes, ",")
}

func (links Links) MissingLinks() Links {
	if len(links) == 1 {
		return Links{}
	}

	nodes := links.GetNodes()

	missingLinks := Links{}

	for i := range nodes {
		for j := i + 1; j < len(nodes); j++ {
			candidate := Link{a: nodes[i], b: nodes[j]}

			if !links.Contains(candidate) {
				missingLinks = append(missingLinks, Link{a: nodes[i], b: nodes[j]})
			}
		}
	}

	return missingLinks
}

type DumbQueue []Links

func (q *DumbQueue) push(item Links) { *q = append(*q, item) }

func (q *DumbQueue) pop() Links {
	item := (*q)[len(*q)-1]
	*q = (*q)[0 : len(*q)-1]
	return item
}
