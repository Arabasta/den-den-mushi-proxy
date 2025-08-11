package main

import (
	"fmt"
	"os"
	"strings"

	"mvdan.cc/sh/v3/syntax"
)

func main() {
	cmds := []string{
		"ls",
		" ls",
		"ls -lrt",
		"ls | su",
		"echo ls > foo.txt",
		" ls -lrt | grep -v '^d' | awk '{print $9}'",
		"echo H4sIAIzmgmgAAysuBQB5vaRlAgAAAA== |  base64 -d -d | /usr/bin/gzi* -d | /bin/bas*\n    "}

	for i := range cmds {
		fmt.Printf("Processing command: %s\n", cmds[i])
		c, err := preprocessCommand(cmds[i])
		for i, c := range c {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing command: %v\n", err)
				return
			}

			fmt.Printf("Command %d: %s\n", i+1, strings.Join(c, " "))
		}
	}
}

func preprocessCommand(cmd string) ([][]string, error) {
	r := strings.NewReader(cmd)

	file, err := syntax.NewParser().Parse(r, "")
	if err != nil {
		return nil, err
	}

	var out [][]string
	for _, stmt := range file.Stmts {
		out = append(out, flattenCmd(stmt.Cmd)...)
	}
	return out, nil
}

func flattenCmd(cmd syntax.Command) [][]string {
	switch x := cmd.(type) {
	case *syntax.CallExpr:
		return [][]string{callToArgs(x)}

	case *syntax.BinaryCmd:
		// Covers pipes (|), sequence (;), &&, || — we just flatten L then R
		left := flattenCmd(x.X.Cmd)
		right := flattenCmd(x.Y.Cmd)
		return append(left, right...)

	case *syntax.Subshell:
		// ( ... ); handle inner statements
		var res [][]string
		for _, s := range x.Stmts {
			res = append(res, flattenCmd(s.Cmd)...)
		}
		return res

	case *syntax.Block:
		// { ... }
		var res [][]string
		for _, s := range x.Stmts {
			res = append(res, flattenCmd(s.Cmd)...)
		}
		return res

	default:
		// Other constructs (if/for/while/etc.) not flattened here
		return nil
	}
}

func callToArgs(call *syntax.CallExpr) []string {
	args := make([]string, 0, len(call.Args))
	for _, w := range call.Args {
		args = append(args, wordToString(w))
	}
	return args
}

func wordToString(w *syntax.Word) string {
	var b strings.Builder
	for _, p := range w.Parts {
		switch t := p.(type) {
		case *syntax.Lit:
			b.WriteString(t.Value)
		case *syntax.SglQuoted:
			// 'text' -> text (no expansion)
			b.WriteString(t.Value)
		case *syntax.DblQuoted:
			for _, q := range t.Parts {
				b.WriteString(partToText(q))
			}
		default:
			b.WriteString(partToText(t))
		}
	}
	return b.String()
}

func partToText(n syntax.Node) string {
	switch t := n.(type) {
	case *syntax.Lit:
		return t.Value
	case *syntax.SglQuoted:
		return t.Value
	case *syntax.ParamExp:
		// keep as-is; no env expansion
		return "${" + t.Param.Value + "}"
	case *syntax.DblQuoted:
		var b strings.Builder
		for _, q := range t.Parts {
			b.WriteString(partToText(q))
		}
		return b.String()
	default:
		// Fallback to printer for things like command subs, arithmetic, etc.
		var sb strings.Builder
		syntax.NewPrinter().Print(&sb, n)
		return sb.String()
	}
}
