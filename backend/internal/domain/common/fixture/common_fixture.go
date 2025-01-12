package fixture

import (
	"fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"

	"mickamy.com/sampay/internal/lib/slices"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	Password = "P@ssw0rd"
)

var (
	random *rand.Rand
)

func init() {
	source := rand.NewSource(1)
	random = rand.New(source)

	gofakeit.GlobalFaker = gofakeit.NewFaker(random, true)
}

func Int() int {
	return random.Int()
}

func Intn(n int) int {
	return random.Intn(n)
}

func String(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[random.Intn(len(letters))]
	}
	return string(b)
}

func FixedULID(seed int) string {
	return FormatNumber(seed, 26)
}

func FormatNumber(number int, digits int) string {
	return fmt.Sprintf("%0*d", digits, number)
}

func RandomStringer[T fmt.Stringer](a []T) T {
	m := map[string]T{}
	selected := gofakeit.GlobalFaker.RandomString(slices.Map(a, func(x T) string {
		m[x.String()] = x
		return x.String()
	}))
	return m[selected]
}

func EmailOrSlug() string {
	return gofakeit.RandomString([]string{gofakeit.GlobalFaker.Email(), gofakeit.GlobalFaker.Username()})
}
