package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
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

	disk := readInput(fileName)

	// part 1
	defragged := disk.expand().defrag()
	fmt.Printf("Checksum after standard defrag: %v\n", defragged.checksum())

	// part 2
	smartDefragged := disk.expand().smartDefrag()
	fmt.Printf("Checksum after smart defrag: %v\n", smartDefragged.checksum())
}

func readInput(filePath string) Disk {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s", filePath)
		return nil
	}
	defer file.Close()

	disk := Disk{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")

		disk = append(disk, row...)
	}

	return disk
}

type Disk []string

func (d Disk) isEmpty(i int) bool {
	return d[i] == "."
}

func (d Disk) draw() {
	for _, value := range d {
		fmt.Print(value)
	}
	fmt.Println()
}

func (d Disk) firstRun(length int) int {
	counter := 0

	for i := range d {
		if d.isEmpty(i) {
			counter += 1
		} else {
			counter = 0
		}

		if counter == length {
			return i - length + 1
		}
	}

	return -1
}

func (d Disk) checksum() int {
	checksum := 0

	for i, value := range d {
		intValue, _ := strconv.Atoi(value)
		checksum += i * intValue
	}

	return checksum
}

func (d Disk) expand() Disk {
	newDisk := Disk{}

	fileId := 0

	for i, value := range d {
		block := ""
		blockCount, _ := strconv.Atoi(value)

		if i%2 == 0 {
			block = strconv.Itoa(fileId)
			fileId += 1
		} else {
			block = "."
		}

		for j := 0; j < blockCount; j++ {
			newDisk = append(newDisk, block)
		}
	}

	return newDisk
}

func (d Disk) defrag() Disk {
	for i := len(d) - 1; i >= 0; i-- {
		// Don't move free space
		if d[i] == "." {
			continue
		}

		freeSpot := d.firstRun(1)

		// free space is ahead of file
		if freeSpot >= i {
			break
		}

		d[freeSpot] = d[i]
		d[i] = "."
	}

	return d
}

func (d Disk) getFileLocation(fileId int) (int, int) {
	start := float64(len(d) + 1)
	end := float64(-1)

	fileIdStr := strconv.Itoa(fileId)

	for i, value := range d {
		if value == fileIdStr {
			start = math.Min(start, float64(i))
			end = math.Max(end, float64(i))
		}
	}

	return int(start), int(end - start + 1)
}

func (d Disk) smartDefrag() Disk {
	fileId := -1
	// Get largest file id
	for i := len(d) - 1; i >= 0; i-- {
		if d[i] != "." {
			fileId, _ = strconv.Atoi(d[i])
			break
		}
	}

	for ; fileId >= 0; fileId-- {
		fileStart, length := d.getFileLocation(fileId)
		freeStart := d.firstRun(length)

		// not enough free space
		if freeStart == -1 {
			continue
		}

		// free space is ahead of file
		if freeStart > fileStart {
			continue
		}

		for i := 0; i < length; i++ {
			d[fileStart+i] = "."
			d[freeStart+i] = strconv.Itoa(fileId)
		}
	}

	return d
}
