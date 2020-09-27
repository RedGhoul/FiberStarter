package providers

import (
	hashing "github.com/thomasvvugt/fiber-hashing"
)

var hash hashing.Driver

func HashProvider() hashing.Driver {
	return hash
}

func SetUpHashProvider() {
	hash = hashing.New()
}
