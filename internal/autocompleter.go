package internal

import (
	ip "go-boli/internal/interpreter"
	"strings"

	"github.com/chzyer/readline"
)

type AutoCompleter struct {
	env *ip.Environment
}

func NewAutoCompleter(env *ip.Environment) *AutoCompleter {
	return &AutoCompleter{env: env}
}

func (ac *AutoCompleter) NewPrefixCompleter() *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItemDynamic(ac.getCompletions),
	)
}

func (ac *AutoCompleter) getCompletions(prefix string) []string {
	var completions []string

	startPos := findStartPos(prefix)
	var firstPart, lastPart string

	if startPos != -1 {
		firstPart = prefix[:startPos]
		lastPart = prefix[startPos:]
	} else {
		firstPart = prefix
		lastPart = ""
	}

	names := ac.env.GetNames()
	for _, name := range names {
		if strings.HasPrefix(name, lastPart) {
			completions = append(completions, firstPart+name)
		}
	}

	return completions
}

func findStartPos(prefix string) int {
	toIgnore := map[rune]bool{
		' ':  true,
		'\t': true,
		'\n': true,
		'(':  true,
		')':  true,
		'[':  true,
		']':  true,
		'{':  true,
		'}':  true,
	}

	prevToIgnore := true
	ret := -1

	for i, ch := range prefix {
		if toIgnore[ch] {
			prevToIgnore = true
			continue
		}
		if prevToIgnore {
			ret = i
			prevToIgnore = false
		}
	}

	return ret
}
