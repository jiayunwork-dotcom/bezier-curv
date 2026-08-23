package cli

import (
	"fmt"
	"strconv"
	"strings"
)

// optSchema 声明某个子命令接受的选项：选项名 → 是否带值。
type optSchema map[string]bool

// argSet 是解析后的命令行参数。
type argSet struct {
	values map[string][]string
	pos    []string
}

// parseArgs 解析命令行参数。选项（-x 或 --x）与位置参数可任意交错，
// 负数值（如 -d -0.3）能正确作为选项值读取。
// 选项取值形式：-d 0.25、-d=0.25、--d 0.25 均可；无值选项视为布尔开关。
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

// isOptionToken 判断 token 是否是选项（以 - 开头且不是单独的 "-"）。
func isOptionToken(tok string) bool {
	return len(tok) > 1 && tok[0] == '-'
}

// splitOption 把选项 token 拆成名字与内联值（-d=0.25 → "d", "0.25", true）。
func splitOption(tok string) (name, val string, hasVal bool) {
	t := strings.TrimLeft(tok, "-")
	if eq := strings.IndexByte(t, '='); eq >= 0 {
		return t[:eq], t[eq+1:], true
	}
	return t, "", false
}

// get 返回选项的最后一个取值。
func (a *argSet) get(name string) (string, bool) {
	vs := a.values[name]
	if len(vs) == 0 {
		return "", false
	}
	return vs[len(vs)-1], true
}

// all 返回选项的全部取值（用于可重复的 -t）。
func (a *argSet) all(name string) []string { return a.values[name] }

// float64 读取浮点选项。
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

// int 读取整数选项。
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

// floats 读取可重复浮点选项。
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

// bool 读取布尔开关。
func (a *argSet) bool(name string, def bool) bool {
	if len(a.values[name]) > 0 {
		return true
	}
	return def
}
