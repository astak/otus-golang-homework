package hw03frequencyanalysis

import (
	"slices"
	"strings"
)

func Top10(s string) []string {
	words := map[string]int{}
	for _, word := range strings.Fields(s) {
		words[word]++
	}

	keys := make([]string, 0, len(words))
	for key := range words {
		keys = append(keys, key)
	}

	slices.SortFunc(keys, func(a, b string) int {
		result := words[b] - words[a]
		if result == 0 {
			result = strings.Compare(a, b)
		}
		return result
	})

	return keys[:min(10, len(keys))]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
