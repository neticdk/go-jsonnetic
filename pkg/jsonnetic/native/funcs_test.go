package native

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockCallNative func() (interface{}, error, error)

func callNative(name string, data []interface{}) func() (interface{}, error, error) {
	return func() (interface{}, error, error) {
		for _, fun := range Funcs() {
			if fun.Name == name {
				// Call the function
				ret, err := fun.Func(data)
				return ret, err, nil
			}
		}

		return nil, nil, fmt.Errorf("could not find native function %s", name)
	}
}

func TestFileExists(t *testing.T) {
	tests := []struct {
		name        string
		call        mockCallNative
		expected    interface{}
		expectError bool
	}{
		{
			name:        "file does exists",
			call:        callNative(FuncFileExists, []interface{}{"fixtures/goodFile.txt"}),
			expected:    true,
			expectError: false,
		},
		{
			name:        "file does not exists",
			call:        callNative(FuncFileExists, []interface{}{"fixtures/doesNotExist.txt"}),
			expected:    false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret, err, callerr := tt.call()

			assert.Empty(t, callerr)
			assert.Equal(t, tt.expected, ret)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegexMatch(t *testing.T) {
	tests := []struct {
		name        string
		call        mockCallNative
		expected    interface{}
		expectError bool
	}{
		{
			name:        "valid regex",
			call:        callNative(FuncRegexMatch, []interface{}{"", "a"}),
			expected:    true,
			expectError: false,
		},
		{
			name:        "valid regex, no match",
			call:        callNative(FuncRegexMatch, []interface{}{"a", "b"}),
			expected:    false,
			expectError: false,
		},
		{
			name:        "invalid regex",
			call:        callNative(FuncRegexMatch, []interface{}{"[0-", "b"}),
			expected:    false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret, err, callerr := tt.call()

			assert.Empty(t, callerr)
			assert.Equal(t, tt.expected, ret)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegexSubst(t *testing.T) {
	tests := []struct {
		name        string
		call        mockCallNative
		expected    interface{}
		expectError bool
	}{
		{
			name:        "valid regex, but no changes",
			call:        callNative(FuncRegexSubst, []interface{}{"a", "b", "c"}),
			expected:    "b",
			expectError: false,
		},
		{
			name:        "valid regex, with changes",
			call:        callNative(FuncRegexSubst, []interface{}{"p[^m]*", "pm", "poe"}),
			expected:    "poem",
			expectError: false,
		},
		{
			name:        "invalid regex",
			call:        callNative(FuncRegexSubst, []interface{}{"p[^m*", "pm", "poe"}),
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret, err, callerr := tt.call()

			assert.Empty(t, callerr)
			assert.Equal(t, tt.expected, ret)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPromqlAggregateBy(t *testing.T) {
	tests := []struct {
		name        string
		call        mockCallNative
		expected    interface{}
		expectError bool
	}{
		{
			name:        "expression with valid aggregation labels",
			call:        callNative(FuncPromqlAggregateBy, []interface{}{"sum by (cluster) (rate(foo[5m]))", "cluster"}),
			expected:    "sum by (cluster) (rate(foo[5m]))",
			expectError: false,
		},
		{
			name:        "expression with missing aggregation labels",
			call:        callNative(FuncPromqlAggregateBy, []interface{}{"sum(rate(foo[5m]))", "cluster"}),
			expected:    "sum by (cluster) (rate(foo[5m]))",
			expectError: false,
		},
		{
			name:        "binary expression with missing aggregation labels",
			call:        callNative(FuncPromqlAggregateBy, []interface{}{"sum(max(foo) * on (bar) group_left (baz) topk(1, max(qux))) > 0", "cluster"}),
			expected:    "sum by (cluster) (max by (cluster) (foo) * on (bar, cluster) group_left (baz) topk by (cluster) (1, max by (cluster) (qux))) > 0",
			expectError: false,
		},
		{
			name:        "expression using without aggregation, and label is not present",
			call:        callNative(FuncPromqlAggregateBy, []interface{}{"sum without (bar) (rate(foo[5m]))", "cluster"}),
			expected:    "sum without (bar) (rate(foo[5m]))",
			expectError: false,
		},
		{
			name:        "expression using without aggregation, and label is present",
			call:        callNative(FuncPromqlAggregateBy, []interface{}{"sum without (bar, cluster) (rate(foo[5m]))", "cluster"}),
			expected:    "sum without (bar) (rate(foo[5m]))",
			expectError: false,
		},
		{
			name:        "valid expression without aggregation",
			call:        callNative(FuncPromqlAggregateBy, []interface{}{"foo", "cluster"}),
			expected:    "foo",
			expectError: false,
		},
		{
			name:        "invalid expression",
			call:        callNative(FuncPromqlAggregateBy, []interface{}{"sum{foo}", "cluster"}),
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ret, err, callerr := tt.call()

			assert.Empty(t, callerr)
			assert.Equal(t, tt.expected, ret)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
