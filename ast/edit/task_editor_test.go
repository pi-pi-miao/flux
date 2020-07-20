package edit_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/influxdata/flux/ast"
	"github.com/influxdata/flux/ast/asttest"
	"github.com/influxdata/flux/ast/edit"
)

var IgnoreInit = []cmp.Option{
	cmpopts.IgnoreFields(ast.VariableAssignment{}, "Init"),
}

func TestGetOption(t *testing.T) {
	testCases := []struct {
		testName string
		optionID string
		file     *ast.File
		want     ast.Node
	}{
		{
			testName: "test getOption",
			optionID: "task",
			file: &ast.File{
				BaseNode: ast.BaseNode{},
				Name:     "foo.flux",
				Body: []ast.Statement{
					&ast.OptionStatement{
						BaseNode: ast.BaseNode{},
						Assignment: &ast.VariableAssignment{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "bar",
							},
							Init: nil,
						},
					},
					&ast.OptionStatement{
						BaseNode: ast.BaseNode{},
						Assignment: &ast.VariableAssignment{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "task",
							},
							Init: nil,
						},
					},
				},
			},
			want: &ast.OptionStatement{
				BaseNode: ast.BaseNode{},
				Assignment: &ast.VariableAssignment{
					BaseNode: ast.BaseNode{},
					ID: &ast.Identifier{
						BaseNode: ast.BaseNode{},
						Name:     "task",
					},
					Init: nil,
				},
			},
		},
		{
			testName: "test getOption with numbers in task name",
			optionID: "numbers222",
			file: &ast.File{
				BaseNode: ast.BaseNode{},
				Name:     "foo.flux",
				Body: []ast.Statement{
					&ast.ExpressionStatement{
						BaseNode:   ast.BaseNode{},
						Expression: nil,
					},
					&ast.OptionStatement{
						BaseNode: ast.BaseNode{},
						Assignment: &ast.VariableAssignment{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "numbers222",
							},
							Init: nil,
						},
					},
				},
			},
			want: &ast.OptionStatement{
				BaseNode: ast.BaseNode{},
				Assignment: &ast.VariableAssignment{
					BaseNode: ast.BaseNode{},
					ID: &ast.Identifier{
						BaseNode: ast.BaseNode{},
						Name:     "numbers222",
					},
					Init: nil,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			got, err := edit.GetOption(tc.file, tc.optionID)
			if err != nil {
				t.Errorf("unexpected error %s", err)
			}

			var ignoreOptions []cmp.Option
			ignoreOptions = append(ignoreOptions, IgnoreInit...)
			ignoreOptions = append(ignoreOptions, asttest.IgnoreBaseNodeOptions...)

			if !cmp.Equal(got, tc.want, ignoreOptions...) {
				t.Errorf("Unexpected value -want/+got:\n%s", cmp.Diff(tc.want, got, ignoreOptions...))
			}
		})
	}
}

