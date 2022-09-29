package scanner

import "fmt"

type Scanner interface {
	Scan(text string) ScanResult
}

type ScanResult struct {
	FilteredText string
	Label        string
	Suggestion   string
}

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

type TextProcessor interface {
	Trim(rtext []rune) TrimResult
	RestoreMatches(TrimResult, []Match) []Match
}

type TrimResult struct {
	TrimText []rune
	Clips    []Clip
}

func (tr TrimResult) String() string {
	return fmt.Sprintf("{%s, %v}", string(tr.TrimText), tr.Clips)
}

type Clip struct {
	Text []rune
	S    int
	E    int
}

func (c Clip) String() string {
	return fmt.Sprintf("{%s, [%d, %d)}", string(c.Text), c.S, c.E)
}

type TokenValidator interface {
	IsValidToken(text []rune, match Match) bool
}
