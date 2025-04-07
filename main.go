package main

import (
	"fmt"
	"moon-cost/auth"
)

func main() {
	rs := auth.RandomSalt{
		Length: 64,
	}

	salt := rs.Generate()

	fmt.Printf("salt: %v\n", salt)
}
