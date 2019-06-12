package main

import (
	"fmt"
	"os"
	"bytes"
	pf "github.com/oltolm/go-pandocfilters"

	"github.com/Azure/golua/lua"
	"github.com/Azure/golua/std"
)

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

var execBuffer bytes.Buffer
var contextState *lua.State

func basePrint(state *lua.State) int {
	var (
		n = state.Top()
		i = 1
	)
	state.GetGlobal("tostring")
	for ; i <= n; i++ {
		state.PushIndex(-1)
		state.PushIndex(i)
		state.Call(1, 1)
		str, ok := state.TryString(-1)
		if !ok {
			panic(fmt.Errorf("'tostring' must return a string to 'print'"))
		}
		execBuffer.WriteString(str)
		state.Pop()
	}
	return 0
}

func addLuatexFuncs(state *lua.State) {
	var baseFuncs = map[string]lua.Func{
		"__lf_tex_write": lua.Func(basePrint),
	}

	// Open base library into globals table.
	state.PushGlobals()
	state.SetFuncs(baseFuncs, 0)
}

func execLuaWithCapture(code string) string {
	execBuffer.Reset()
	contextState.ExecText(code)
	return execBuffer.String()
}

func Caps(key string, value interface{}, format string, meta interface{}) interface{} {
	//fmt.Fprintf(os.Stderr, "%v %v\n", key, value)
	if key == "CodeBlock" || key == "Code"{
		blockType := value.([]interface{})[0].([]interface{})[1].([]interface{})[0].(string)
		//fmt.Fprintf(os.Stderr, "%v\n", blockType)
		if (blockType == "evallua") {
			code := value.([]interface{})[1].(string)
			out := execLuaWithCapture(code)
			if key == "CodeBlock" {
				return pf.RawBlock(format, out)
			}
			return pf.RawInline(format, out)
		}
	}
	return nil
}

func main() {
	var opts = []lua.Option{lua.WithTrace(false), lua.WithVerbose(false)}
	contextState = lua.NewState(opts...)
	defer contextState.Close()
	std.Open(contextState)
	addLuatexFuncs(contextState)

	patch_startup := `
print = __lf_tex_write
`
	security_startup := `
os.execute = function(...) end
`
	contextState.ExecText(patch_startup)
	contextState.ExecText(security_startup)
	pf.ToJSONFilter(Caps)
}
