package go_test

import (
	"fmt"
	"strings"
)

type Person struct {
	age/* @TEST_ANNOTATION inline comment #1 */ int
	name string /* @TEST_ANNOTATION inline comment #2 */
}

func mix() int {
	// @TEST_ANNOTATION decode the message and clean up after yourself!
	return 0
}

/*
 * @TEST_ANNOTATION drop a star if you know about this code wars challenge
 * Digital Cypher assigns to each letter of the alphabet unique number.
 * Instead of letters in encrypted word we write the corresponding number
 * Then we add to each obtained digit consecutive digits from the key
 * */
func decode(code []byte, key int) string {
	n := len(code)
	msg := make([]byte, n)
	keyStr := strconv.Itoa(key)
	keyLen := len(keyStr)

	for i := 0; i < n; i++ {
		msg[i] = code[i] - keyStr[i%keyLen] + '0' + 'a' - 1
	}

	return string(msg)
}

// This comment should not be parsed since it does not have an annotation
