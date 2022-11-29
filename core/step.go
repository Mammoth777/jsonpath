package core

import (
	"fmt"
	"log"
	"regexp"
)


type Step struct {
	key string
	action actionType
	// args   []any
}

func (s *Step) GetKey() string {
	return s.key
}

func (s *Step) GetAction() actionType {
	return s.action
}

func (s *Step) Action(action actionType) *Step {
	s.action = action
	return s
}

func (s *Step) Key(key string) *Step {
	s.key = key
	return s
}

func (s *Step) String() string {
	return fmt.Sprintf("Step{action: %s, key: %s}", s.action, s.key)
}

func (s *Step) stickDot(origin string, index int) (int) {
	j := index + 1
	for ; j < len(origin); j++ {
		if origin[j] == '.' || origin[j] == '[' {
			s.Key(origin[index + 1:j])
			index = j - 1
			break
		}
	}
	if j == len(origin) {
		s.Key(origin[index + 1:])
		index = j - 1
	}
	s.Action(KEY_ACTION)
	return index
}

func (s *Step) stickBracket(origin string, index int) (int) {
	j := index + 1
	for ; j < len(origin); j++ {
		if origin[j] == ']' {
			s.Key(origin[index+1:j])
			index = j
			break
		}
	}
	if j == len(origin) {
		s.Key(origin[index + 1:])
		index = j - 1
	}
	// check action type
	var (
		regStr = `^'(.*)'$`
		regNum = `^\d+$`
	)
	if matched, _ := regexp.Match(regStr, []byte(s.key)); matched {
		s.
			Action(KEY_ACTION).
			Key(s.key[1:len(s.key) - 1])
	} else if matched, _ := regexp.Match(regNum, []byte(s.key)); matched {
		s.Action(IDX_ACTION)
	} else {
		log.Println("Invalid key: ", s.key)
	}
	return index
}

func NewStep() *Step {
	return &Step{}
}
