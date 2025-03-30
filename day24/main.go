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

	aMap := readInput(fileName)

	fmt.Printf("The z wires output %d\n", Part1(aMap))
}

func readInput(filePath string) Connections {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	aMap := Connections{}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			value, _ := strconv.Atoi(strings.TrimSpace(parts[1]))

			aMap[parts[0]] = Connection{value: uint8(value)}
		}

		if strings.Contains(line, "->") {
			parts := strings.Split(line, " ")
			aMap[parts[4]] = Connection{
				value:    255,
				operator: parts[1],
				operand1: parts[0],
				operand2: parts[2],
			}
		}
	}

	return aMap
}

func Part1(aMap Connections) int {
	binaryDigits := ""

	for _, k := range aMap.Keys() {
		if k[0] != "z"[0] {
			continue
		}
		newDigit := strconv.Itoa(int(SettleValue(&aMap, k)))
		binaryDigits = fmt.Sprintf("%s%s", newDigit, binaryDigits)
	}
	fmt.Println(binaryDigits)

	decimalValue, _ := strconv.ParseInt(binaryDigits, 2, 64)

	return int(decimalValue)
}

func SettleValue(aMap *Connections, wire string) uint8 {
	value := (*aMap)[wire].value
	operator := (*aMap)[wire].operator
	operand1 := (*aMap)[wire].operand1
	operand2 := (*aMap)[wire].operand2

	operand1Val := (*aMap)[operand1].value
	operand2Val := (*aMap)[operand2].value

	if value != 255 {
		return value
	}

	if operand1Val == 255 {
		operand1Val = SettleValue(aMap, operand1)
	}

	if operand2Val == 255 {
		operand2Val = SettleValue(aMap, operand2)
	}

	switch operator {
	case "OR":
		value = operand1Val | operand2Val
	case "XOR":
		value = operand1Val ^ operand2Val
	case "AND":
		value = operand1Val & operand2Val
	}

	tmp := (*aMap)[wire]
	tmp.value = value
	(*aMap)[wire] = tmp

	return value
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

type Connection struct {
	value    uint8
	operator string
	operand1 string
	operand2 string
}

func (conn *Connection) Calculate() {
	if conn.operand1 != "" && conn.operand2 != "" {
		switch conn.operator {
		case "XOR":
			(*conn).value = 0
		case "OR":
			(*conn).value = 1
		case "AND":
			(*conn).value = 2
		}
	}
}
