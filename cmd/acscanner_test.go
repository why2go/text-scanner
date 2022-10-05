package main

import (
	"fmt"
	"testing"

	"gitee.com/piecat/text-scanner/matcher/actrie"
	"gitee.com/piecat/text-scanner/processor"
	"gitee.com/piecat/text-scanner/scanner/acscanner"
	"gitee.com/piecat/text-scanner/validator"
)

func TestACScanner(t *testing.T) {
	patternMatcher := actrie.NewACTrie()
	patternMatcher.Put([]rune("fuck"))
	patternMatcher.Put([]rune("习大大"))
	patternMatcher.Put([]rune("av"))
	patternMatcher.Put([]rune("isis"))

	patternMatcher.ConstructFailureLinks()

	confusingCharMatcher := actrie.NewACTrie()
	confusingCharMatcher.Put([]rune(" "))
	confusingCharMatcher.Put([]rune("+"))
	confusingCharMatcher.Put([]rune("@"))
	confusingCharMatcher.Put([]rune("#"))
	confusingCharMatcher.Put([]rune("%"))
	confusingCharMatcher.Put([]rune("&"))
	confusingCharMatcher.Put([]rune("*"))

	confusingCharMatcher.ConstructFailureLinks()

	preProcessor := processor.NewConfusingCharProcessor(confusingCharMatcher)

	tokenValidator := &validator.Validator{}

	myScanner := acscanner.NewACScanner(
		patternMatcher,
		confusingCharMatcher,
		preProcessor,
		tokenValidator,
		10000,
	)

	fmt.Println(myScanner.Scan("av"))
	fmt.Println()
	fmt.Println(myScanner.Scan("fuck"))
	fmt.Println()
	fmt.Println(myScanner.Scan("习大大"))
	fmt.Println()
	fmt.Println(myScanner.Scan("isis"))
	fmt.Println()

	fmt.Println(myScanner.Scan("have"))
	fmt.Println()
	fmt.Println(myScanner.Scan("this is"))
	fmt.Println()
	fmt.Println(myScanner.Scan("习@大#大"))
	fmt.Println()
}
