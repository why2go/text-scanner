package triescanner

import (
	scanner "gitee.com/piecat/text-scanner"
	"gitee.com/piecat/text-scanner/internal/trie"
)

// 问题：
// 1. 对于干扰字符，在检查text之前做了去除，要考虑干扰字符集的更新对于屏蔽词库的影响，
// 因为可能添加的屏蔽词中包含干扰字符，比如 f*c*u*k，作为屏蔽词加入了词库，那么一开始是能
// 检测出f*c*u*k的，但是如果，此时把*设置为干扰字符，那么送检的词就变成了fuck，因为屏蔽词不是fuck，则无法检测出fuck

// 2. 书写系统，构词方法

// 3. 词库的更新

type TrieScanner struct {
	keyWordTrie     *trie.Trie
	ignoredCharTrie *trie.Trie
}

func NewTrieScanner() *TrieScanner {
	ts := &TrieScanner{}
	return ts
}

func (ts *TrieScanner) Scan(text string) *scanner.ScanResult {
	// 0. 将string转换为[]rune(以unicode的code point作为文本处理基本单元，更方便依据unicode标准做字符分类), 得到rtext
	// 1. 规范化待检测文本rtext，去除干扰字符，得到nrtext(normolized rune text)和extractResult
	// 2. 将nrtext过一遍keyWordTrie, 得到所有的匹配项nrmatches
	// 3. 结合nrtext, nrmatches以及extractResult, 得到在rtext上的rmatches
	// 4. 对rmatches做进一步判定，确定是否真的需要屏蔽
	// 5. 返回ScanResult
	return nil
}
