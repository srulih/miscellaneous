package calculator

import "fmt"

type Frame struct {
	table map[string]interface{}
}

func (f *Frame) AddVar(v string, val interface{}) {
	f.table[v] = val
}

func (f *Frame) GetVar(v string) (interface{}, error) {
	if val, ok := f.table[v]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("unbound variable %s", v)
}

type Env struct {
	frames       []Frame
	currentFrame *Frame
}

func (e *Env) CreateFrame() {
	e.frames = append(e.frames, Frame{table: make(map[string]interface{})})
	e.currentFrame = &e.frames[len(e.frames)-1]
}

func (e *Env) RemoveFrame() {
	e.frames = e.frames[:len(e.frames)-1]
	e.currentFrame = &e.frames[len(e.frames)-1]
}

//TODO should we check if any frames exist?
func (e *Env) AddVar(v string, val interface{}) {
	e.currentFrame.AddVar(v, val)
}

func (e *Env) GetVar(v string) (interface{}, error) {
	var res interface{}
	var err error
	for i := len(e.frames) - 1; i >= 0; i-- {
		res, err = e.frames[i].GetVar(v)
		if err == nil {
			return res, nil
		}
	}

	return nil, err
}
