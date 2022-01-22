package runtime

import (
	"testing"

	"github.com/VadimZvf/golang/parser"
	"github.com/VadimZvf/golang/runtime_bridge_mock"
	"github.com/VadimZvf/golang/source_mock"
	"github.com/VadimZvf/golang/stdout_mock"
)

func TestPrintText(t *testing.T) {
	var bridge, err = runCode("print(\"Hello World!\")")

	if err != nil {
		t.Errorf("Code failed with error: \"%s\"", err.Error())
	}

	if bridge.GetLastPring() != "Hello World!" {
		t.Errorf("Code should print message \"Hello World!\", but received: \"%s\"", bridge.GetLastPring())
	}
}

func TestCallFunction(t *testing.T) {
	var bridge, err = runCode(`
	function foo() {
		print("Hello Foo!")
	}
	
	foo()
	`)

	if err != nil {
		t.Errorf("Code failed with error: \"%s\"", err.Error())
	}

	if bridge.GetLastPring() != "Hello Foo!" {
		t.Errorf("Code should print message \"Hello Foo!\", but received: \"%s\"", bridge.GetLastPring())
	}
}

func TestFabricFunction(t *testing.T) {
	var bridge, err = runCode(`
	function welcomeFabric(prefix) {
		return function welcome(name) {
			print(prefix + " " + name + "!")
		}
	}
	
	var welcome = welcomeFabric("Hi")
	welcome("Go")
	`)

	if err != nil {
		t.Errorf("Code failed with error: \"%s\"", err.Error())
	}

	if bridge.GetLastPring() != "Hi Go!" {
		t.Errorf("Code should print message \"Hi Go!\", but received: \"%s\"", bridge.GetLastPring())
	}
}

func TestCurryingFunction(t *testing.T) {
	var bridge, err = runCode(`
	function sum(first) {
		return function sum(second) {
			print(first + second)
		}
	}
	
	sum(5)(8)
	`)

	if err != nil {
		t.Errorf("Code failed with error: \"%s\"", err.Error())
	}

	if bridge.GetLastPring() != "13.000000" {
		t.Errorf("Code should print message \"13.000000\", but received: \"%s\"", bridge.GetLastPring())
	}
}

func TestFunctionClosure(t *testing.T) {
	var bridge, err = runCode(`
	var first = 4
	var second = 10

	function bar() {
		var second = 3
		print(first + second)
	}
	
	bar()
	`)

	if err != nil {
		t.Errorf("Code failed with error: \"%s\"", err.Error())
	}

	if bridge.GetLastPring() != "7.000000" {
		t.Errorf("Code should print message \"7.000000\", but received: \"%s\"", bridge.GetLastPring())
	}
}

func TestCallbackFunction(t *testing.T) {
	var bridge, err = runCode(`
	function baz(cb) {
		cb()
	}
	
	function callback() {
		print("Hi callback")
	}

	baz(callback)
	`)

	if err != nil {
		t.Errorf("Code failed with error: \"%s\"", err.Error())
	}

	if bridge.GetLastPring() != "Hi callback" {
		t.Errorf("Code should print message \"Hi callback\", but received: \"%s\"", bridge.GetLastPring())
	}
}

func TestParenthesizedExpression(t *testing.T) {
	var bridge, err = runCode(`
	var first = 4
	var second = 10

	var a = (
		(first + 10) / (second - 8)
	) + 5
	
	print(a)
	`)

	if err != nil {
		t.Errorf("Code failed with error: \"%s\"", err.Error())
	}

	if bridge.GetLastPring() != "12.000000" {
		t.Errorf("Code should print message \"12.000000\", but received: \"%s\"", bridge.GetLastPring())
	}
}

func TestNoFunctionReferenceError(t *testing.T) {
	var _, err = runCode(`
	function welcome() {
    print("Hi")
	}

	welcome("Hi")("Tleng")
	`)

	if err == nil {
		t.Errorf("Should return a error")
	}

	if err.Error() != "Error in calling expression, here is no reference to a function" {
		t.Errorf("Should return reference error, but received: \"%s\"", err.Error())
	}
}

func runCode(code string) (*runtime_bridge_mock.Bridge, error) {
	var src = source_mock.GetSourceMock(code)
	var bridge = runtime_bridge_mock.CreateBridge()
	var stdout = stdout_mock.CreateStdout()

	var parser = parser.CreateParser(src, &stdout)

	var astRoot, astError = parser.Parse(false)

	if astError != nil {
		return &bridge, astError
	}

	var rt = CreateRuntime(&bridge)
	var runtimeErr = rt.Run(astRoot)

	if runtimeErr != nil {
		return &bridge, runtimeErr
	}

	return &bridge, nil
}
