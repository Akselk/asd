package backupModule

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func Errors(e error) bool {
	if e != nil {
		return true
	}
	return false
}
func ReadInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	var result []int
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return result, err
		}
		result = append(result, x)
	}
	return result, scanner.Err()
}

func ReadInternalQuBackupFile() (internalQueue []int) {

	file, err := os.Open("./InternalBackup.txt")
	if !Errors(err) {
		readQueue := bufio.NewReader(file)
		internalQueue, err = ReadInts(readQueue)
		Errors(err)
		file.Close()
	} else {
		internalQueue = []int{0}
	}
	return
}

func WriteInternalQueueBackupFile(InternalQueue []int) {
	file, err := os.Create("./InternalBackup.txt")
	if Errors(err) {
		fmt.Println("ERROR: Unable to create backupfile -> Check admin rights")
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	for floorNumber := 0; floorNumber < len(InternalQueue); floorNumber++ {
		write.WriteString(strconv.Itoa(InternalQueue[floorNumber]) + "\n")
	}
	write.Flush()
}
