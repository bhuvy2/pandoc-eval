package main

import (
    "fmt"
    "os"

    "github.com/Azure/golua/lua"
    "github.com/Azure/golua/std"
)

var (
    trace bool = false
    debug bool = false
    tests bool = false
)

func must(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
        os.Exit(1)
    }
}

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
        fmt.Print("MyFunc: ")
        fmt.Print(str)
        state.Pop()
    }
    fmt.Println()
    return 0
}

func Open(state *lua.State) int {
    var baseFuncs = map[string]lua.Func{
        "__lf_tex_write": lua.Func(basePrint),
    }

    // Open base library into globals table.
    state.PushGlobals()
    state.SetFuncs(baseFuncs, 0)

    return 1
}

func main() {
    var opts = []lua.Option{lua.WithTrace(trace), lua.WithVerbose(debug)}
    state := lua.NewState(opts...)
    defer state.Close()
    std.Open(state)
    Open(state)

    patch_startup := `
os.execute = function(...) end
tex = {}
tex.write = __lf_tex_write
`
    security_startup := `
os.execute = function(...) end
`

    state.ExecText(patch_startup)
    state.ExecText(security_startup)
    state.ExecText("tex.write(\"Hello\")")
}
