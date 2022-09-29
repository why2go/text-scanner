package util

import (
	"sort"

	scanner "gitee.com/piecat/text-scanner"
)

func MergeOrderedMatches(matches []scanner.Match) []scanner.Match {
	return mergeMatches(matches, false)
}

func MergeUnorderedMatches(matches []scanner.Match) []scanner.Match {
	return mergeMatches(matches, true)
}

func mergeMatches(matches []scanner.Match, needSort bool) []scanner.Match {
	if len(matches) == 0 || len(matches) == 1 {
		return matches
	}
	if needSort {
		sort.Slice(matches[:], func(i, j int) bool { return matches[i].S < matches[j].S })
	}
	var mergedMatches []scanner.Match
	s0 := matches[0].S
	e0 := matches[0].E
	for i := 1; i < len(matches); i++ {
		s1 := matches[i].S
		e1 := matches[i].E
		if s1 <= e0 {
			if e0 < e1 {
				e0 = e1
			}
		} else {
			mergedMatches = append(mergedMatches, scanner.Match{S: s0, E: e0})
			s0 = s1
			e0 = e1
		}
	}
	mergedMatches = append(mergedMatches, scanner.Match{S: s0, E: e0})
	return mergedMatches
}
