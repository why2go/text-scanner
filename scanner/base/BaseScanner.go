package base

import (
	scanner "gitee.com/piecat/text-scanner"
	"gitee.com/piecat/text-scanner/util"
)

type BaseScanner struct {
	patternMatcher       scanner.Matcher
	confusingCharMatcher scanner.Matcher
	processor            scanner.PreProcessor
	validator            scanner.TokenValidator
	maxTextLength        int
}

func (bs *BaseScanner) Scan(text string) (scanner.ScanResult, error) {
	if len(text) == 0 {
		return scanner.ScanResult{FilteredText: text}, nil
	}
	if len(text) > bs.maxTextLength {
		return scanner.ScanResult{}, scanner.ErrTextTooLong
	}
	rtext := []rune(text)
	trimResult := bs.processor.Process(rtext)
	trimMatches := bs.patternMatcher.FindMatches(trimResult.TrimText)
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
		if bs.validator.IsValidToken(rtext, originalMatches[i]) {
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
	return scanResult, nil
}
