package rand

import (
	"math/rand"
)

func RandInt(low, high int) int {
	return rand.Intn(high-low) + low
}
