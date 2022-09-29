package actrie

import (
	"container/list"
	"sync"
)

// 实现 Aho-Corasick automaton
// 1. construct the trie,
// 2. construct suffix links, and
// 3. construct output links.

// goto functin: g
// failure functin: f
// output function: output

// 前缀树，输入是utf8编码的string，构建时使用rune，便于根据unicode标准来判断字符类型
// 暂时使用读写锁对整个trie做并发控制

func NewTrie() *ACTrie {
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
	// 使用循环来添加
	for i := range rkey {
		cp := rkey[i]
		nextNode, ok := curNode.next[cp]
		if !ok {
			nextNode = newNode()
			nextNode.parent = curNode
			nextNode.cp = cp
			curNode.next[cp] = nextNode
		}
		curNode = nextNode
	}
	if !curNode.isTerminal {
		t.size++
		curNode.isTerminal = true
		curNode.outputs = append(curNode.outputs, rkey)
	}
}

// 找出text中包含的key

type Match struct {
	Start int
	End   int
}

func (t *ACTrie) FindMatches(rtext []rune) []Match {
	if len(rtext) == 0 {
		return nil
	}
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	curNode := t.root
	// 在这里做输出
	for _, cp := range rtext {
		nextNode, ok := curNode.next[cp]
		if ok {
			curNode = nextNode
		} else {
			curNode = curNode.failure
		}
	}

	return nil
}

func (t *ACTrie) Contains(rkey []rune) bool {
	if len(rkey) == 0 {
		return false
	}
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	return t.contains(t.root, rkey, 0)
}

func (t *ACTrie) contains(node *node, rkey []rune, idx int) bool {
	if node == nil {
		return false
	}
	if idx == len(rkey) {
		return node.isTerminal
	}
	b := rkey[idx]
	return t.contains(node.next[b], rkey, idx+1)
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
	cp         rune
	isTerminal bool
	parent     *node
	next       map[rune]*node
	failure    *node
	outputs    [][]rune
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
