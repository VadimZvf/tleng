package tokenizer_buffer

import (
	"reflect"
	"testing"

	"github.com/VadimZvf/golang/source_mock"
)

func TestCreateBuffer(t *testing.T) {
	var source = source_mock.GetSimpleSource()
	buffer := CreateBuffer(source)
	if buffer.isSourceEnd != false {
		t.Errorf("Buffer interface should contain isEnd flag")
	}

	if reflect.TypeOf(buffer.value).Name() != "string" {
		t.Errorf("Buffer interface should contain value prop")
	}
}

func TestGetValue(t *testing.T) {
	var source = source_mock.GetSimpleSource()
	buffer := CreateBuffer(source)
	value := buffer.GetValue()
	if value != "" {
		t.Errorf("Buffer should return empty value on start. But receive: \"%s\"", value)
	}

	source.NextSymbolValue = 'a'
	buffer.Next()
	buffer.AddSymbol()
	value = buffer.GetValue()
	if value != "a" {
		t.Errorf("Buffer should add symbol")
	}

	source.NextSymbolValue = 'b'
	buffer.Next()
	buffer.AddSymbol()
	source.NextSymbolValue = 'c'
	buffer.Next()
	buffer.AddSymbol()
	value = buffer.GetValue()
	if value != "abc" {
		t.Errorf("Buffer receive symbol from source %s", value)
	}
}

func TestNext(t *testing.T) {
	var source = source_mock.GetSimpleSource()
	buffer := CreateBuffer(source)

	source.NextSymbolValue = 'a'
	buffer.Next()
	symbol := buffer.GetSymbol()
	if symbol != 'a' {
		t.Errorf("Buffer should return symbol")
	}

	isEnd := buffer.GetIsEnd()
	if isEnd != false {
		t.Errorf("Buffer should return isEnd")
	}

	source.IsEnd = true
	source.NextSymbolValue = 'g'

	buffer.Next()
	symbol = buffer.GetSymbol()

	if symbol == 'g' {
		t.Errorf("Buffer should't return symbol at end")
	}

	isEnd = buffer.GetIsEnd()
	if isEnd != true {
		t.Errorf("Buffer should notify about end")
	}
}

func TestGetPosition(t *testing.T) {
	var source = source_mock.GetSimpleSource()
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

func TestTrimNext(t *testing.T) {
	var source = source_mock.GetSourceMock("            wow    ")
	buffer := CreateBuffer(source)
	buffer.TrimNext()
	symbol := buffer.GetSymbol()
	if symbol != 'w' {
		t.Errorf("Buffer should trim spaces. But received: %c", symbol)
	}

	position := buffer.GetPosition()
	if position != 12 {
		t.Errorf("Buffer should save position after clear")
	}
}

func TestTrimNextWothoutSpace(t *testing.T) {
	var source = source_mock.GetSourceMock("test text")
	buffer := CreateBuffer(source)
	buffer.TrimNext()
	symbol := buffer.GetSymbol()
	if symbol != 't' {
		t.Errorf("Buffer should trim spaces")
	}

	position := buffer.GetPosition()
	if position != 0 {
		t.Errorf("Buffer should save position after clear")
	}
}

func TestClear(t *testing.T) {
	var source = source_mock.GetSimpleSource()
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
