package actrie

import (
	"container/list"
	"sync"

	scanner "gitee.com/piecat/text-scanner"
)

// 实现 Aho-Corasick automaton
// 1. construct the trie,
// 2. construct suffix links, and
// 3. construct output links.

// goto functin: g
// failure functin: f
// output function: output

func NewACTrie() *ACTrie {
	t := &ACTrie{
		root: newNode(),
		size: 0,
		rwmu: &sync.RWMutex{},
	}
	return t
}

type ACTrie struct {
	root *node
	size uint32
	rwmu *sync.RWMutex
}

func (t *ACTrie) Put(rkey []rune) {
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
			nextNode.cp = cp
			curNode.next[cp] = nextNode
		}
		curNode = nextNode
	}
	if !curNode.isTerminal {
		t.size++
		curNode.isTerminal = true
		curNode.suffixLengths = append(curNode.suffixLengths, len(rkey))
	}
}

func (t *ACTrie) FindMatches(rtext []rune) []scanner.Match {
	if len(rtext) == 0 {
		return nil
	}
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	state := t.root
	var matches []scanner.Match
	for i := range rtext {
		cp := rtext[i]
		for _, ok := state.next[cp]; !ok; _, ok = state.next[cp] {
			if state == t.root {
				break
			}
			state = state.failure
		}
		state = state.next[cp]
		if state != nil {
			for _, l := range state.suffixLengths {
				matches = append(matches, scanner.Match{S: i - l + 1, E: i + 1})
			}
		} else {
			state = t.root
		}
	}
	return matches
}

func (t *ACTrie) Contains(rkey []rune) bool {
	if len(rkey) == 0 {
		return false
	}
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	curNode := t.root
	for i := range rkey {
		nextNode, ok := curNode.next[rkey[i]]
		if !ok {
			return false
		}
		curNode = nextNode
	}
	return curNode.isTerminal
}

// 构建failure link
func (t *ACTrie) ConstructFailureLinks() {
	t.rwmu.Lock()
	defer t.rwmu.Unlock()
	if t.size == 0 {
		return
	}
	queue := list.New()
	for _, child := range t.root.next {
		child.failure = t.root
		queue.PushBack(child)
	}
	for queue.Len() > 0 {
		curNode := queue.Remove(queue.Front()).(*node)
		for _, nextNode := range curNode.next {
			curFailure := curNode.failure
			cp := nextNode.cp
			for {
				nextFailure, ok := curFailure.next[cp]
				if ok {
					nextNode.failure = nextFailure
					nextNode.suffixLengths = append(nextNode.suffixLengths, nextFailure.suffixLengths...)
					break
				} else {
					if curFailure == t.root {
						nextNode.failure = t.root
						break
					}
					curFailure = curFailure.failure
				}
			}
			queue.PushBack(nextNode)
		}
	}
}

// 判断前缀树是否包含任何key
func (t *ACTrie) IsEmpty() bool {
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	return t.size == 0
}

// 返回当前前缀树中包含的key的数量
func (t *ACTrie) Size() uint32 {
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	return t.size
}

// 前缀树中的节点类型
type node struct {
	cp            rune
	isTerminal    bool
	next          map[rune]*node
	failure       *node
	suffixLengths []int // 用于找到matches
}

func newNode() *node {
	n := &node{
		isTerminal: false,
		next:       make(map[rune]*node, 0),
	}
	return n
}
