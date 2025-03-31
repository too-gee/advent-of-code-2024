package main

import (
	"bufio"
	"fmt"
	"math"
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
	fmt.Printf("The swapped wires are %v\n", Part2(conns))
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
	copy := conns.Copy()

	return strconv.Itoa(copy.DecValue("z"))
}

func Part2(conns Connections) string {
	for i := range conns.Bits() - 1 {

		num := int(math.Pow(float64(2), float64(i)))

		for _, letter := range []string{"x", "y"} {
			var tc TestCase
			if letter == "x" {
				tc = TestCase{x: num, y: 0, expectedZ: num}
			} else {
				tc = TestCase{x: 0, y: num, expectedZ: num}
			}

			actualZ := conns.Test(tc)
			fmt.Printf("Case :: %d + %d = %d ? %v (got %d)\n", tc.x, tc.y, tc.expectedZ, tc.expectedZ == actualZ, actualZ)
		}
	}

	return ""
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

func decToBin(num int, digits int) string {
	binary := strconv.FormatInt(int64(num), 2)
	padded := strings.Repeat("0", digits) + binary
	return padded[len(padded)-digits:]
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

func (conn Connection) Copy() Connection {
	return Connection{
		value:    conn.value,
		operator: conn.operator,
		operand1: conn.operand1,
		operand2: conn.operand2,
	}
}

type Connections map[string]Connection

func (conns Connections) Copy() Connections {
	copy := make(Connections)

	for k, v := range conns {
		copy[k] = v.Copy()
	}

	return copy
}

func (conns Connections) Keys() []string {
	keys := make([]string, 0, len(conns))

	for k := range conns {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func (conns Connections) DecValue(key string) int {
	binaryDigits := ""

	for _, k := range conns.Keys() {
		if k[0] != key[0] {
			continue
		}
		newDigit := strconv.Itoa(int(SettleValue(&conns, k)))
		binaryDigits = fmt.Sprintf("%s%s", newDigit, binaryDigits)
	}

	return binToDec(binaryDigits)
}

func (conns Connections) Bits() int {
	bits := 0
	for i := range conns {
		if i[0] == "x"[0] {
			bits++
		}
	}

	return bits
}

func (conns Connections) Test(tc TestCase) int {
	copy := conns.Copy()

	for _, letter := range []string{"x", "y"} {
		var value int

		if letter == "x" {
			value = tc.x
		} else {
			value = tc.y
		}

		digits := strings.Split(decToBin(value, conns.Bits()), "")
		for i := range digits {
			digitName := fmt.Sprintf("%s%02d", letter, i)
			tmp := copy[digitName]
			intVal, _ := strconv.Atoi(digits[i])
			tmp.value = uint8(intVal)
			copy[digitName] = tmp
		}
	}

	return copy.DecValue("z")
}

func (conns Connections) IsValid() bool {
	bits := conns.Bits()
	var tc TestCase
	var num int

	for _, letter := range []string{"x", "y"} {
		for i := range bits - 1 {
			num = int(math.Pow(float64(2), float64(i)))

			if letter == "x" {
				tc = TestCase{x: num, y: 0, expectedZ: num}
			} else {
				tc = TestCase{x: 0, y: num, expectedZ: num}
			}

			if num != conns.Test(tc) {
				return false
			}
		}

		num = int(math.Pow(float64(2), float64(bits))) - 1

		if letter == "x" {
			tc = TestCase{x: num, y: 1, expectedZ: num + 1}
		} else {
			tc = TestCase{x: 1, y: num, expectedZ: num + 1}
		}

		if num != conns.Test(tc) {
			return false
		}
	}

	return true
}

func (conns *Connections) Swap(a string, b string) {
	tmp := (*conns)[a]
	(*conns)[a] = (*conns)[b]
	(*conns)[b] = tmp
}

type TestCase struct {
	x         int
	y         int
	expectedZ int
}
