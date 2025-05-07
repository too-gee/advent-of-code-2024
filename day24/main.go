package main

import (
	"bufio"
	"fmt"
	"os"
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

	conns := readInput(fileName)

	fmt.Printf("The z wires output %s\n", Part1(conns))
}

func readInput(filePath string) Connections {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	conns := Connections{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			value, _ := strconv.Atoi(strings.TrimSpace(parts[1]))

			conns[parts[0]] = Connection{value: uint8(value)}
		}

		if strings.Contains(line, "->") {
			parts := strings.Split(line, " ")
			conns[parts[4]] = Connection{
				value:    255,
				operator: parts[1],
				operand1: parts[0],
				operand2: parts[2],
			}
		}
	}

	return conns
}

func Part1(conns Connections) string {
	binaryDigits := ""

	for _, k := range conns.Keys() {
		if k[0] != "z"[0] {
			continue
		}
		newDigit := strconv.Itoa(int(SettleValue(&conns, k)))
		binaryDigits = fmt.Sprintf("%s%s", newDigit, binaryDigits)
	}

	return strconv.Itoa(binToDec(binaryDigits))
}

func SettleValue(conns *Connections, wire string) uint8 {
	value := (*conns)[wire].value
	operator := (*conns)[wire].operator
	operand1 := (*conns)[wire].operand1
	operand2 := (*conns)[wire].operand2

	operand1Val := (*conns)[operand1].value
	operand2Val := (*conns)[operand2].value

	if value != 255 {
		return value
	}

	if operand1Val == 255 {
		operand1Val = SettleValue(conns, operand1)
	}

	if operand2Val == 255 {
		operand2Val = SettleValue(conns, operand2)
	}

	switch operator {
	case "OR":
		value = operand1Val | operand2Val
	case "XOR":
		value = operand1Val ^ operand2Val
	case "AND":
		value = operand1Val & operand2Val
	}

	tmp := (*conns)[wire]
	tmp.value = value
	(*conns)[wire] = tmp

	return value
}

func binToDec(binary string) int {
	num, _ := strconv.ParseInt(binary, 2, 64)
	return int(num)
}

type Connection struct {
	value    uint8
	operator string
	operand1 string
	operand2 string
}

type Connections map[string]Connection

func (conns Connections) Keys() []string {
	keys := make([]string, 0, len(conns))

	for k := range conns {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}
