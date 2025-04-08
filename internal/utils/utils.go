package utils

import (
	"github.com/neticdk/go-stdlib/xslices"
	"github.com/prometheus/prometheus/promql/parser"
)

// ExprNodeInspectorFunc returns a PromQL inspector.
func ExprNodeInspectorFunc(label string) func(node parser.Node, path []parser.Node) error {
	return func(node parser.Node, _ []parser.Node) error {
		switch n := node.(type) {
		case *parser.AggregateExpr:
			return prepareAggregationExpr(n, label)
		case *parser.BinaryExpr:
			return prepareBinaryExpr(n, label)
		default:
			return nil
		}
	}
}

// prepareAggregationExpr modifies the aggregation expression to include the given label.
func prepareAggregationExpr(e *parser.AggregateExpr, label string) error {
	if e.Without {
		// If the label is present in the omission, we should remove it.
		e.Grouping = xslices.Filter[string](e.Grouping, func(s string) bool { return s != label })
		return nil
	}

	for _, lbl := range e.Grouping {
		if lbl == label {
			return nil
		}
	}

	e.Grouping = append(e.Grouping, label)
	return nil
}

// prepareBinaryExpr modifies the binary expression to include the given label.
func prepareBinaryExpr(e *parser.BinaryExpr, label string) error {
	if e.VectorMatching == nil || !e.VectorMatching.On {
		return nil
	}

	// Skip if the aggregation label is already present in the group_left/right() or on() clause
	for _, lbl := range append(e.VectorMatching.MatchingLabels, e.VectorMatching.Include...) {
		if lbl == label {
			return nil
		}
	}

	e.VectorMatching.MatchingLabels = append(e.VectorMatching.MatchingLabels, label)
	return nil
}
