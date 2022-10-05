package processor

import (
	"fmt"
	"testing"
)

func TestPreProcessor(t *testing.T) {
	tp := NewConfusingCharProcessor()

	tp.AddIgnoredText(" ")
	tp.AddIgnoredText("$")
	tp.AddIgnoredText("🏴󠁧󠁢󠁥󠁮󠁧󠁿")
	tp.AddIgnoredText("🚴")
	tp.AddIgnoredText("🚴🏻")
	tp.AddIgnoredText("🚴‍♂️")
	tp.AddIgnoredText("🚴🏻‍♂️")

	fmt.Printf("tp.Process([]rune(\"f u c k\")): %v\n", tp.Process([]rune("f u c k")))
	fmt.Println()
	fmt.Printf("tp.Process([]rune(\"f   \")): %v\n", tp.Process([]rune("f   "))) // 结尾的字符没有记录
	fmt.Println()
	fmt.Printf("tp.Process([]rune(\"f c k\")): %v\n", tp.Process([]rune("f c k")))
	fmt.Println()
	fmt.Printf("tp.Process([]rune(\"f🏴󠁧󠁢󠁥󠁮󠁧󠁿u🏴󠁧󠁢󠁥󠁮󠁧󠁿c🏴󠁧󠁢󠁥󠁮󠁧󠁿k\")): %v\n",
		tp.Process([]rune("f🏴󠁧󠁢󠁥󠁮󠁧󠁿u🏴󠁧󠁢󠁥󠁮󠁧󠁿c🏴󠁧󠁢󠁥󠁮󠁧󠁿k")))
	fmt.Println()
	fmt.Printf("tp.Process([]rune(\"f$u$$c$$ $$k$$$\")): %v\n",
		tp.Process([]rune("f$u$$c$$ $$k$$$")))
	fmt.Println()
	fmt.Printf("tp.Process(\"$$$$$  🏴󠁧󠁢󠁥󠁮󠁧󠁿 $$$\"): %v\n",
		tp.Process([]rune("$$$$$  🏴󠁧󠁢󠁥󠁮󠁧󠁿 $$$")))
	fmt.Println()
	fmt.Printf("tp.Process([]rune(\"f🚴u🚴‍♂️c🚴🏻‍♂️k$\")): %v\n",
		tp.Process([]rune("f🚴u🚴‍♂️c🚴🏻‍♂️k$")))
	fmt.Println()
	fmt.Printf("tp.Process([]rune(\"f👴🏻u👴🏻c👴🏻k👴🏻\")): %v\n",
		tp.Process([]rune("f👴🏻u👴🏻c👴🏻k👴🏻")))
}
