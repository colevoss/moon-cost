package main

import (
	"bytes"
	"fmt"
)

func main() {
	buf := bytes.Buffer{}

	buf.Write([]byte("HELLo"))

	other := make([]byte, 5)

	buf.Read(other)

	fmt.Printf("other: %v\n", string(other))
}
