package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var (
	patternTrim    = regexp.MustCompile(`^[[:punct:]]+|[[:punct:]]+$`)
	patternAllPunc = regexp.MustCompile(`^[[:punct:]]{2,}$`)
)

func Top10(s string) []string {
	tokens := split(s)
	words := getInOrder(tokens)
	return append([]string(nil), words[:min(10, len(words))]...)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func split(s string) map[string]int {
	s = strings.ToLower(s)
	words := map[string]int{}
	for _, word := range strings.Fields(s) {
		if !patternAllPunc.MatchString(word) {
			word = patternTrim.ReplaceAllLiteralString(word, "")
		}
		if word != "" {
			words[word]++
		}
	}
	return words
}

func getInOrder(words map[string]int) []string {
	keys := make([]string, 0, len(words))
	for key := range words {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		result := words[keys[j]] - words[keys[i]]
		if result == 0 {
			result = strings.Compare(keys[i], keys[j])
		}
		return result < 0
	})

	return keys
}
