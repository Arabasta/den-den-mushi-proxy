package main

import (
	"fmt"
	"os"
	"strings"

	"mvdan.cc/sh/v3/syntax"
)

func main() {
	cmd := "echo foo${bar}'baz'|ls|su|ls -lrt" // change to `echo foo && bar` to test BinaryCmd

	r := strings.NewReader(cmd)
	f, err := syntax.NewParser().Parse(r, "")
	if err != nil {
		panic(err)
	}

	printer := syntax.NewPrinter()

	for _, stmt := range f.Stmts {
		switch x := stmt.Cmd.(type) {
		case *syntax.CallExpr:
			fmt.Println("== CallExpr ==")
			handleCall(x, printer)
		case *syntax.BinaryCmd:
			fmt.Println("== BinaryCmd ==")
			fmt.Println("Op:", x.Op)
			// handle left and right separately
			fmt.Println("-- Left --")
			printCmd(x.X.Cmd, printer)
			fmt.Println("-- Right --")
			printCmd(x.Y.Cmd, printer)
		default:
			fmt.Printf("Unhandled type: %T\n", x)
		}
	}
}

func handleCall(call *syntax.CallExpr, printer *syntax.Printer) {
	var cmdArr []string
	for i, word := range call.Args {
		fmt.Printf("Word number %d:\n", i)
		var sb strings.Builder
		for _, part := range word.Parts {
			fmt.Printf("%-20T - ", part)
			printer.Print(os.Stdout, part)
			fmt.Println()
			sb.WriteString(getPartText(part))
		}
		cmdArr = append(cmdArr, sb.String())
	}
	fmt.Println("Final cmd array:", cmdArr)
}

func printCmd(cmd syntax.Command, printer *syntax.Printer) {
	switch x := cmd.(type) {
	case *syntax.CallExpr:
		handleCall(x, printer)
	default:
		fmt.Printf("Nested type not handled: %T\n", x)
	}
}

// same as before
func getPartText(node syntax.Node) string {
	switch x := node.(type) {
	case *syntax.Lit:
		return x.Value
	case *syntax.SglQuoted:
		return x.Value
	case *syntax.DblQuoted:
		var sb strings.Builder
		for _, p := range x.Parts {
			sb.WriteString(getPartText(p))
		}
		return sb.String()
	case *syntax.ParamExp:
		return "${" + x.Param.Value + "}"
	default:
		return ""
	}
}
