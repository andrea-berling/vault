package levenshtein

import (
	"fmt"
	"sort"
	"unicode/utf8"
)

type wordScore struct {
	Word  string
	Score int
}

// BestNPossibilities returns n possibilities where input scores equal to or less than threshold. A threshold of -1 return all matches.
// If no matches are found or n==0 an error is returned.
func BestNPossibilities(input string, possibilities []string, threshold int, n int) ([]wordScore, error) {
	if n == 0 {
		return nil, fmt.Errorf("n should be greater than 0 or set to -1 for all possibilities")
	}

	if n > len(possibilities) || n == -1 {
		n = len(possibilities)
	}

	var guesses []wordScore

	for _, guess := range possibilities {
		score := levenshteinValue(input, guess)
		if score <= threshold || threshold == -1 {
			guesses = append(guesses, wordScore{
				Word:  guess,
				Score: score,
			})
		}
	}

	if len(guesses) == 0 {
		return nil, fmt.Errorf("no matches found using threshold: %d", threshold)
	}

	sort.SliceStable(guesses, func(i, j int) bool {
		return guesses[i].Score < guesses[j].Score
	})

	if len(guesses) > n {
		return guesses[:n], nil
	}

	return guesses, nil
}

func levenshteinValue(a, b string) int {
	f := make([]int, utf8.RuneCountInString(b)+1)

	for j := range f {
		f[j] = j
	}

	for _, ca := range a {
		j := 1
		fj1 := f[0]
		f[0]++
		for _, cb := range b {
			mn := min(f[j]+1, f[j-1]+1)
			if cb != ca {
				mn = min(mn, fj1+1)
			} else {
				mn = min(mn, fj1)
			}

			fj1, f[j] = f[j], mn
			j++
		}
	}

	return f[len(f)-1]
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
