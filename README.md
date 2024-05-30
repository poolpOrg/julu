# JuLu Programming Language Specification


**WARNING:**
**THIS IS A WORK IN PROGRESS, IT DOES NOT WORK, DO NOT BOTHER BUILDING.**


## Overview

JuLu is a programming language designed for a short learning curve,
targeting system development with efficiency comparable to C.
It combines elements of Python and Golang,
allowing for low-level operations,
including inlining assembly code,
crafting structures,
and memory layout management.


## "Hello world !"
```go
package main

fn main => println("Hello world !")
```


## Types

- Strictly typed
- Basic types:
  - int
  - float
  - int8
  - int16
  - int32
  - int64
  - uint8
  - uint16
  - uint32
  - uint64
  - float32
  - float64
  - char
  - bool
  - string
- Composite types through structs and unions
- Type inference


## Functions
Support for returning multiple values matching a return signature

```go
fn foobar(x : int, y : int) -> (int) {
    return x+y
}

fn foobar(x : int, y : int) -> (int, int) {
    return x, y
}
```


## Conditionals

- if / else if / else
```go
// long form for conditional code blocks
if true {
    println("do this !")
}

if x {
    println("do this !")
} else if y {
    println("or that !")
}

if x {
    println("do this !")
} else if y {
    println("or that !")
} else {
    println("or even that !")
}

// short form for single expression / statements
if true => println("only do this !")

if x => println("only do this !")
else if y => println("don't do that !")

if x => println("x!")
else if y => println("y!")
else => println("z!")


// long and short forms can be mixed:
if !ok => return -1
else => {
    println("yeah !")
    return 0
}

if ok => {
    println("yeah !")
    return 0
} else => return -1
```



- match:
```go
match x {
    case x==1 => println("match")
    case x!=1 => println("mismatch")
} else => println("found no match!")

match x {
    case x==1 {
        println("match")
    }
    case x!=1 {
        println("mismatch")
    }
} else => println("found no match!")


```

## Loops

- loop:
```go
loop {
    println("infinite loop")
}
```

- while:
```go
x = 0
while x < 42 {
    println("loop while x < 42")
    x++
}
```


- until:
```go
x = 42
until x > 42 {
    println("loop until x > 42")
    x++
}
```

- for x in y {}
```go
for x in [1, 2, 3, 4] => println(x * 2)

for x in [1, 2, 3, 4] {
    println(x * 2)
    println(x << 1)
}
```

## Concurrency

- select statement for multiplexing
- channels
- tasks (similar to goroutines)


## Context management

- Resource allocation with deferred release upon block termination

```
with open() as fp => fp.write("foobar")
// fp closes at the end of this block
```

## Comments
- single line `#` or `//`
- multi-line `/* */`


## Operators

- All C operators: +, -, *, +, %, ~, &, |, <<, >>
- Circular shift: <<< and >>>
- Power: **
- Logical operators: &&, ||, ! can also be written as and, or, not


## Strings

- regular strings
- raw strings (no escaping)
- f-strings (strings with expression expansion)


## Pointers

- allow manipulation of pointers when necessary


## Method binding

- methods can be bound to structs:

```
struct File {
    fd: int
}

extern fn open(filename: string, mode: string) -> File
extern fn read(fd: int, size: int) -> string
extern fn close(fd: int)

File => {
    fn read(self, size: int) -> string {
        return read(self.fd, size)
    }

    fn close(self) {
        close(self.fd)
    }
}
```
- self is a keyword used for method binding withing structs

## Language keywords
- fn
- let
- mut
- extern
- return
- if
- else
- else if
- switch
- case
- default
- for
- in
- range
- loop
- while
- break
- continue
- with
- as
- true
- false
- struct
- union
- type
- int
- int8
- int16
- int32
- int64
- uint8
- uint16
- uint32
- uint64
- float
- float32
- float64
- char
- bool
- string
- select
- task
- import
- package
- self
- chan
- <-
