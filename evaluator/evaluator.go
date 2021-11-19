package evaluator

import (
	"fmt"
	"lookageek.com/ode/ast"
	"lookageek.com/ode/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval function is the entry point to which the parsed AST node is passed
// it walks the tree recursively and evaluated the nodes
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBooleanToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)

		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)

		if isError(left) {
			return left
		}

		right := Eval(node.Right)

		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node)

	case *ast.IfExpression:
		return evalIfExpression(node)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)

		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	}

	return nil
}

// evalIfExpression evaluates an IfExpression node
// it first evaluates the condition in the IfExpression and
// if true evaluates the Consequence BlockStatement or else
// evaluates the Alternative BlockStatement
func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

// evalProgram takes a slice of Statement nodes and evaluates them one
// by one
func evalProgram(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)

		// if a return statement or an error is encountered,
		//stop the evaluation of the program further on
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil {
			rt := result.Type()

			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

// converts the literal boolean to a Boolean Object
func nativeBooleanToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

// evalPrefixExpression evals a prefix expression which consists of ! or - expression
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// evalInfixExpression evals normal math operators +, -, *, /
// also comparision operators of ==, !=, <, >, for integer operands
// and == & != for boolean operands
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	// take care of operands being integers in either math operators or in boolean operators
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	// now operands can only be boolean object hence create boolean object with native comparision
	case operator == "==":
		return nativeBooleanToBooleanObject(left == right)
	case operator == "!=":
		return nativeBooleanToBooleanObject(left != right)
	case left.Type() != right.Type():
		// types in infix expression should match, or else it raises error
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		// in cases which have wrong operands or wrong operator, an error is raised
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalIntegerInfixExpression handles evaluating both math operators and comparision operators
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBooleanToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBooleanToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBooleanToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBooleanToBooleanObject(leftVal != rightVal)
	default:
		// in default case where operator symbol is not any of the above results in an evaluation error
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// newError is a constructor for the error object raised during evaluation
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}
