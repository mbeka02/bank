package utils

import (
	"fmt"
	"math/rand"
	"strings"
	//"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

/*func Initialize() *rand.Rand {
	//create a new random number generator  with a custom seed
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	return rng
}*/

func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // return random int btwn min and max
}

func RandString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		character := alphabet[rand.Intn(k)]
		sb.WriteByte(character)

	}
	return sb.String()

}

func RandName() string {

	name := RandString(10)
	return name
}

func RandMoney() int64 {
	money := RandInt(2000, 6000)
	return money
}

func RandCurrency() string {
	curr := []string{"USD", "EUR", "KSH"}
	return curr[rand.Intn(len(curr))]
}

func RandEmail() string {

	return fmt.Sprintf("%s@gmail.com", RandString(10))
}
