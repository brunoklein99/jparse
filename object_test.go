package jparse

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetObjectWithPathMultiLevelExpression(t *testing.T) {
	m := map[string]interface{}{
		"aaa": map[string]interface{}{
			"bbb": map[string]interface{}{
				"ccc": map[string]interface{}{
					"what": "42",
				},
			},
		},
	}
	obj := &Obj{m: m}
	obj, err := obj.GetObjectWithPath("aaa", "bbb", "ccc")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "42", obj.MustStringWithName("what"))
}

func TestGetObjectWithPathSingleLevelExpression(t *testing.T) {
	m := map[string]interface{}{
		"aaa": map[string]interface{}{
			"what": "42",
		},
	}
	obj := &Obj{m: m}
	obj, err := obj.GetObjectWithPath("aaa")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "42", obj.MustStringWithName("what"))
}
