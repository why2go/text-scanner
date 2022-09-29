package processor

import (
	"fmt"
	"testing"
)

func TestTrim(t *testing.T) {
	tp := NewTrimProcessor()

	tp.AddIgnoredText(" ")
	tp.AddIgnoredText("$")
	tp.AddIgnoredText("🏴󠁧󠁢󠁥󠁮󠁧󠁿")
	tp.AddIgnoredText("🚴")
	tp.AddIgnoredText("🚴🏻")
	tp.AddIgnoredText("🚴‍♂️")
	tp.AddIgnoredText("🚴🏻‍♂️")

	fmt.Printf("tp.Trim([]rune(\"f u c k\")): %v\n", tp.Trim([]rune("f u c k")))
	fmt.Println()
	fmt.Printf("tp.Trim([]rune(\"f   \")): %v\n", tp.Trim([]rune("f   "))) // 结尾的字符没有记录
	fmt.Println()
	fmt.Printf("tp.Trim([]rune(\"f c k\")): %v\n", tp.Trim([]rune("f c k")))
	fmt.Println()
	fmt.Printf("tp.Trim([]rune(\"f🏴󠁧󠁢󠁥󠁮󠁧󠁿u🏴󠁧󠁢󠁥󠁮󠁧󠁿c🏴󠁧󠁢󠁥󠁮󠁧󠁿k\")): %v\n", tp.Trim([]rune("f🏴󠁧󠁢󠁥󠁮󠁧󠁿u🏴󠁧󠁢󠁥󠁮󠁧󠁿c🏴󠁧󠁢󠁥󠁮󠁧󠁿k")))
	fmt.Println()
	fmt.Printf("tp.Trim([]rune(\"f$u$$c$$ $$k$$$\")): %v\n", tp.Trim([]rune("f$u$$c$$ $$k$$$")))
	fmt.Println()
	fmt.Printf("tp.Trim(\"$$$$$  🏴󠁧󠁢󠁥󠁮󠁧󠁿 $$$\"): %v\n", tp.Trim([]rune("$$$$$  🏴󠁧󠁢󠁥󠁮󠁧󠁿 $$$")))
	fmt.Println()
	fmt.Printf("tp.Trim([]rune(\"f🚴u🚴‍♂️c🚴🏻‍♂️k$\")): %v\n", tp.Trim([]rune("f🚴u🚴‍♂️c🚴🏻‍♂️k$")))
}
