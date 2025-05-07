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

var descendants map[string][]string

func main() {
	var fileName string

	if len(os.Args) == 2 {
		fileName = os.Args[1]
	} else {
		fileName = "input.txt"
	}

	conns := readInput(fileName)

	descendants := make(map[string][]string)

	for k, v := range conns {
		if v.operand1 == "" {
			continue
		}

		mapAppend(descendants, v.operand1, k)
		mapAppend(descendants, v.operand2, k)
	}

	fmt.Printf("The z wires output %s\n", Part1(conns))
	fmt.Printf("The swapped wires are %v\n", Part2(conns))

	//tc := TestCase{x: 0, y: 0, expectedZ: 0}
	//fmt.Println(conns.Test(tc))
	valid, _ := conns.IsValid()
	fmt.Println(valid)
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

	return strconv.Itoa(copy.GetValue("z").dec)
}

func Part2(conns Connections) string {
	copy := conns.Copy()

	fmt.Println(copy.GetValue("z"))

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

func (conns Connections) Bits() int {
	bits := 0
	for i := range conns {
		if i[0] == "x"[0] {
			bits++
		}
	}

	return bits
}

func (conns Connections) Test(tc TestCase) (bool, []int, Value) {
	copy := conns.Copy()

	copy.SetValue("x", tc.x)
	copy.SetValue("y", tc.y)

	actualZ := copy.GetValue("z")

	fmt.Printf("%d + %d = %d\n", tc.x.dec, tc.y.dec, actualZ.dec)

	if actualZ.dec == tc.expectedZ.dec {
		return true, nil, actualZ
	}

	digits := conns.Bits()

	wrongDigits := []int{}

	for i := range digits {
		if tc.expectedZ.bin[i] != actualZ.bin[i] {
			wrongDigits = append(wrongDigits, i)
		}
	}

	return false, wrongDigits, actualZ
}

func (conns *Connections) SetValue(prefix string, value Value) {
	for i := range value.bin {
		key := fmt.Sprintf("%s%d", prefix, i)
		rawBitValue, _ := strconv.Atoi(value.bin[i:i])
		bitValue := uint8(rawBitValue)

		if (*conns)[key].value == bitValue {
			continue
		}

		tmp := (*conns)[key]
		tmp.value = uint8(value.dec)
		(*conns)[key] = tmp
	}
}

func (conns Connections) GetValue(prefix string) Value {
	binaryDigits := ""

	for _, k := range conns.Keys() {
		if k[0] != prefix[0] {
			continue
		}
		newDigit := strconv.Itoa(int(SettleValue(&conns, k)))
		binaryDigits = fmt.Sprintf("%s%s", newDigit, binaryDigits)
	}

	return Value{
		dec: binToDec(binaryDigits),
		bin: binaryDigits,
	}
}

func (conns Connections) IsValid() (bool, int) {
	bits := conns.Bits()
	var tc TestCase
	var num int

	for _, letter := range []string{"x", "y"} {
		for i := range bits - 1 {
			num = int(math.Pow(float64(2), float64(i)))

			if letter == "x" {
				tc = TestCase{
					x:         valueFromDec(num, bits),
					y:         valueFromDec(0, bits),
					expectedZ: valueFromDec(num, bits),
				}
			} else {
				tc = TestCase{
					x:         valueFromDec(0, bits),
					y:         valueFromDec(num, bits),
					expectedZ: valueFromDec(num, bits),
				}
			}

			pass, _, result := conns.Test(tc)

			if !pass {
				fmt.Println(result)
				return false, i
			}
		}

		num = int(math.Pow(float64(2), float64(bits))) - 1

		if letter == "x" {
			tc = TestCase{
				x:         valueFromDec(num, bits),
				y:         valueFromDec(1, bits),
				expectedZ: valueFromDec(num, bits),
			}
		} else {
			tc = TestCase{
				x:         valueFromDec(1, bits),
				y:         valueFromDec(num, bits),
				expectedZ: valueFromDec(num, bits),
			}
		}

		pass, _, result := conns.Test(tc)

		if !pass {
			fmt.Println(result)
			return false, bits
		}
	}

	return true, 0
}

func (conns *Connections) Swap(a string, b string) {
	tmp := (*conns)[a]
	(*conns)[a] = (*conns)[b]
	(*conns)[b] = tmp
}

type TestCase struct {
	x         Value
	y         Value
	expectedZ Value
}

type Value struct {
	dec int
	bin string
}

func valueFromDec(value int, bits int) Value {
	return Value{
		dec: value,
		bin: decToBin(value, bits),
	}
}

func valueFromBin(value string) Value {
	return Value{
		dec: binToDec(value),
		bin: value,
	}
}

func mapAppend(descMap map[string][]string, key string, value string) {
	if _, ok := descMap[key]; !ok {
		descMap[key] = make([]string, 0)
	}

	if !slices.Contains(descMap[key], value) {
		descMap[key] = append(descMap[key], value)
	}
}
