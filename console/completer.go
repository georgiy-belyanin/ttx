package console

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/abiosoft/readline/runes"
	"github.com/tarantool/go-tarantool/v2"
)

// This is a patched version of https://github.com/chzyer/readline/blob/main/complete_helper.go.
//
// Changes:
// * Don't add ' ' at the end of suggested completions.
// * Suggest completions in the middle of the code.

type DynamicCompleteFunc func(string) []string

type PrefixCompleterInterface interface {
	Print(prefix string, level int, buf *bytes.Buffer)
	Do(line []rune, pos int) (newLine [][]rune, length int)
	GetName() []rune
	GetChildren() []PrefixCompleterInterface
	SetChildren(children []PrefixCompleterInterface)
}

type DynamicPrefixCompleterInterface interface {
	PrefixCompleterInterface
	IsDynamic() bool
	GetDynamicNames(line []rune) [][]rune
}

type TarantoolCompleter struct {
	Name     []rune
	Dynamic  bool
	Callback DynamicCompleteFunc
	Children []PrefixCompleterInterface
}

func (p *TarantoolCompleter) Tree(prefix string) string {
	buf := bytes.NewBuffer(nil)
	p.Print(prefix, 0, buf)
	return buf.String()
}

func Print(p PrefixCompleterInterface, prefix string, level int, buf *bytes.Buffer) {
	if strings.TrimSpace(string(p.GetName())) != "" {
		buf.WriteString(prefix)
		if level > 0 {
			buf.WriteString("├")
			buf.WriteString(strings.Repeat("─", (level*4)-2))
			buf.WriteString(" ")
		}
		buf.WriteString(string(p.GetName()) + "\n")
		level++
	}
	for _, ch := range p.GetChildren() {
		ch.Print(prefix, level, buf)
	}
}

func (p *TarantoolCompleter) Print(prefix string, level int, buf *bytes.Buffer) {
	Print(p, prefix, level, buf)
}

func (p *TarantoolCompleter) IsDynamic() bool {
	return p.Dynamic
}

func (p *TarantoolCompleter) GetName() []rune {
	return p.Name
}

func (p *TarantoolCompleter) GetDynamicNames(line []rune) [][]rune {
	var names = [][]rune{}
	for _, name := range p.Callback(string(line)) {
		names = append(names, []rune(name))
	}
	return names
}

func (p *TarantoolCompleter) GetChildren() []PrefixCompleterInterface {
	return p.Children
}

func (p *TarantoolCompleter) SetChildren(children []PrefixCompleterInterface) {
	p.Children = children
}

func PcItem(name string, pc ...PrefixCompleterInterface) *TarantoolCompleter {
	name += " "
	return &TarantoolCompleter{
		Name:     []rune(name),
		Dynamic:  false,
		Children: pc,
	}
}

func PcItemDynamic(callback DynamicCompleteFunc, pc ...PrefixCompleterInterface) *TarantoolCompleter {
	return &TarantoolCompleter{
		Callback: callback,
		Dynamic:  true,
		Children: pc,
	}
}

func (p *TarantoolCompleter) Do(line []rune, pos int) (newLine [][]rune, offset int) {
	return doInternal(p, line, pos, line)
}

func Do(p PrefixCompleterInterface, line []rune, pos int) (newLine [][]rune, offset int) {
	return doInternal(p, line, pos, line)
}

func doInternal(p PrefixCompleterInterface, line []rune, pos int, origLine []rune) (newLine [][]rune, offset int) {
	lpos := max(pos, 0)
	// TODO: this should be a little bit smarter and end the word capture not
	// only on spaces but on other symbols like ')' too.
	for lpos > 0 && !unicode.IsSpace(line[lpos-1]) {
		lpos--
	}

	line = line[lpos:pos]
	goNext := false
	var lineCompleter PrefixCompleterInterface
	for _, child := range p.GetChildren() {
		childNames := make([][]rune, 1)

		childDynamic, ok := child.(DynamicPrefixCompleterInterface)
		if ok && childDynamic.IsDynamic() {
			childNames = childDynamic.GetDynamicNames(line)
		} else {
			childNames[0] = child.GetName()
		}

		for _, childName := range childNames {
			if runes.Equal(childName, line) {
				continue
			}
			newLine = append(newLine, childName[len(line):])
			offset = len(line)
			lineCompleter = child
		}
	}

	if len(newLine) != 1 {
		return
	}

	tmpLine := make([]rune, 0, len(line))
	for i := offset; i < len(line); i++ {
		if line[i] == ' ' {
			continue
		}

		tmpLine = append(tmpLine, line[i:]...)
		return doInternal(lineCompleter, tmpLine, len(tmpLine), origLine)
	}

	if goNext {
		return doInternal(lineCompleter, nil, 0, origLine)
	}
	return
}

func completionLua(line string) string {
	var completionLuaBase = `return unpack(require('console').completion_handler([[%s]], 0, %d))`
	return fmt.Sprintf(completionLuaBase, line, len(line))
}

func CreateCompleter(conn *tarantool.Connection) *TarantoolCompleter {
	return PcItem("", PcItemDynamic(func(text string) []string {
		var res []string
		err := conn.Do(tarantool.NewEvalRequest(completionLua(text))).GetTyped(&res)
		if err != nil {
			return []string{}
		}
		return res
	}))
}
