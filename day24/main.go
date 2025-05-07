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
	fmt.Printf("The swapped gates are %s\n", Part2(conns))
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

func (conns Connections) Find(operand1 string, operand2 string, operator string) string {
	if operand1 == "" && operand2 == "" {
		return ""
	}

	for i := range conns {
		if (operator == conns[i].operator) &&
			(operand1 == "" || operand1 == conns[i].operand1 || operand1 == conns[i].operand2) &&
			(operand2 == "" || operand2 == conns[i].operand1 || operand2 == conns[i].operand2) {
			return i
		}
	}

	return ""
}

func Part2(conns Connections) string {
	bits := 0
	for i := range conns {
		if i[0] == "x"[0] {
			bits++
		}
	}

	oops := []string{}
	carry_output := ""

	for i := range bits {
		xIn := fmt.Sprintf("x%02d", i)
		yIn := fmt.Sprintf("y%02d", i)
		expectedZOut := fmt.Sprintf("z%02d", i)

		if i == 0 {
			zOut := conns.Find(xIn, yIn, "XOR")
			cOut := conns.Find(xIn, yIn, "AND")

			if zOut != expectedZOut {
				oops = append(oops, zOut)
			}

			if cOut[0] == "z"[0] {
				oops = append(oops, cOut)
				carry_output = ""
			} else {
				carry_output = cOut
			}
		} else {
			iOut := conns.Find(xIn, yIn, "XOR")
			kOut := conns.Find(xIn, yIn, "AND")

			jOut := conns.Find(carry_output, "", "AND")
			if jOut == "" || carry_output == "" {
				jOut = conns.Find(iOut, "", "AND")
			}
			jGate := conns[jOut]

			if iOut != jGate.operand1 && iOut != jGate.operand2 {
				oops = append(oops, iOut)
			}

			zOut := conns.Find(jGate.operand1, jGate.operand2, "XOR")
			if zOut != expectedZOut {
				oops = append(oops, zOut)
			}

			cOut := conns.Find(jOut, "", "OR")
			if cOut == "" {
				cOut = conns.Find(kOut, "", "OR")
			}
			cGate := conns[cOut]

			if cGate.operand1 != kOut && cGate.operand2 != kOut {
				oops = append(oops, kOut)
			}

			if cGate.operand1 != jOut && cGate.operand2 != jOut {
				oops = append(oops, jOut)
			}

			if cOut != "" && i != bits-1 && cOut[0] == "z"[0] {
				oops = append(oops, cOut)
				carry_output = ""
			} else {
				carry_output = cOut
			}
		}
	}

	sort.Strings(oops)
	return strings.Join(oops, ",")
}
