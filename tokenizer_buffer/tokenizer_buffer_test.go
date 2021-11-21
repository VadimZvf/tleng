package tokenizer_buffer

import (
	"reflect"
	"testing"
)

type SourceMock struct {
	isEnd      bool
	nextSymbol string
}

func (source *SourceMock) NextSymbol() (symbol string, isEnd bool) {
	return source.nextSymbol, source.isEnd
}

func GetSourceMock() *SourceMock {
	return &SourceMock{false, "a"}
}

func TestCreateBuffer(t *testing.T) {
	var source = GetSourceMock()
	buffer := CreateBuffer(source)
	if buffer.isEnd != false {
		t.Errorf("Buffer interface should contain isEnd flag")
	}

	if reflect.TypeOf(buffer.value).Name() != "string" {
		t.Errorf("Buffer interface should contain value prop")
	}
}

func TestGetValue(t *testing.T) {
	var source = GetSourceMock()
	buffer := CreateBuffer(source)
	value := buffer.GetValue()
	if value != "" {
		t.Errorf("Buffer should return embty value on start")
	}

	buffer.Next()
	buffer.AddSymbol()
	value = buffer.GetValue()
	if value != "a" {
		t.Errorf("Buffer should add symbol")
	}

	source.nextSymbol = "b"
	buffer.Next()
	buffer.AddSymbol()
	source.nextSymbol = "c"
	buffer.Next()
	buffer.AddSymbol()
	value = buffer.GetValue()
	if value != "abc" {
		t.Errorf("Buffer receive symbol from source %s", value)
	}
}

func TestNext(t *testing.T) {
	var source = GetSourceMock()
	buffer := CreateBuffer(source)

	buffer.Next()
	symbol := buffer.GetSymbol()
	if symbol != "a" {
		t.Errorf("Buffer should return symbol")
	}

	isEnd := buffer.GetIsEnd()
	if isEnd != false {
		t.Errorf("Buffer should return isEnd")
	}

	source.isEnd = true
	source.nextSymbol = "g"

	buffer.Next()
	symbol = buffer.GetSymbol()

	if symbol != "g" {
		t.Errorf("Buffer should return symbol")
	}

	isEnd = buffer.GetIsEnd()
	if isEnd != true {
		t.Errorf("Buffer should notify about end")
	}
}

func TestGetPosition(t *testing.T) {
	var source = GetSourceMock()
	buffer := CreateBuffer(source)
	positin := buffer.GetPosition()
	if positin != 0 {
		t.Errorf("Buffer should return start position 0")
	}

	buffer.Next()
	buffer.Next()
	positin = buffer.GetPosition()
	if positin != 2 {
		t.Errorf("Buffer should return start position 2")
	}

	buffer.Next()
	positin = buffer.GetPosition()
	if positin != 3 {
		t.Errorf("Buffer should return start position 3")
	}
}

type SourceMockWithContent struct {
	fullText string
	position int
}

func (source *SourceMockWithContent) NextSymbol() (symbol string, isEnd bool) {
	source.position = source.position + 1

	if source.position >= len(source.fullText) {
		return "", true
	}

	return string(source.fullText[source.position]), false
}

func TestTrimNext(t *testing.T) {
	var source = SourceMockWithContent{
		fullText: "      \n     wow    ",
		position: -1,
	}
	buffer := CreateBuffer(&source)
	buffer.TrimNext()
	symbol := buffer.GetSymbol()
	if symbol != "w" {
		t.Errorf("Buffer should trim spaces")
	}

	position := buffer.GetPosition()
	if position != 12 {
		t.Errorf("Buffer should save position after clear")
	}
}

func TestTrimNextWothoutSpace(t *testing.T) {
	var source = SourceMockWithContent{
		fullText: "test text",
		position: -1,
	}
	buffer := CreateBuffer(&source)
	buffer.TrimNext()
	symbol := buffer.GetSymbol()
	if symbol != "t" {
		t.Errorf("Buffer should trim spaces")
	}

	position := buffer.GetPosition()
	if position != 0 {
		t.Errorf("Buffer should save position after clear")
	}
}

func TestClear(t *testing.T) {
	var source = GetSourceMock()
	buffer := CreateBuffer(source)
	buffer.Next()
	buffer.Next()
	buffer.Next()
	buffer.Next()
	buffer.Clear()
	value := buffer.GetValue()
	if value != "" {
		t.Errorf("Buffer should clear value")
	}

	position := buffer.GetPosition()
	if position != 4 {
		t.Errorf("Buffer should save position after clear")
	}
}
