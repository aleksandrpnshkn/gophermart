package services

import (
	"strconv"
)

func IsValidLuhnNumber(n string) bool {
	sum := 0
	checkMod := len(n) % 2

	for idx, ds := range n {
		d, err := strconv.Atoi(string(ds))
		if err != nil {
			return false
		}

		if idx%2 != checkMod {
			sum += d
			continue
		}

		d *= 2
		if d > 9 {
			d -= 9
		}
		sum += d
	}

	return sum%10 == 0
}
