#!/bin/bash

echo "Start tests"

echo "Root test"
echo "======================"
go test
echo ""

echo "Tokenizer buffer tests"
echo "======================"
cd tokenizer_buffer
go test
cd ..
echo ""

echo "Function declaration token"
echo "======================"
cd token_function_declaration
go test
cd ..
echo ""

echo "Return declaration token"
echo "======================"
cd token_return
go test
cd ..
echo ""

echo "Variable declaration token"
echo "======================"
cd token_variable_declaration
go test
cd ..
echo ""

echo "Property read token"
echo "======================"
cd token_read_property
go test
cd ..
echo ""

echo "Number token"
echo "======================"
cd token_number
go test
cd ..
echo ""

echo "String token"
echo "======================"
cd token_string
go test
cd ..
echo ""

echo "Tokenizer"
echo "======================"
cd tokenizer
go test
cd ..
echo ""

echo "Parser"
echo "======================"
cd parser
go test
cd ..
echo ""

echo "Runtime"
echo "======================"
cd runtime
go test
cd ..
echo ""
