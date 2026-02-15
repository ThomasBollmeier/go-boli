package internal

import (
	"fmt"
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

	if prefix == "" {
		return ac.env.GetNames()
	}
	parts := strings.Fields(prefix)
	lastPart := parts[len(parts)-1]

	names := ac.env.GetNames()
	for _, name := range names {
		if strings.HasPrefix(name, lastPart) {
			completions = append(completions, name)
			fmt.Println(name)
		}
	}

	return completions
}
