# TLENg
Useless programming language

## Example

```js
function welcome(name) {
    print("Hello " + name)
}

var user = "World"

welcome(user)
```
## Demo page
[Demo](https://vadimzvf.github.io/tleng/)

## Syntax
Variable declaration
```js
var a = 1
```
Function declaration
```js
function summ(a, b) {
    return a + b
}
```

## Data types
Integer number only
```js
var a = 130
```
String
```js
var b = "Some value"
```
Function. This example will print: `3.000000`
```js
function foo(inner) {
    return inner() + 1
}

function innerFunc() {
    return 2
}

var a = foo(innerFunc)

print(a)
```
Unknown
```js
var b
```

## Build in methods
Log values
```js
print("Some text")
```