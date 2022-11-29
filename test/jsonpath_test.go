package test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/Mammoth777/jsonpath/core"

	"github.com/stretchr/testify/assert"
)

var origin = `
	{
		"kind": "Pod",
		"apiVersion": "v1",
		"metadata": {
			"name": "nginx",
			"namespace": "default",
			"labels": {
				"app": "nginx"
			}
		},
		"spec": ["a", "b", "c"],
		"arr": [
			[1,2,3],
			[4,5,6]
		]
	}
`
var data any

func init() {
	err := json.Unmarshal([]byte(origin), &data)
	if err != nil {
		log.Panic(err)
	}
}

func TestCompileObj(t *testing.T) {
	steps := core.Compile("$.metadata.name")
	for _, s := range steps {
		log.Println(s.String())
	}
	if len(steps) != 5 {
		t.Error("Compile failed")
	}
}

func TestCompileDynamicObj(t *testing.T) {
	steps := core.Compile("$.metadata['name']")
	for _, s := range steps {
		log.Println(s.String())
	}
	should := assert.New(t)
	should.Equal(len(steps), 3)
	should.Equal(steps[2].GetKey(), "name")
	should.Equal(steps[2].GetAction(), core.KEY_ACTION)
}

func TestCompileArray(t *testing.T) {
	steps := core.Compile("$.metadata['name'][34]")
	for _, s := range steps {
		log.Println(s.String())
	}
	should := assert.New(t)
	should.Equal(len(steps), 4)
	should.Equal(steps[2].GetKey(), "name")
	should.Equal(steps[2].GetAction(), core.KEY_ACTION)
	should.Equal(steps[3].GetKey(), "34")
	should.Equal(steps[3].GetAction(), core.IDX_ACTION)
}

func TestReadObj(t *testing.T) {
	should := assert.New(t)
	value, err := core.Read(data, "$.metadata['name']")
	log.Println(value, "value")
	should.Nil(err)
	should.Equal(value, "nginx")
}

func TestReadArray(t *testing.T) {
	value, err := core.Read(data, "$.spec[1]")
	log.Println(value, "value")
	should := assert.New(t)
	should.Nil(err)
	should.Equal(value, "b")
}

func TestReadArray2(t *testing.T) {
	value, err := core.Read(data, "$.arr[1][1]")
	log.Println(value, "value")
	should := assert.New(t)
	should.Nil(err)
	should.Equal(value, float64(5))
}

func TestWriteObj(t *testing.T) {
	should := assert.New(t)
	core.Write(data, "$.metadata['name']", "nginx2")
	value, err := core.Read(data, "$.metadata['name']")
	log.Println(value, "value")
	should.Nil(err)
	should.Equal(value, "nginx2")
}

func TestWriteList(t *testing.T) {
	should := assert.New(t)
	core.Write(data, "$.arr[1][1]", "a")
	value, err := core.Read(data, "$.arr[1][1]")
	should.Nil(err)
	should.Equal(value, "a")
}
