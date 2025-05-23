package evaluator

import (
	"monkey/ast"
	"monkey/object"
	"fmt"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right){
			return right
		}
		return evalPrefixExpression(node.Operator, right, env)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left){
			return left
		}
		right := Eval(node.Right, env)

		if isError(right){
			return right
		}
		return evalInfixExpression(node.Operator, left, right, env)

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env) 
		if isError(val){
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}
	return nil
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object, env *object.Environment) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right, env)
	case "-":
		return evalMinusPrefixOperatorExpression(right, env)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object, env *object.Environment) object.Object {
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

func evalMinusPrefixOperatorExpression(right object.Object, env *object.Environment) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(
    operator string,
    left, right object.Object, env *object.Environment,
) object.Object {
    switch {
    case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
        return evalIntegerInfixExpression(operator, left, right, env)
		case operator == "==":
				return nativeBoolToBooleanObject(left == right)
		case operator == "!=":
				return nativeBoolToBooleanObject(left != right)
		case left.Type() != right.Type():
				return newError("Type mismatch: %s %s %s", left.Type(), operator, right.Type())
    default:
				return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
    }
}

func evalIntegerInfixExpression( operator string, left, right object.Object, env *object.Environment) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
			return &object.Integer{Value: leftVal + rightVal}
	case "-":
			return &object.Integer{Value: leftVal - rightVal}
	case "*":
			return &object.Integer{Value: leftVal + rightVal}
	case "/":
			return &object.Integer{Value: leftVal / rightVal}
	case "<":
			return nativeBoolToBooleanObject ( leftVal < rightVal )
	case ">":
			return nativeBoolToBooleanObject ( leftVal > rightVal )
	case "==":
			return nativeBoolToBooleanObject ( leftVal == rightVal )
	case "!=":
			return nativeBoolToBooleanObject ( leftVal != rightVal )
	default:
		return newError("unkown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

		if isError(condition){
			return condition
		}
	if isTruthy(condition){
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil{
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy (obj object.Object) bool {
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

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:	
			return result
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env) 

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func newError(format string, a...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier( node *ast.Identifier, env *object.Environment,) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}
	return val
}
