package main

import (
	"fmt"
	"strconv"
)

func validateToken(token string) (int, error) {
	id, err := strconv.Atoi(token)

	if err != nil {
		fmt.Println(token)
		fmt.Println(err)
	}

	return id, err
}
