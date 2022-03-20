package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type wordEntry struct {
	word  string
	count int
}

var re = regexp.MustCompile(`([а-яА-Я-])+`)

func getWords(entries []*wordEntry) []string {
	words := make([]string, 0, 10)
	for _, entry := range entries {
		words = append(words, entry.word)
	}
	return words
}

func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func sortEntries(entries []*wordEntry) []*wordEntry {
	sort.Slice(entries, func(i, j int) bool {
		switch {
		case entries[i].count > entries[j].count:
			return true
		case entries[i].count == entries[j].count:
			return strings.Compare(entries[i].word, entries[j].word) < 0
		default:
			return false
		}
	})
	return entries
}

func Top10(inputString string) []string {
	wordEntriesMap := make(map[string]*wordEntry)
	words := re.FindAllString(strings.ToLower(strings.ReplaceAll(inputString, "- ", " ")), -1)
	for _, word := range words {
		if entry, ok := wordEntriesMap[word]; ok {
			entry.count++
		} else {
			wordEntriesMap[word] = &wordEntry{
				word:  word,
				count: 1,
			}
		}
	}
	entries := make([]*wordEntry, 0, len(wordEntriesMap))

	for _, entry := range wordEntriesMap {
		entries = append(entries, entry)
	}

	sortedEntries := sortEntries(entries)
	endIndex := IfThenElse(len(sortedEntries) >= 10, 10, len(sortedEntries)).(int)
	return getWords(sortedEntries[0:endIndex])
}
