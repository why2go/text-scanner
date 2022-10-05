package processor

import (
	"fmt"
	"testing"

	"gitee.com/piecat/text-scanner/matcher/actrie"
)

func TestPreProcessor(t *testing.T) {
	confusingCharMatcher := actrie.NewACTrie()
	confusingCharMatcher.Put([]rune(" "))
	confusingCharMatcher.Put([]rune("$"))
	confusingCharMatcher.Put([]rune("+"))
	confusingCharMatcher.Put([]rune("@"))
	confusingCharMatcher.Put([]rune("#"))
	confusingCharMatcher.Put([]rune("%"))
	confusingCharMatcher.Put([]rune("&"))
	confusingCharMatcher.Put([]rune("*"))
	confusingCharMatcher.Put([]rune("🏴󠁧󠁢󠁥󠁮󠁧󠁿"))
	confusingCharMatcher.Put([]rune("🚴"))
	confusingCharMatcher.Put([]rune("🚴🏻"))
	confusingCharMatcher.Put([]rune("🚴‍♂️"))
	confusingCharMatcher.Put([]rune("🚴🏻‍♂️"))

	confusingCharMatcher.ConstructFailureLinks()

	tp := NewConfusingCharProcessor(confusingCharMatcher)

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