func TestSetDeleteOption(t *testing.T) {
	testCases := []struct {
		testName   string
		testType   string
		optionID   string
		got        *ast.File
		want       *ast.File
		opt        ast.Node
	}{
		{
			testName: "test set option",
			testType: "setOption",
			optionID: "foo",
			got: &ast.File{
				BaseNode: ast.BaseNode{},
				Name:     "foo.flux",
				Body: []ast.Statement{
					&ast.OptionStatement{
						BaseNode: ast.BaseNode{},
						Assignment: &ast.VariableAssignment{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "foo",
							},
							Init: nil,
						},
					},
				},
			},
			opt: &ast.OptionStatement{
				BaseNode: ast.BaseNode{},
				Assignment: &ast.VariableAssignment{
					BaseNode: ast.BaseNode{},
					ID: &ast.Identifier{
						BaseNode: ast.BaseNode{},
						Name:     "wantTask",
					},
					Init: nil,
				},
			},
			want: &ast.File{
				BaseNode: ast.BaseNode{},
				Name:     "foo.flux",
				Body: []ast.Statement{
					&ast.OptionStatement{
						BaseNode: ast.BaseNode{},
						Assignment: &ast.VariableAssignment{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "wantTask",
							},
							Init: nil,
						},
					},
				},
			},
		},
		{
			testName: "test setOption with numbers in task name",
			testType: "setOption",
			optionID: "bar",
			got: &ast.File{
				BaseNode: ast.BaseNode{},
				Name:     "foo.flux",
				Body: []ast.Statement{
					&ast.OptionStatement{
						BaseNode: ast.BaseNode{},
						Assignment: &ast.VariableAssignment{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "bar",
							},
							Init: nil,
						},
					},
				},
			},
			opt: &ast.OptionStatement{
				BaseNode: ast.BaseNode{},
				Assignment: &ast.VariableAssignment{
					BaseNode: ast.BaseNode{},
					ID: &ast.Identifier{
						BaseNode: ast.BaseNode{},
						Name:     "numbers",
					},
					Init: nil,
				},
			},
			want: &ast.File{
				BaseNode: ast.BaseNode{},
				Name:     "foo.flux",
				Body: []ast.Statement{
					&ast.OptionStatement{
						BaseNode: ast.BaseNode{},
						Assignment: &ast.VariableAssignment{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "numbers",
							},
							Init: nil,
						},
					},
				},
			},
		},
		{
			testName: "test deleteOption",
			testType: "deleteOption",
			optionID: "numbers",
			got: &ast.File{
				BaseNode: ast.BaseNode{},
				Name:     "foo.flux",
				Body: []ast.Statement{
					&ast.OptionStatement{
						BaseNode: ast.BaseNode{},
						Assignment: &ast.VariableAssignment{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "bar",
							},
							Init: nil,
						},
					},
					&ast.ExpressionStatement{
						BaseNode:   ast.BaseNode{},
						Expression: nil,
					},
					&ast.OptionStatement{
						BaseNode: ast.BaseNode{},
						Assignment: &ast.VariableAssignment{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "numbers",
							},
							Init: nil,
						},
					},
				},
			},
			want: &ast.File{
				BaseNode: ast.BaseNode{},
				Name:     "foo.flux",
				Body: []ast.Statement{
					&ast.OptionStatement{
						BaseNode: ast.BaseNode{},
						Assignment: &ast.VariableAssignment{
							BaseNode: ast.BaseNode{},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{},
								Name:     "bar",
							},
							Init: nil,
						},
					},
					&ast.ExpressionStatement{
						BaseNode:   ast.BaseNode{},
						Expression: nil,
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			switch tc.testType {
			case "setOption":
				edit.SetOption(tc.got, tc.optionID, tc.opt)
			case "deleteOption":
				edit.DeleteOption(tc.got, tc.optionID)
			default:
				t.Fatal("Test type must be set to 'setOption' or 'deleteOption'.")
			}

			var ignoreOptions []cmp.Option
			ignoreOptions = append(ignoreOptions, IgnoreInit...)
			ignoreOptions = append(ignoreOptions, asttest.IgnoreBaseNodeOptions...)

			if !cmp.Equal(tc.got, tc.want, ignoreOptions...) {
				t.Errorf("Unexpected value -want/+got:\n%s", cmp.Diff(tc.want, tc.got, ignoreOptions...))
			}
		})
	}
}

