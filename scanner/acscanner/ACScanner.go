package acscanner

import (
	"fmt"

	scanner "gitee.com/piecat/text-scanner"
	"gitee.com/piecat/text-scanner/util"
)

const (
	defaultMaxTextLength = int(10000) // bytes
)

type ACScanner struct {
	patternMatcher       scanner.Matcher
	confusingCharMatcher scanner.Matcher
	processor            scanner.PreProcessor
	validator            scanner.TokenValidator
	maxTextLength        int
}

func NewACScanner(
	patternMatcher scanner.Matcher,
	confusingCharMatcher scanner.Matcher,
	processor scanner.PreProcessor,
	validator scanner.TokenValidator,
	maxTextLength int,
) *ACScanner {
	s := &ACScanner{
		patternMatcher:       patternMatcher,
		confusingCharMatcher: confusingCharMatcher,
		processor:            processor,
		validator:            validator,
	}
	if maxTextLength <= 0 {
		maxTextLength = defaultMaxTextLength
	}
	s.maxTextLength = maxTextLength
	return s
}

func (acs *ACScanner) Scan(text string) (scanner.ScanResult, error) {
	if len(text) == 0 {
		return scanner.ScanResult{FilteredText: text}, nil
	}
	if len(text) > acs.maxTextLength {
		return scanner.ScanResult{}, scanner.ErrTextTooLong
	}
	fmt.Printf("originalText: %s\n", text)
	rtext := []rune(text)
	trimResult := acs.processor.Process(rtext)
	fmt.Printf("trimResult: %v\n", trimResult)
	trimMatches := acs.patternMatcher.FindMatches(trimResult.TrimText)
	fmt.Printf("trimMatches: %v\n", trimMatches)
	// 将trimMatches转换为original matches
	var originalMatches []scanner.Match
	for i := range trimMatches {
		match := scanner.Match{
			S: trimResult.OriginalIndex[trimMatches[i].S],
			E: trimResult.OriginalIndex[trimMatches[i].E-1] + 1,
		}
		originalMatches = append(originalMatches, match)
	}
	fmt.Printf("originalMatches: %v\n", originalMatches)
	// 为避免误判，判断每个match是否能有效
	validMatches := originalMatches[:0]
	for i := range originalMatches {
		if acs.validator.IsValidToken(rtext, originalMatches[i]) {
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
