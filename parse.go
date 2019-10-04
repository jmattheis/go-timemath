package timemath

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NowKey the string that represents now.
const NowKey = "now"

// Parse parses time.
func Parse(now time.Time, value string, startOf bool, weekday time.Weekday) (time.Time, error) {
	if !strings.HasPrefix(value, NowKey) {
		if strings.Contains(value, NowKey) {
			return time.Time{}, errors.New("'now' must be at the start")
		}
		return time.Time{}, errors.New("value must be a valid rfc3339 date or start with 'now'")
	}

	ctx := parseContext{
		value:     value,
		index:     len(NowKey),
		parseType: typeOperation,
		operation: OperationNone,
		result:    now,
	}

	for ctx.index < len(value) {
		var err error
		ctx, err = process(ctx, startOf, weekday)
		if err != nil {
			return time.Time{}, err
		}
	}
	if ctx.parseType == typeNumber {
		return ctx.result, errors.New("expected number at the end but got nothing")
	}
	if ctx.parseType == typeUnit {
		return ctx.result, errors.New("expected unit at the end but got nothing")
	}
	return ctx.result, nil
}

func process(ctx parseContext, startOf bool, weekday time.Weekday) (parseContext, error) {
	switch ctx.parseType {
	case typeOperation:
		return processOperation(ctx)
	case typeNumber:
		return processNumber(ctx)
	case typeUnit:
		return processUnit(ctx, startOf, weekday)
	default:
		panic("unknown parse type")
	}
}

func processOperation(ctx parseContext) (parseContext, error) {
	op := operation(ctx.value[ctx.index])
	switch op {
	case OperationAdd, OperationDivide, OperationSubtract:
		ctx.index = ctx.index + 1
		ctx.operation = op
		if op == OperationDivide {
			ctx.parseType = typeUnit
		} else {
			ctx.parseType = typeNumber
		}
		return ctx, nil
	default:
		return ctx, fmt.Errorf("expected operation / + - at index %d but was %q", ctx.index, op)
	}
}

func processUnit(ctx parseContext, startOf bool, weekday time.Weekday) (parseContext, error) {
	u := Unit(ctx.value[ctx.index])
	switch u {
	case Second, Minute, Hour, Day, Week, Month, Year:
		newResult := ctx.result
		switch ctx.operation {
		case OperationSubtract:
			newResult = u.Subtract(ctx.result, ctx.number)
		case OperationAdd:
			newResult = u.Add(ctx.result, ctx.number)
		case OperationDivide:
			if startOf {
				newResult = u.StartOf(ctx.result, weekday)
			} else {
				newResult = u.EndOf(ctx.result, weekday)
			}
		default:
			panic("unknown operation")
		}
		return parseContext{
			index:     ctx.index + 1,
			result:    newResult,
			number:    -1,
			operation: OperationNone,
			parseType: typeOperation,
			value:     ctx.value,
		}, nil
	default:
		return ctx, fmt.Errorf("expected unit y M w d h m s at index %d but was %q", ctx.index, u)
	}

}

func processNumber(ctx parseContext) (parseContext, error) {
	first := ctx.value[ctx.index : ctx.index+1]
	_, err := strconv.Atoi(first)
	if err != nil {
		return ctx, fmt.Errorf("expected number at index %d but was %s", ctx.index, first)
	}
	numberIndex := 1
	for len(ctx.value) > ctx.index+numberIndex+1 && validNumber(ctx.value[ctx.index:ctx.index+numberIndex+1]) {
		numberIndex++
	}

	number, err := strconv.Atoi(ctx.value[ctx.index : ctx.index+numberIndex])

	ctx.number = number
	ctx.index = ctx.index + numberIndex
	ctx.parseType = typeUnit
	return ctx, err
}

func validNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

type parseContext struct {
	value     string
	index     int
	parseType parseType
	operation operation
	number    int
	result    time.Time
}

type parseType string
type operation rune

var (
	typeUnit      parseType = "unit"
	typeOperation parseType = "operation"
	typeNumber    parseType = "number"

	// OperationDivide the divide operator.
	OperationDivide operation = '/'
	// OperationAdd the add operator.
	OperationAdd operation = '+'
	// OperationSubtract the subtract operator.
	OperationSubtract operation = '-'
	// OperationNone the none operator.
	OperationNone operation = '_'
)
