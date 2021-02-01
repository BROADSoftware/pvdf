package cmd

import (
	"testing"
	"fmt"
)

func Test1(t *testing.T) {
	fmt.Printf("Test1\n")
	for x := int64(1); x < 10000000000000; x = x * 10 {
		fmt.Printf("%d -> %s\n", x, bytes2human(x))
		x = x + x/2
		fmt.Printf("%d -> %s\n", x, bytes2human(x))
	}
}
