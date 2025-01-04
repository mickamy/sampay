package ulid

import (
	"crypto/rand"
	"time"

	lib "github.com/oklog/ulid/v2"
)

func New() string {
	entropy := lib.Monotonic(rand.Reader, 0)
	return lib.MustNew(lib.Timestamp(time.Now()), entropy).String()
}
