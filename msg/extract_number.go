package msg

import (
	"regexp"
	"strconv"
)

type ExtractResult int

const (
	ExtractResultSuccess ExtractResult = iota
	ExtractResultNotFoundNumber
	ExtractResultMultipleNumber
)

var (
	regexpNumber = regexp.MustCompile(`[\d]+`)
)

func ExtractNumber(message string) (ExtractResult, int) {
	numbers := regexpNumber.FindAllStringSubmatch(message, -1)
	if len(numbers) >= 2 {
		return ExtractResultMultipleNumber, 0
	} else if len(numbers) == 1 {
		number, _ := strconv.Atoi(numbers[0][0])
		return ExtractResultSuccess, number
	}
	return ExtractResultNotFoundNumber, 0
}
