package cli

import (
	"fmt"
	"strconv"
	"strings"
)

type optSchema map[string]bool

type argSet struct {
	values map[string][]string
	pos    []string
}

func parseArgs(args []string, schema optSchema) (*argSet, error) {
	a := &argSet{values: map[string][]string{}}
	i := 0
	for i < len(args) {
		tok := args[i]
		if !isOptionToken(tok) {
			a.pos = append(a.pos, tok)
			i++
			continue
		}
		name, inline, hasInline := splitOption(tok)
		if _, known := schema[name]; !known {
			return nil, fmt.Errorf("unknown option %q", tok)
		}
		expectsValue := schema[name]
		switch {
		case hasInline && expectsValue:
			a.values[name] = append(a.values[name], inline)
		case hasInline && !expectsValue:
			return nil, fmt.Errorf("option -%s does not take a value", name)
		case expectsValue:
			if i+1 >= len(args) {
				return nil, fmt.Errorf("option -%s requires a value", name)
			}
			i++
			a.values[name] = append(a.values[name], args[i])
		default:
			a.values[name] = append(a.values[name], "true")
		}
		i++
	}
	return a, nil
}

func isOptionToken(tok string) bool {
	return len(tok) > 1 && tok[0] == '-'
}

func splitOption(tok string) (name, val string, hasVal bool) {
	t := strings.TrimLeft(tok, "-")
	if eq := strings.IndexByte(t, '='); eq >= 0 {
		return t[:eq], t[eq+1:], true
	}
	return t, "", false
}

func (a *argSet) get(name string) (string, bool) {
	vs := a.values[name]
	if len(vs) == 0 {
		return "", false
	}
	return vs[len(vs)-1], true
}

func (a *argSet) all(name string) []string { return a.values[name] }

func (a *argSet) float64(name string, def float64) float64 {
	v, ok := a.get(name)
	if !ok {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return def
	}
	return f
}

func (a *argSet) int(name string, def int) int {
	v, ok := a.get(name)
	if !ok {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func (a *argSet) floats(name string) []float64 {
	vs := a.all(name)
	out := make([]float64, 0, len(vs))
	for _, s := range vs {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			out = append(out, f)
		}
	}
	return out
}

func (a *argSet) bool(name string, def bool) bool {
	if len(a.values[name]) > 0 {
		return true
	}
	return def
}
