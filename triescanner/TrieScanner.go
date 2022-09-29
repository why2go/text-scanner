package triescanner

import (
	scanner "gitee.com/piecat/text-scanner"
	"gitee.com/piecat/text-scanner/internal/util"
)

// 问题：
// 1. 对于干扰字符，在检查text之前做了去除，要考虑干扰字符集的更新对于屏蔽词库的影响，
// 因为可能添加的屏蔽词中包含干扰字符，比如 f*c*u*k，作为屏蔽词加入了词库，那么一开始是能
// 检测出f*c*u*k的，但是如果，此时把*设置为干扰字符，那么送检的词就变成了fuck，因为屏蔽词不是fuck，则无法检测出fuck

// 2. 书写系统，构词方法

// 3. 词库的更新

type TrieScanner struct {
	keyWordMatcher scanner.Matcher
	trimProcessor  scanner.TextProcessor
	tokenValidator scanner.TokenValidator
}

func NewTrieScanner(
	keyWordMatcher scanner.Matcher,
	trimProcessor scanner.TextProcessor,
	tokenValidator scanner.TokenValidator,
) *TrieScanner {
	ts := &TrieScanner{
		keyWordMatcher: keyWordMatcher,
		trimProcessor:  trimProcessor,
		tokenValidator: tokenValidator,
	}
	return ts
}

func (ts *TrieScanner) Scan(text string) scanner.ScanResult {
	if len(text) == 0 {
		return scanner.ScanResult{FilteredText: text}
	}
	rtext := []rune(text)
	trimResult := ts.trimProcessor.Trim(rtext)
	trimMatches := ts.keyWordMatcher.FindMatches(trimResult.TrimText)
	// 将trimMatches转换为original matches
	var originalMatches []scanner.Match
	for i := range trimMatches {
		match := scanner.Match{
			S: trimResult.OriginalIndex[trimMatches[i].S],
			E: trimResult.OriginalIndex[trimMatches[i].E],
		}
		originalMatches = append(originalMatches, match)
	}
	// 为避免误判，判断每个match是否能有效
	validMatches := originalMatches[:0]
	for i := range originalMatches {
		if ts.tokenValidator.IsValidToken(rtext, originalMatches[i]) {
			validMatches = append(validMatches, originalMatches[i])
		}
	}
	mergedMatches := util.MergeOrderedMatches(validMatches)
	//  将每个match替换成*
	if len(mergedMatches) > 0 {
		for i := range mergedMatches {
			for j := mergedMatches[i].S; j < mergedMatches[i].E; j++ {
				// 对于emoji，会变成多个*
				rtext[j] = rune('*')
			}
		}
	}
	scanResult := scanner.ScanResult{
		FilteredText: string(rtext),
	}
	return scanResult
}
