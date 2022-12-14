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
	confusingCharMatcher.Put([]rune("ğ´ó §ó ¢ó ¥ó ®ó §ó ¿"))
	confusingCharMatcher.Put([]rune("ğ´"))
	confusingCharMatcher.Put([]rune("ğ´ğ»"))
	confusingCharMatcher.Put([]rune("ğ´ââï¸"))
	confusingCharMatcher.Put([]rune("ğ´ğ»ââï¸"))

	confusingCharMatcher.ConstructFailureLinks()

	tp := NewConfusingCharProcessor(confusingCharMatcher)

	fmt.Printf("tp.Process([]rune(\"f u c k\")): %v\n", tp.Process([]rune("f u c k")))
	fmt.Println()
	fmt.Printf("tp.Process([]rune(\"f   \")): %v\n", tp.Process([]rune("f   "))) // ç»å°¾çå­ç¬¦æ²¡æè®°å½
	fmt.Println()
	fmt.Printf("tp.Process([]rune(\"f c k\")): %v\n", tp.Process([]rune("f c k")))
	fmt.Println()
	fmt.Printf("tp.Process([]rune(\"fğ´ó §ó ¢ó ¥ó ®ó §ó ¿uğ´ó §ó ¢ó ¥ó ®ó §ó ¿cğ´ó §ó ¢ó ¥ó ®ó §ó ¿k\")): %v\n",
		tp.Process([]rune("fğ´ó §ó ¢ó ¥ó ®ó §ó ¿uğ´ó §ó ¢ó ¥ó ®ó §ó ¿cğ´ó §ó ¢ó ¥ó ®ó §ó ¿k")))
	fmt.Println()
	fmt.Printf("tp.Process([]rune(\"f$u$$c$$ $$k$$$\")): %v\n",
		tp.Process([]rune("f$u$$c$$ $$k$$$")))
	fmt.Println()
	fmt.Printf("tp.Process(\"$$$$$  ğ´ó §ó ¢ó ¥ó ®ó §ó ¿ $$$\"): %v\n",
		tp.Process([]rune("$$$$$  ğ´ó §ó ¢ó ¥ó ®ó §ó ¿ $$$")))
	fmt.Println()
	fmt.Printf("tp.Process([]rune(\"fğ´uğ´ââï¸cğ´ğ»ââï¸k$\")): %v\n",
		tp.Process([]rune("fğ´uğ´ââï¸cğ´ğ»ââï¸k$")))
	fmt.Println()
	fmt.Printf("tp.Process([]rune(\"fğ´ğ»uğ´ğ»cğ´ğ»kğ´ğ»\")): %v\n",
		tp.Process([]rune("fğ´ğ»uğ´ğ»cğ´ğ»kğ´ğ»")))
}
