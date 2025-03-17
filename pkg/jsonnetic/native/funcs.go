//nolint:errcheck
package native

import (
	"regexp"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/neticdk/go-common/pkg/file"
	"github.com/neticdk/go-jsonnetic/internal/utils"
	"github.com/prometheus/prometheus/promql/parser"
)

const (
	FuncFileExists        = "fileExists"
	FuncRegexMatch        = "regexMatch"
	FuncRegexSubst        = "regexSubst"
	FuncPromqlAggregateBy = "promqlAggregateBy"
)

func Funcs() []*jsonnet.NativeFunction {
	return []*jsonnet.NativeFunction{
		fileExists(),
		regexMatch(),
		regexSubst(),
		promqlAggregateBy(),
	}
}

// fileExists returns whether the given file exists.
func fileExists() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   FuncFileExists,
		Params: ast.Identifiers{"path"},
		Func: func(args []any) (any, error) {
			return file.Exists(args[0].(string))
		},
	}
}

// regexMatch returns whether the given string is matched by the given re2 regular expression.
func regexMatch() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   FuncRegexMatch,
		Params: ast.Identifiers{"regex", "string"},
		Func: func(args []any) (any, error) {
			return regexp.MatchString(args[0].(string), args[1].(string))
		},
	}
}

// regexSubst replaces all matches of the re2 regular expression with another string.
func regexSubst() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   FuncRegexSubst,
		Params: ast.Identifiers{"regex", "src", "repl"},
		Func: func(args []any) (any, error) {
			regex, src, repl := args[0].(string), args[1].(string), args[2].(string)

			r, err := regexp.Compile(regex)
			if err != nil {
				return "", err
			}
			return r.ReplaceAllString(src, repl), nil
		},
	}
}

// promqlAggregateBy modifies a PromQL expression to include a given label.
func promqlAggregateBy() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   FuncPromqlAggregateBy,
		Params: ast.Identifiers{"expr", "label"},
		Func: func(args []any) (any, error) {
			expr, err := parser.ParseExpr(args[0].(string))
			if err != nil {
				return "", err
			}
			f := utils.ExprNodeInspectorFunc(args[1].(string))
			parser.Inspect(expr, f)

			return expr.String(), nil
		},
	}
}
