package visitor

import (
	"fmt"
	"strconv"

	ast "github.com/codecrafters-io/interpreter-starter-go/cmd/generator/out"
	"github.com/codecrafters-io/interpreter-starter-go/internal/token"
)

type Interpreter struct {
	errors []error
}

func (i *Interpreter) VisitExprBinary(expr *ast.Binary) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	var arithmeticOp func(a, b float64) float64
	var comparisonOp func(a, b float64) bool
	var equalityOp func(a, b any) bool
	var stringConcat func(a, b string) string
	switch expr.Operator.Type {
	case token.MINUS:
		arithmeticOp = func(a, b float64) float64 { return a - b }
	case token.SLASH:
		arithmeticOp = func(a, b float64) float64 { return a / b }
	case token.STAR:
		arithmeticOp = func(a, b float64) float64 { return a * b }
	case token.PLUS:
		switch left.(type) {
		case string:
			stringConcat = func(a, b string) string { return a + b }
		case float64:
			arithmeticOp = func(a, b float64) float64 { return a + b }
		default:
			return nil, errorFunc(left, "unsupported operand type", expr.Operator.Line)
		}
	case token.GREATER:
		comparisonOp = func(a, b float64) bool { return a > b }
	case token.GREATER_EQUAL:
		comparisonOp = func(a, b float64) bool { return a >= b }
	case token.LESS:
		comparisonOp = func(a, b float64) bool { return a < b }
	case token.LESS_EQUAL:
		comparisonOp = func(a, b float64) bool { return a <= b }
	case token.EQUAL:
		equalityOp = i.isEqual
	case token.BANG_EQUAL:
		equalityOp = func(a, b any) bool {
			return !i.isEqual(a, b)
		}
	}

	if arithmeticOp != nil {
		if leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		} else {
			return arithmeticOp(leftVal, rightVal), nil
		}
	}
	if comparisonOp != nil {
		if leftVal, rightVal, err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		} else {
			return comparisonOp(leftVal, rightVal), nil
		}
	}
	if equalityOp != nil {
		return equalityOp(left, right), nil
	}
	if stringConcat != nil {
		if leftVal, rightVal, err := i.checkStringOperands(expr.Operator, left, right); err != nil {
			return nil, err
		} else {
			return stringConcat(leftVal, rightVal), nil
		}
	}
	panic("unreachable")
}

func (i *Interpreter) VisitExprGrouping(expr *ast.Grouping) (any, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitExprLiteral(expr *ast.Literal) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitExprUnary(expr *ast.Unary) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Type {
	case token.BANG:
		return !i.isTruthy(right), nil
	case token.MINUS:
		if val, err := i.checkNumberOperand(expr.Operator, right); err != nil {
			return nil, err
		} else {
			return -val, nil
		}
	}
	panic("unreachable")
}

var (
	_ ast.Visitor[any] = &Interpreter{}
)

func (i *Interpreter) evaluate(expr ast.Expr) (any, error) {
	return expr.Accept(i)
}
func (i *Interpreter) Interpret(expr ast.Expr) (any, error) {
	return i.evaluate(expr)
}

func (i *Interpreter) isTruthy(right any) bool {
	if right == nil {
		return false
	}
	if val, ok := right.(bool); ok {
		return val
	}
	return true
}
func (i *Interpreter) isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a == b
}

func (i *Interpreter) checkNumberOperand(operator token.Token, operand any) (float64, error) {
	if val, ok := operand.(float64); !ok {
		return 0, errorFunc(operand, fmt.Sprintf("%s expects operand to be number",
			operator.String()), operator.Line)
	} else {
		return val, nil
	}
}
func (i *Interpreter) checkNumberOperands(operator token.Token, left, right any) (leftVal, rightVal float64, err error) {
	if leftVal, err = i.checkNumberOperand(operator, left); err != nil {
		return 0, 0, err
	}
	if rightVal, err = i.checkNumberOperand(operator, right); err != nil {
		return 0, 0, err
	}
	return leftVal, rightVal, nil
}
func (i *Interpreter) checkStringOperand(operator token.Token, operand any) (string, error) {
	if val, ok := operand.(string); !ok {
		return "", errorFunc(operand, fmt.Sprintf("%s expects operand to be string",
			operator.String()), operator.Line)
	} else {
		return val, nil
	}
}
func (i *Interpreter) checkStringOperands(operator token.Token, left any, right any) (leftVal, rightVal string, err error) {
	if leftVal, err = i.checkStringOperand(operator, left); err != nil {
		return "", "", err
	}
	if rightVal, err = i.checkStringOperand(operator, right); err != nil {
		return "", "", err
	}
	return leftVal, rightVal, nil
}
func errorFunc(actual interface{}, expectation string, line int) error {
	return fmt.Errorf("%s, but got %s",
		expectation, actual.(fmt.Stringer).String())
}

func (i *Interpreter) Stringer(obj any) string {
	return Stringer(obj)
}

func Stringer(obj any) string {
	if obj == nil {
		return "nil"
	}
	switch obj.(type) {
	case string:
		return fmt.Sprintf("%s", obj)
	case float64:
		val := strconv.FormatFloat(obj.(float64), 'f', -1, 64)
		return val
	case bool:
		return fmt.Sprintf("%t", obj)
	default:
		return fmt.Sprintf("%v", obj)
	}
}
