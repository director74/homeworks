package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	var topSlice []string

	type textPart struct {
		word string
		cnt  int
	}

	dict := make(map[string]int)

	for _, str := range strings.Fields(text) {
		dict[str] += 1
	}

	result := make([]textPart, 0)

	for word, cnt := range dict {
		result = append(result, textPart{word: word, cnt: cnt})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].cnt > result[j].cnt || result[i].cnt == result[j].cnt && strings.Compare(result[i].word, result[j].word) < 0
	})

	j := len(result)
	for i := 0; i < 10 && j > 0; i++ {
		topSlice = append(topSlice, result[i].word)
		j--
	}

	return topSlice
}
