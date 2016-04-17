package helpModule

func FindIndexSmallestNumber(arrayofnumber []int) (indexOfSmallest int) {

	var number int = 1000
	indexOfSmallest = -1

	for index := 0; index < len(arrayofnumber); index++ {
		if arrayofnumber[index] < number && arrayofnumber[index] > 0 {
			number = arrayofnumber[index]
			indexOfSmallest = index

		}
	}
	return
}

func FindIndexOfNumber(arrayofnumber []int, number int) (indexofnumber int) {
	indexofnumber = -1
	for index := 0; index < len(arrayofnumber); index++ {
		if arrayofnumber[index] == number {
			indexofnumber = index
			return
		}
	}
	return
}

func FindBiggestNumber(arrayofnumber []int) (biggestNr int) {

	for index := 0; index < len(arrayofnumber); index++ {
		if arrayofnumber[index] > biggestNr {
			biggestNr = arrayofnumber[index]

		}
	}
	return
}

func FindSmallestNumber(arrayofnumber []int) (smallestNr int) {

	for index := 0; index < len(arrayofnumber); index++ {
		if arrayofnumber[index] < smallestNr && arrayofnumber[index] != 0 {
			smallestNr = arrayofnumber[index]

		}
	}
	return
}

func IndexToFloor(index int) (floor int) {
	if index > 4 {
		floor = index - 4
	} else {
		floor = index
	}
	return
}

func HighestIndex(arrayofnumber []int) (indexOfNumber int) {
	indexOfNumber = -1

	for index := len(arrayofnumber); index > 0; index-- {
		if arrayofnumber[index-1] > 0 {
			indexOfNumber = index - 1
			return
		}
	}
	return
}

func LowestIndex(arrayofnumber []int) (indexOfNumber int) {
	indexOfNumber = -1
	for index := 0; index < len(arrayofnumber); index++ {
		if arrayofnumber[index] > 0 {
			indexOfNumber = index
			return
		}
	}
	return
}
func Abs(number int) (absnumber int) {
	if number < 0 {
		absnumber = number * -1
	} else {
		absnumber = number
	}
	return
}

func FindLowestIPindex(ip []string) (lowestip int) {
	if len(ip) == 1 {
		return 0
	}
	lowestip = 0
	for index := 0; index < len(ip); index++ {
		if ip[index] < ip[lowestip] {
			lowestip = index
		}
	}
	return

}

func IsEmpty(ExternalQueue []int, InternalQueue []int) bool {
	if FindBiggestNumber(ExternalQueue)+
		FindBiggestNumber(InternalQueue) == 0 {
		return true
	}
	return false
}

func FindIndexSmallestNumberLargerThan(arrayofnumber []int, largenumber int) (smallestNumberIndex int) {
	smallestNr := 32767
	for index := 0; index < len(arrayofnumber); index++ {
		if arrayofnumber[index] > largenumber && arrayofnumber[index] > -1 {
			if arrayofnumber[index] < smallestNr && arrayofnumber[index] != 0 {
				smallestNr = arrayofnumber[index]
			}
		}
	}
	smallestNumberIndex = FindIndexOfNumber(arrayofnumber[:], smallestNr)
	return
}
