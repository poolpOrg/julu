# JuLu Programming Language Specification

## Overview

JuLu is a programming language designed for a short learning curve,
targeting system development with efficiency comparable to C.
It combines elements of Python and Golang,
allowing for low-level operations,
including inlining assembly code,
crafting structures,
and memory layout management.


## Syntax

```
fn main => println("Hello world !")

// Alternatively
fn main() {
    println("Hello world !")
}
```

Parentheses are optional when there are no parameters,
and `=>` indicates a block with a single expression.


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


## Loops

- for i in range { }
- loop { }
- loop condition { }


## Conditionals

- switch case
- if
- else if
- else


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

- All C operators
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
