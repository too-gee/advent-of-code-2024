package main

import (
	"bufio"
	"fmt"
	"math"
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

	aMap := readInput(fileName)

	fmt.Println(aMap)

	fmt.Printf("There are %s sets of 3 computers with at least 1 't' member\n", Part1(aMap))
	fmt.Printf("The password for the largest party is %s\n", Part2(aMap))
}

func readInput(filePath string) map[string][]string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	list := map[string][]string{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		nodes := strings.Split(line, "-")

		for i := range nodes {
			a := i
			b := int(math.Abs(float64(i - 1)))

			if _, ok := list[nodes[a]]; !ok {
				list[nodes[a]] = []string{}
			}

			list[nodes[a]] = append(list[nodes[a]], nodes[b])
		}
	}

	for k := range list {
		slices.Sort(list[k])
		list[k] = slices.Compact(list[k])
	}

	return list
}

func Part1(aMap map[string][]string) string {
	parties := []string{}

	for node, neighbors := range aMap {
		if node[0] != "t"[0] {
			continue
		}

		for _, neighbor := range neighbors {
			for _, second := range aMap[neighbor] {
				if second == node {
					continue
				}

				if slices.Contains(aMap[second], node) {
					nodes := []string{node, neighbor, second}
					sort.Strings(nodes)
					password := strings.Join(nodes, ",")

					if !slices.Contains(parties, password) {
						parties = append(parties, password)
					}
				}
			}
		}
	}

	return strconv.Itoa(len(parties))
}

func Part2(aMap map[string][]string) string {
	password := ""
	biggestParty := 0

	for i := range aMap {
		party := append(aMap[i], i)

		for j := 0; j < len(party)-1; j++ {
			party = commonNodes(party, append(aMap[party[j]], party[j]))
		}

		if len(party) > biggestParty {
			biggestParty = len(party)
			sort.Strings(party)
			password = strings.Join(party, ",")
		}
	}

	return password
}

func commonNodes(a, b []string) []string {
	lookup := map[string]bool{}

	for i := range a {
		lookup[a[i]] = true
	}

	common := []string{}

	for i := range b {
		if lookup[b[i]] {
			common = append(common, b[i])
		}
	}

	return common
}
