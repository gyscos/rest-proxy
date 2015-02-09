package main

import "math/rand"

var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Generates a random token with the given length,
// using only alphanumeric characters.
func randomId(l int) string {
	result := make([]rune, l)

	for i := 0; i < l; i++ {
		result[i] = rune(letters[rand.Intn(len(letters))])
	}

	return string(result)
}
