package processor

import (
	scanner "gitee.com/piecat/text-scanner"
	"gitee.com/piecat/text-scanner/internal/trie"
	"gitee.com/piecat/text-scanner/internal/util"
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

// 去除原始文本中的所有干扰字符
func (p *TrimProcessor) Trim(rtext []rune) scanner.TrimResult {
	matches := p.ignoredTextTrie.FindMatches(rtext)
	// fmt.Printf("origin matches: %v\n", matches)
	matches = util.MergeOrderedMatches(matches)
	// fmt.Printf("merged matches: %v\n", matches)
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
