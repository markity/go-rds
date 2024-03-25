package main

import (
	"fmt"
	"strconv"
)

func main() {
	i, err := strconv.ParseInt("31212381231231239312", 10, 64)
	fmt.Println(i, err)
}
