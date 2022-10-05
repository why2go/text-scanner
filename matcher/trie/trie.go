package trie

import (
	"sync"

	scanner "gitee.com/piecat/text-scanner"
)

type Trie struct {
	root *node
	size uint32
	rwmu *sync.RWMutex
}

func NewTrie() *Trie {
	t := &Trie{
		root: newNode(),
		size: 0,
		rwmu: &sync.RWMutex{},
	}
	return t
}

// 将key加入到trie
func (t *Trie) Put(rkey []rune) {
	if len(rkey) == 0 {
		return
	}
	t.rwmu.Lock()
	defer t.rwmu.Unlock()
	curNode := t.root
	for i := range rkey {
		cp := rkey[i]
		nextNode, ok := curNode.next[cp]
		if !ok {
			nextNode = newNode()
			curNode.next[cp] = nextNode
		}
		curNode = nextNode
	}
	if !curNode.isTerminal {
		t.size++
		curNode.isTerminal = true
	}
}

// 找出text中包含的key
// 暴力法寻找文本中的模式匹配项
func (t *Trie) FindMatches(rtext []rune) []scanner.Match {
	if len(rtext) == 0 {
		return nil
	}
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	var matches []scanner.Match = nil
	for i := range rtext {
		curNode := t.root
		for j := i; j < len(rtext); j++ {
			if nextNode, ok := curNode.next[rtext[j]]; ok {
				curNode = nextNode
				if curNode.isTerminal {
					matches = append(matches, scanner.Match{S: i, E: j + 1})
				}
			} else {
				break
			}
		}
	}
	return matches
}

// 跳过干扰字符寻找，考虑干扰字符的话，还是穷尽所有可能最合适
type IgnoredCharSet struct {
	m map[rune]struct{}
}

// FIXME：可以加入emoji吗？或者多个词组
func NewIgnoredCharSet(chars []rune) *IgnoredCharSet {
	set := &IgnoredCharSet{}
	set.m = make(map[rune]struct{})
	for i := range chars {
		set.m[chars[i]] = struct{}{}
	}
	return set
}

func (s *IgnoredCharSet) Contains(r rune) bool {
	_, ok := s.m[r]
	return ok
}

func (s *IgnoredCharSet) Size() int {
	return len(s.m)
}

func (t *Trie) FindMatchesWithIgnoredChars(text string, ignoredCharSet *IgnoredCharSet) []scanner.Match {
	if len(text) == 0 {
		return nil
	}
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	rtext := []rune(text)
	var matches []scanner.Match = nil
	for i := range rtext {
		matches = append(matches, t.findMatchesWithIgnoredChars(t.root, rtext, i, i, ignoredCharSet)...)
	}
	return matches
}

func (t *Trie) findMatchesWithIgnoredChars(
	node *node, rtext []rune, start, idx int, ignoredCharSet *IgnoredCharSet) []scanner.Match {
	if node == nil {
		return nil
	}
	var matches []scanner.Match
	if node.isTerminal {
		matches = append(matches, scanner.Match{S: start, E: idx})
	}
	if idx == len(rtext) {
		return matches
	}
	if ignoredCharSet.Contains(rtext[idx]) {
		matches = append(matches, t.findMatchesWithIgnoredChars(node, rtext, start, idx+1, ignoredCharSet)...)
	}
	if nextNode, ok := node.next[rtext[idx]]; ok {
		matches = append(matches, t.findMatchesWithIgnoredChars(nextNode, rtext, start, idx+1, ignoredCharSet)...)
	}
	return matches
}

func (t *Trie) Contains(key string) bool {
	if len(key) == 0 {
		return false
	}
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	rkey := []rune(key)
	return t.contains(t.root, rkey, 0)
}

func (t *Trie) contains(node *node, rkey []rune, idx int) bool {
	if node == nil {
		return false
	}
	if idx == len(rkey) {
		return node.isTerminal
	}
	b := rkey[idx]
	return t.contains(node.next[b], rkey, idx+1)
}

// 对于AC自动机而言，删除操作并不合适
// 将key从前缀树中删除
func (t *Trie) Delete(key string) {
	if len(key) == 0 {
		return
	}
	t.rwmu.Lock()
	defer t.rwmu.Unlock()
	rkey := []rune(key)
	t.root = t.delete(t.root, rkey, 0)
}

func (t *Trie) delete(node *node, rkey []rune, idx int) *node {
	if node == nil {
		return nil
	}
	if idx == len(rkey) {
		if node.isTerminal {
			node.isTerminal = false
			t.size--
		}
	} else {
		b := rkey[idx]
		nextNode, ok := node.next[b]
		if !ok {
			return node
		}
		node.next[b] = t.delete(nextNode, rkey, idx+1)
		if node.next[b] == nil {
			delete(node.next, b)
		}
	}
	if node.isTerminal || node.hasChildren() {
		return node
	}
	return nil
}

// 清空trie
func (t *Trie) Clear() {
	t.rwmu.Lock()
	defer t.rwmu.Unlock()
	t.root = newNode()
	t.size = 0
}

// 判断前缀树是否包含任何key
func (t *Trie) IsEmpty() bool {
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	return t.size == 0
}

// 返回当前前缀树中包含的key的数量
func (t *Trie) Size() uint32 {
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	return t.size
}

// 前缀树中的节点类型
type node struct {
	isTerminal bool
	next       map[rune]*node
}

func newNode() *node {
	n := &node{
		isTerminal: false,
		next:       make(map[rune]*node, 0),
	}
	return n
}

func (n *node) hasChildren() bool {
	return len(n.next) != 0
}
