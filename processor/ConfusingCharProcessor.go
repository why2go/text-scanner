package processor

import (
	scanner "gitee.com/piecat/text-scanner"
	"gitee.com/piecat/text-scanner/util"
)

// 处理停顿字符（干扰字符）
type ConfusingCharProcessor struct {
	matcher scanner.Matcher
}

func NewConfusingCharProcessor(matcher scanner.Matcher) *ConfusingCharProcessor {
	p := &ConfusingCharProcessor{
		matcher: matcher,
	}
	return p
}

// // 将想要在文本扫描中被忽略的字符加入到字典树中，用于检出这些字符，需要在Proecess前调用
// func (p *ConfusingCharProcessor) AddIgnoredText(text string) {
// 	rtext := []rune(text)
// 	p.confusingTextTrie.Put(rtext)
// }

// 去除原始文本中的所有干扰字符
func (p *ConfusingCharProcessor) Process(rtext []rune) scanner.PreProcessResult {
	matches := p.matcher.FindMatches(rtext)
	matches = util.MergeOrderedMatches(matches)
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
	result := scanner.PreProcessResult{
		TrimText:      trimText,
		OriginalIndex: originalIndex,
		Clips:         clips,
	}
	return result
}
