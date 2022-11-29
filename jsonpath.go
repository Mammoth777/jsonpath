package jsonpath

import (
	"errors"
	"strconv"

	"github.com/Mammoth777/jsonpath/core"
)

func Compile(path string) []*core.Step {
	steps := make([]*core.Step, 0)
	for i := 0; i < len(path); i++ {
		c := path[i]
		step := core.NewStep()
		if c == '$' {
			step.Action(core.ROOT_ACTION).Key(string(c))
		} else if c == '.' {
			i = step.StickDot(path, i)
		} else if c == '[' {
			i = step.StickBracket(path, i)
		}
		steps = append(steps, step)
	}
	return steps
}

func getValue(data any, steps []*core.Step) (any, error) {
	if len(steps) == 0 {
		return data, errors.New("steps is empty")
	}
	var (
		temp any
		ok   bool
	)
	for _, step := range steps {
		switch step.GetAction() {
		case core.ROOT_ACTION:
			temp = data
		case core.KEY_ACTION:
			temp, ok = temp.(map[string]any)[step.GetKey()]
			if !ok {
				return temp, errors.New("key not found")
			}
		case core.IDX_ACTION:
			list, ok := temp.([]any)
			if !ok {
				return temp, errors.New("not a list")
			}
			index, err := strconv.Atoi(step.GetKey())
			if err != nil {
				return temp, errors.New("index is not a number")
			}
			temp = list[index]
		default:
			return temp, errors.New("unknown action")
		}
	}
	return temp, nil
}

func setValue(data any, steps []*core.Step, value any) (any, error) {
	if len(steps) == 0 {
		return data, errors.New("steps is empty")
	}
	var (
		temp any
		ok   bool
	)
	for i, step := range steps {
		switch step.GetAction() {
		case core.ROOT_ACTION:
			temp = data
		case core.KEY_ACTION:
			if i == len(steps)-1 {
				temp.(map[string]any)[step.GetKey()] = value
			} else {
				temp, ok = temp.(map[string]any)[step.GetKey()]
				if !ok {
					return data, errors.New("key not found")
				}
			}
		case core.IDX_ACTION:
			if i == len(steps)-1 {
				list, ok := temp.([]any)
				if !ok {
					return data, errors.New("not a list")
				}
				index, err := strconv.Atoi(step.GetKey())
				if err != nil {
					return data, errors.New("index is not a number")
				}
				list[index] = value
			} else {
				list, ok := temp.([]any)
				if !ok {
					return data, errors.New("not a list")
				}
				index, err := strconv.Atoi(step.GetKey())
				if err != nil {
					return data, errors.New("index is not a number")
				}
				temp = list[index]
			}
		default:
			return data, errors.New("unknown action")
		}
	}
	return data, nil
}

func Read(data any, path string) (any, error) {
	steps := Compile(path)
	return getValue(data, steps)
}

func Write(data any, path string, value any) (any, error) {
	steps := Compile(path)
	return setValue(data, steps, value)
}
