package jparse

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
)

type Obj struct {
	m map[string]interface{}
}

func FromReader(r io.Reader) (*Obj, error) {
	m := make(map[string]interface{})
	err := json.NewDecoder(r).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("jparse: unable to create object from reader: %w", err)
	}
	return &Obj{m: m}, nil
}

func (o *Obj) ObjectWithName(name string) (*Obj, error) {
	v, ok := o.m[name]
	if !ok {
		return nil, fmt.Errorf("jparse: %s not found", name)
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("jparse: unable to cast %s as object", name)
	}
	return &Obj{m: m}, nil
}

func (o *Obj) getInterfaceWithName(name string) (interface{}, error) {
	v, ok := o.m[name]
	if !ok {
		return nil, fmt.Errorf("jparse: %s not found", name)
	}
	return v, nil
}

func (o *Obj) getInterfaceWithPath(args ...string) (interface{}, error) {
	obj := o
	var err error
	for _, arg := range args[:len(args)-1] {
		obj, err = obj.GetObjectWithPath(arg)
		if err != nil {
			return nil, fmt.Errorf("jparse: unable to get interface with name '%s': %w", arg, err)
		}
	}
	arg := args[len(args)-1]
	v, ok := obj.m[arg]
	if !ok {
		return nil, fmt.Errorf("jparse: unable to get interface with name '%s': %w", arg, err)
	}
	return v, nil
}

func (o *Obj) GetObjectWithPath(args ...string) (*Obj, error) {
	v, err := o.getInterfaceWithPath(args...)
	if err != nil {
		return nil, fmt.Errorf("jparse: unable to find interface with path '%v'", args)
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("jparse: unable to cast '%v' to map[string]interface{}", v)
	}
	return &Obj{m: m}, nil
}

func (o *Obj) MustObjectWithPath(args ...string) *Obj {
	obj, err := o.GetObjectWithPath(args...)
	if err != nil {
		panic(err)
	}
	return obj
}

func (o *Obj) MustObjectWithName(name string) *Obj {
	obj, err := o.ObjectWithName(name)
	if err != nil {
		panic(err)
	}
	return obj
}

func (o *Obj) singleKeyWithRegex(pattern string) (string, error) {
	k := o.Keys()
	var m []string
	for _, key := range k {
		ok, err := regexp.MatchString(pattern, key)
		if err != nil {
			return "", fmt.Errorf("jparse: unable to match string: %w", err)
		}
		if ok {
			m = append(m, key)
		}
	}
	if len(m) == 0 {
		return "", fmt.Errorf("jparse: no match for pattern")
	}
	if len(m) > 1 {
		return "", fmt.Errorf("jparse: more than a single match for pattern")
	}
	return m[0], nil
}

func (o *Obj) ObjectWithRegex(pattern string) (*Obj, error) {
	name, err := o.singleKeyWithRegex(pattern)
	if err != nil {
		return nil, err
	}
	return o.ObjectWithName(name)
}

func (o *Obj) MustObjectWithRegex(pattern string) *Obj {
	obj, err := o.ObjectWithRegex(pattern)
	if err != nil {
		panic(err)
	}
	return obj
}

func (o *Obj) StringWithName(name string) (string, error) {
	v, ok := o.m[name]
	if !ok {
		return "", fmt.Errorf("jparse: %s not found", name)
	}
	s := v.(string)
	if !ok {
		return "", fmt.Errorf("jparse: unable to cast %s to string", name)
	}
	return s, nil
}

func (o *Obj) StringWithPath(args ...string) (string, error) {
	v, err := o.getInterfaceWithPath(args...)
	if err != nil {
		return "", fmt.Errorf("jparse: unable to find string with path '%v': %w", args, err)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("jparse: unable to cast %v to string", args)
	}
	return s, nil
}

func (o *Obj) MustStringWithPath(args ...string) string {
	s, err := o.StringWithPath(args...)
	if err != nil {
		panic(err)
	}
	return s
}

func (o *Obj) MustStringWithName(name string) string {
	s, err := o.StringWithName(name)
	if err != nil {
		panic(err)
	}
	return s
}

func (o *Obj) FloatWithName(name string) (float64, error) {
	v, ok := o.m[name]
	if !ok {
		return 0, fmt.Errorf("jparse: %s not found", name)
	}
	f, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("jparse: unable to cast %s to float64", name)
	}
	return f, nil
}

