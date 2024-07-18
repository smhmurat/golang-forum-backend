package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func MergeLowercaseAndAddRandomNumber(text string) string {
	// Split the text into words by spaces
	words := strings.Fields(text)
	// Join the words and convert to lowercase
	mergedText := strings.ToLower(strings.Join(words, ""))

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	// Generate a random 4-digit number
	randomNumber := fmt.Sprintf("%04d", rand.Intn(10000))

	// Combine the merged text with the random number
	return mergedText + randomNumber
}
