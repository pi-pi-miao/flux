package edit

import (
	"fmt"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/ast"
	"github.com/influxdata/flux/codes"
)

// GetOption finds and returns the AST node for the option value.
func GetOption(file *ast.File, name string) (ast.Node, error) {
	for _, st := range file.Body {
		val, ok := st.(*ast.OptionStatement)
		// wrong statement type; skip this statement
		if !ok {
			continue
		}
		assign := val.Assignment
		va, ok := assign.(*ast.VariableAssignment)
		if va.ID.Name == name {
			if ok {
				return val, nil
			}
		}
	}

	return nil, &flux.Error{
		Code: codes.Internal,
		Msg:  "Option not found",
	}
}

// SetOption replaces an existing option definition with the provided option node or adds
// the option if it doesn't exist. The pkg AST is mutated in place.
func SetOption(file *ast.File, id string, opt ast.Node) {
	// check for the correct file
	for i, st := range file.Body {
		val, ok := st.(*ast.OptionStatement)
		// wrong statement type; skip this statement
		if !ok {
			continue
		}
		assign := val.Assignment
		va, ok := assign.(*ast.VariableAssignment)
		if ok {
			if va.ID.Name == id {
				if optSt, ok := opt.(*ast.OptionStatement); ok {
					file.Body[i] = optSt
					return
				}
			}
		}

	}

	file.Body = append(file.Body, &ast.OptionStatement{
		BaseNode:   ast.BaseNode{},
		Assignment: nil,
	})
}

// DeleteOption removes an option if it exists. The pkg AST is mutated in place.
func DeleteOption(file *ast.File, name string) error {
	for i, st := range file.Body {
		val, ok := st.(*ast.OptionStatement)
		if !ok {
			continue
		}
		assign := val.Assignment
		va, ok := assign.(*ast.VariableAssignment)
		if va.ID.Name == name {
			if ok {
				file.Body = append(file.Body[:i], file.Body[i+1:]...)
				return nil
			}
		}
	}

	return &flux.Error{
		Code: codes.Internal,
		Msg:  fmt.Sprintf("Could not delete %v, Option not found.", name),
	}
}

// GetProperty finds and returns the AST node for the property value.
func GetProperty(obj *ast.ObjectExpression, key string) (*ast.Property, error) {
	for _, prop := range obj.Properties {
		if key == prop.Key.Key() {
			return prop, nil
		}
	}
	return nil, &flux.Error{
		Code: codes.Internal,
		Msg:  "Property not found",
	}
}

// SetOption replaces an existing property definition with the provided object expression or adds
// the property if it doesn't exist. The object expression AST is mutated in place.
func SetProperty(obj *ast.ObjectExpression, key string, value ast.Expression) {
	for _, prop := range obj.Properties {
		if key == prop.Key.Key() {
			prop.Value = value
			return
		}
	}
	obj.Properties = append(obj.Properties, &ast.Property{
		BaseNode: obj.BaseNode,
		Key:      &ast.Identifier{Name: key},
		Value:    value,
	})
}

// DeleteProperty removes a property from the object expression if it exists. The object AST is mutated in place.
func DeleteProperty(obj *ast.ObjectExpression, key string) error {
	for i, prop := range obj.Properties {
		if key == prop.Key.Key() {
			obj.Properties = append(obj.Properties[:i], obj.Properties[i+1:]...)
			return nil
		}
	}
	return &flux.Error{
		Code: codes.Internal,
		Msg:  fmt.Sprintf("Could not delete %v, Property not found.", key),
	}
}
