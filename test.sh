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