func (o *Obj) IntWithName(name string) (int, error) {
	v, ok := o.m[name]
	if !ok {
		return 0, fmt.Errorf("jparse: %s not found", name)
	}
	f, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("jparse: unable to cast %s to int", name)
	}
	return f, nil
}

func (o *Obj) MustFloatWithName(name string) float64 {
	f, err := o.FloatWithName(name)
	if err != nil {
		panic(err)
	}
	return f
}

func (o *Obj) SliceOfObjectWithName(name string) ([]*Obj, error) {
	v, ok := o.m[name]
	if !ok {
		return nil, fmt.Errorf("jparse: %s not found", name)
	}
	l, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("jparse: unable to cast %s to list", name)
	}
	var objs []*Obj
	for _, i := range l {
		m, ok := i.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("jparse: unable to cast list item to map[string]interface{}")
		}
		objs = append(objs, &Obj{m: m})
	}
	return objs, nil
}

func (o *Obj) SliceOfObjectWithRegex(pattern string) ([]*Obj, error) {
	name, err := o.singleKeyWithRegex(pattern)
	if err != nil {
		return nil, err
	}
	return o.SliceOfObjectWithName(name)
}

func (o *Obj) SliceOfObjectWithPath(args ...string) ([]*Obj, error) {
	v, err := o.getInterfaceWithPath(args...)
	if err != nil {
		return nil, fmt.Errorf("jparse: unable to find interface with path '%v'", args)
	}
	l, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("jparse: unable to cast %v to list", args)
	}
	var objs []*Obj
	for _, i := range l {
		m, ok := i.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("jparse: unable to cast list item to map[string]interface{}")
		}
		objs = append(objs, &Obj{m: m})
	}
	return objs, nil
}

func (o *Obj) MustSliceOfObjectWithPath(args ...string) []*Obj {
	s, err := o.SliceOfObjectWithPath(args...)
	if err != nil {
		panic(err)
	}
	return s
}

func (o *Obj) MustSliceOfObjectWithName(name string) []*Obj {
	objs, err := o.SliceOfObjectWithName(name)
	if err != nil {
		panic(err)
	}
	return objs
}

func (o *Obj) SliceOfStringWithName(name string) ([]string, error) {
	v, ok := o.m[name]
	if !ok {
		return nil, fmt.Errorf("jparse: %s not found", name)
	}
	s, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("jparse: unable to cast '%s' to []interface{}", name)
	}
	var strs []string
	for _, i := range s {
		str, ok := i.(string)
		if !ok {
			return nil, fmt.Errorf("jparse: unable to cast '%s' item to string", name)
		}
		strs = append(strs, str)
	}
	return strs, nil
}

func (o *Obj) SliceOfStringWithPath(args ...string) ([]string, error) {
	v, err := o.getInterfaceWithPath(args...)
	if err != nil {
		return nil, fmt.Errorf("jparse: unable to find interface with path '%v'", args)
	}
	s, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("jparse: unable to cast '%v' to []interface{}", args)
	}
	var strs []string
	for _, i := range s {
		str, ok := i.(string)
		if !ok {
			return nil, fmt.Errorf("jparse: unable to cast '%v' item to string", args)
		}
		strs = append(strs, str)
	}
	return strs, nil
}

func (o *Obj) MustSliceOfStringWithPath(args ...string) []string {
	s, err := o.SliceOfStringWithPath(args...)
	if err != nil {
		panic(err)
	}
	return s
}

func (o *Obj) Keys() []string {
	var ks []string
	for k, _ := range o.m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

func (o *Obj) KeysWithRegex(pattern string) ([]string, error) {
	var ks []string
	for k, _ := range o.m {
		ok, err := regexp.MatchString(pattern, k)
		if err != nil {
			return nil, fmt.Errorf("jparse: unable to match: %w", err)
		}
		if ok {
			ks = append(ks, k)
		}
	}
	sort.Strings(ks)
	return ks, nil
}

func (o *Obj) LoadStruct(v interface{}) error {
	data, err := json.Marshal(o.m)
	if err != nil {
		fmt.Errorf("unable to marshal: %w", err)
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("unable to unmarshal: %w", err)
	}
	return nil
}
