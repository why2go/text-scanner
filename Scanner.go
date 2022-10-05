package scanner

import (
	"errors"
	"fmt"
)

var (
	ErrTextTooLong = errors.New("text too long")
)

// Scanner用来扫描需要审核的文本
type Scanner interface {
	Scan(text string) (ScanResult, error)
}

type ScanResult struct {
	FilteredText string
	// Label        string
	// Suggestion   string
}

// 找出文本中的模式串
type Matcher interface {
	FindMatches(rtext []rune) []Match
}

type Match struct {
	S int
	E int
}

func (m Match) String() string {
	return fmt.Sprintf("[%d, %d)", m.S, m.E)
}

// 文本预处理
type PreProcessor interface {
	Process(rtext []rune) PreProcessResult
}

type PreProcessResult struct {
	TrimText      []rune
	OriginalIndex []int // 对应TrimIndex的每个rune在原字符串中的下标索引
	Clips         []Clip
}

func (tr PreProcessResult) String() string {
	return fmt.Sprintf("{%s, %v, %v}", string(tr.TrimText), tr.Clips, tr.OriginalIndex)
}

type Clip struct {
	Text []rune
	S    int
	E    int
}

func (c Clip) String() string {
	return fmt.Sprintf("{%s, [%d, %d)}", string(c.Text), c.S, c.E)
}

// 判断token是否合法
type TokenValidator interface {
	IsValidToken(text []rune, match Match) bool
}