func TestGetProperty(t *testing.T) {
	testCases := []struct {
		testName string
		key      string
		want     *ast.Property
		obj      *ast.ObjectExpression
	}{
		{
			testName: "test getProperty with boolean",
			key:      "b",
			want: &ast.Property{
				BaseNode: ast.BaseNode{},
				Key:      &ast.StringLiteral{Value: "b"},
				Value:    &ast.BooleanLiteral{Value: true},
			},
			obj: &ast.ObjectExpression{
				BaseNode: ast.BaseNode{},
				With:     nil,
				Properties: []*ast.Property{
					{
						Key:   &ast.StringLiteral{Value: "a"},
						Value: &ast.StringLiteral{Value: "hello"},
					},
					{
						Key:   &ast.StringLiteral{Value: "b"},
						Value: &ast.BooleanLiteral{Value: true},
					},
				},
			},
		},
		{
			testName: "test getProperty with integer",
			key:      "foo",
			want: &ast.Property{
				BaseNode: ast.BaseNode{},
				Key:      &ast.StringLiteral{Value: "foo"},
				Value:    &ast.StringLiteral{Value: "hello"},
			},
			obj: &ast.ObjectExpression{
				BaseNode: ast.BaseNode{},
				With:     nil,
				Properties: []*ast.Property{
					{
						Key:   &ast.StringLiteral{Value: "foo"},
						Value: &ast.StringLiteral{Value: "hello"},
					},
					{
						Key:   &ast.StringLiteral{Value: "bar"},
						Value: &ast.IntegerLiteral{Value: 100},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			got, err := edit.GetProperty(tc.obj, tc.key)
			if err != nil {
				t.Errorf("unexpected error %s", err)
			}

			var ignoreOptions []cmp.Option
			ignoreOptions = append(ignoreOptions, IgnoreInit...)
			ignoreOptions = append(ignoreOptions, asttest.IgnoreBaseNodeOptions...)

			if !cmp.Equal(got, tc.want, ignoreOptions...) {
				t.Errorf("Unexpected value -want/+got:\n%s", cmp.Diff(tc.want, got, ignoreOptions...))
			}
		})
	}
}

func TestSetDeleteProperty(t *testing.T) {
	testCases := []struct {
		testName string
		key      string
		testType string
		want     ast.Node
		obj      *ast.ObjectExpression
		value    ast.Expression
	}{
		{
			testName: "test setProperty with float",
			testType: "setProperty",
			key:      "foo",
			obj: &ast.ObjectExpression{
				BaseNode: ast.BaseNode{},
				With:     nil,
				Properties: []*ast.Property{
					{
						Key:   &ast.StringLiteral{Value: "foo"},
						Value: &ast.StringLiteral{Value: "hello"},
					},
					{
						Key:   &ast.StringLiteral{Value: "bar"},
						Value: &ast.IntegerLiteral{Value: 100},
					},
				},
			},
			value: &ast.FloatLiteral{
				Value: 1.23,
			},
			want: &ast.ObjectExpression{
				BaseNode: ast.BaseNode{},
				With:     nil,
				Properties: []*ast.Property{
					{
						Key:   &ast.StringLiteral{Value: "foo"},
						Value: &ast.FloatLiteral{Value: 1.23},
					},
					{
						Key:   &ast.StringLiteral{Value: "bar"},
						Value: &ast.IntegerLiteral{Value: 100},
					},
				},
			},
		},
		{
			testName: "test setProperty with date time",
			testType: "setProperty",
			key:      "otherTest",
			obj: &ast.ObjectExpression{
				BaseNode: ast.BaseNode{},
				With:     nil,
				Properties: []*ast.Property{
					{
						Key:   &ast.StringLiteral{Value: "test"},
						Value: &ast.StringLiteral{Value: "hello"},
					},
					{
						Key:   &ast.StringLiteral{Value: "otherTest"},
						Value: &ast.IntegerLiteral{Value: 100},
					},
				},
			},
			value: &ast.DateTimeLiteral{
				Value: time.Date(2017, 8, 8, 8, 8, 8, 8, time.UTC),
			},
			want: &ast.ObjectExpression{
				BaseNode: ast.BaseNode{},
				With:     nil,
				Properties: []*ast.Property{
					{
						Key:   &ast.StringLiteral{Value: "test"},
						Value: &ast.StringLiteral{Value: "hello"},
					},
					{
						Key: &ast.StringLiteral{Value: "otherTest"},
						Value: &ast.DateTimeLiteral{
							Value: time.Date(2017, 8, 8, 8, 8, 8, 8, time.UTC),
						},
					},
				},
			},
		},
		{
			testName: "test deleteProperty with duration",
			testType: "deleteProperty",
			key:      "test",
			obj: &ast.ObjectExpression{
				BaseNode: ast.BaseNode{},
				With:     nil,
				Properties: []*ast.Property{
					{
						Key: &ast.StringLiteral{Value: "test"},
						Value: &ast.DurationLiteral{
							Values: []ast.Duration{{
								Magnitude: 1,
								Unit:      "s",
							}},
						},
					},
					{
						Key:   &ast.StringLiteral{Value: "otherTest"},
						Value: &ast.IntegerLiteral{Value: 100},
					},
				},
			},
			want: &ast.ObjectExpression{
				BaseNode: ast.BaseNode{},
				With:     nil,
				Properties: []*ast.Property{
					{
						Key:   &ast.StringLiteral{Value: "otherTest"},
						Value: &ast.IntegerLiteral{Value: 100},
					},
				},
			},
		},
		{
			testName: "test deleteProperty with float",
			testType: "deleteProperty",
			key:      "bar",
			obj: &ast.ObjectExpression{
				BaseNode: ast.BaseNode{},
				With:     nil,
				Properties: []*ast.Property{
					{
						Key: &ast.StringLiteral{Value: "foo"},
						Value: &ast.DurationLiteral{
							Values: []ast.Duration{{
								Magnitude: 1,
								Unit:      "s",
							}},
						},
					},
					{
						Key:   &ast.StringLiteral{Value: "bar"},
						Value: &ast.IntegerLiteral{Value: 100},
					},
				},
			},
			want: &ast.ObjectExpression{
				BaseNode: ast.BaseNode{},
				With:     nil,
				Properties: []*ast.Property{
					{
						Key: &ast.StringLiteral{Value: "foo"},
						Value: &ast.DurationLiteral{
							Values: []ast.Duration{{
								Magnitude: 1,
								Unit:      "s",
							}},
						},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			var err error
			switch tc.testType {
			case "setProperty":
				edit.SetProperty(tc.obj, tc.key, tc.value)
			case "deleteProperty":
				err = edit.DeleteProperty(tc.obj, tc.key)
			default:
				t.Fatal("Test type must be set to 'setProperty' or 'deleteProperty'.")
			}
			if err != nil {
				t.Errorf("unexpected error %s", err)
			}

			var ignoreOptions []cmp.Option
			ignoreOptions = append(ignoreOptions, IgnoreInit...)
			ignoreOptions = append(ignoreOptions, asttest.IgnoreBaseNodeOptions...)

			if !cmp.Equal(tc.obj, tc.want, ignoreOptions...) {
				t.Errorf("Unexpected value -want/+got:\n%s", cmp.Diff(tc.want, tc.obj, ignoreOptions...))
			}
		})
	}
}
