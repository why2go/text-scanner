package processor

import (
	"fmt"
	"sort"

	scanner "gitee.com/piecat/text-scanner"
	"gitee.com/piecat/text-scanner/internal/trie"
)

type TrimProcessor struct {
	ignoredTextTrie *trie.Trie
}

func NewTrimProcessor() *TrimProcessor {
	p := &TrimProcessor{
		ignoredTextTrie: trie.NewTrie(),
	}
	return p
}

// 将想要在文本扫描中被忽略的字符加入到字典树中，用于检出这些字符
func (p *TrimProcessor) AddIgnoredText(text string) {
	rtext := []rune(text)
	p.ignoredTextTrie.Put(rtext)
}

func (p *TrimProcessor) Trim(rtext []rune) scanner.TrimResult {
	matches := p.ignoredTextTrie.FindMatches(rtext)
	fmt.Printf("origin matches: %v\n", matches)
	matches = p.mergeMatches(matches)
	fmt.Printf("merged matches: %v\n", matches)
	var trimText []rune
	var originalIndex []int
	var clips []scanner.Clip
	k := 0
	for i := range matches {
		s := len(trimText)
		trimText = append(trimText, rtext[k:matches[i].S]...)
		for j := k; j < matches[i].S; j++ {
			originalIndex = append(originalIndex, j)
		}
		k = matches[i].E
		clip := scanner.Clip{
			Text: rtext[matches[i].S:matches[i].E],
			S:    s,
			E:    len(trimText),
		}
		clips = append(clips, clip)
	}
	if k < len(rtext) {
		trimText = append(trimText, rtext[k:]...)
		for i := k; i < len(rtext); i++ {
			originalIndex = append(originalIndex, i)
		}
	}
	result := scanner.TrimResult{
		TrimText:      trimText,
		OriginalIndex: originalIndex,
		Clips:         clips,
	}
	return result
}

func (p *TrimProcessor) mergeMatches(matches []scanner.Match) []scanner.Match {
	if len(matches) == 0 || len(matches) == 1 {
		return matches
	}
	sort.Slice(matches[:], func(i, j int) bool { return matches[i].S < matches[j].S })
	fmt.Printf("sorted matches: %v\n", matches)
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

// 将trim后的文本得到的matches，还原为原始文本的matches
func (p *TrimProcessor) RestoreMatches(trimResult scanner.TrimResult, trimMatches []scanner.Match) []scanner.Match {
	// tp.Trim([]rune("f$u$$c$$ $$k$$$")): {fuck, [{$, [0, 1)} {$$, [1, 2)} {$$ $$, [2, 3)} {$$$, [3, 4)}]}
	clips := trimResult.Clips
	var matches []scanner.Match
	for i := range clips {
		fmt.Printf("i: %v\n", i)
	}
	return matches
}
