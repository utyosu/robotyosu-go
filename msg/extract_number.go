package msg

import (
	"regexp"
	"strconv"
)

var (
	regexpNumber = regexp.MustCompile(`[\d]+`)
)

func ExtractNumber(message string) int {
	numbers := regexpNumber.FindAllStringSubmatch(message, -1)
	if len(numbers) == 1 {
		number, _ := strconv.Atoi(numbers[0][0])
		return number
	}
	return 0
}
