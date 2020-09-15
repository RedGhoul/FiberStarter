package providers

import (
	"fmt"

	hashing "github.com/thomasvvugt/fiber-hashing"
)

var hash hashing.Driver

func HashProvider() hashing.Driver {
	return hash
}

func SetHashProvider() {
	hash = hashing.New()
	fmt.Println(hash)
}
