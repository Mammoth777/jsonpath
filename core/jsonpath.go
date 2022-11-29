package core

import (
	"errors"
	"strconv"
)

type actionType string

var (
	ROOT_ACTION = actionType("root")
	KEY_ACTION  = actionType("key")
	IDX_ACTION  = actionType("index")
)

func Compile(path string) []*Step {
	steps := make([]*Step, 0)
	for i := 0; i < len(path); i++ {
		c := path[i]
		step := NewStep()
		if c == '$' {
			step.Action(ROOT_ACTION).Key(string(c))
		} else if c == '.' {
			i = step.stickDot(path, i)
		} else if c == '[' {
			i = step.stickBracket(path, i)
		}
		steps = append(steps, step)
	}
	return steps
}

func getValue(data any, steps []*Step) (any, error) {
	if len(steps) == 0 {
		return data, errors.New("steps is empty")
	}
	var (
		temp any
		ok   bool
	)
	for _, step := range steps {
		switch step.GetAction() {
		case ROOT_ACTION:
			temp = data
		case KEY_ACTION:
			temp, ok = temp.(map[string]any)[step.GetKey()]
			if !ok {
				return temp, errors.New("key not found")
			}
		case IDX_ACTION:
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

func setValue(data any, steps []*Step, value any) (any, error) {
	if len(steps) == 0 {
		return data, errors.New("steps is empty")
	}
	var (
		temp any
		ok   bool
	)
	for i, step := range steps {
		switch step.GetAction() {
		case ROOT_ACTION:
			temp = data
		case KEY_ACTION:
			if i == len(steps)-1 {
				temp.(map[string]any)[step.GetKey()] = value
			} else {
				temp, ok = temp.(map[string]any)[step.GetKey()]
				if !ok {
					return data, errors.New("key not found")
				}
			}
		case IDX_ACTION:
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
