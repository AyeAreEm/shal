# Shal

I got sick of pulling up calculator apps and using `expr` on linux, so here's a simple executable for doing calculations.<br>
Why not use a repl like Python or Node? I primarily use NixOs and don't have languages installed globally. That is also why I chose to not make this in those language. All I need is an executable with some stuff to do math.<br>

## Features
### Supported
- Addition
- Subtraction
- Multiplication
- Division
- Variables
- Binary / Octal / Hexadecimal Numbers
- Functions
- `1_000_000` Syntax
- `clear` (clear stdout)
- `exit`  (exit)

### Coming
- Modulus
- Builtin Constants (pi, e, ...)
- Builtin Functions (sqrt, log, ...)

## Why not the Calc's that already exist?
Well, I didn't know about them until after making this... Funny how one of them is also made in Go.<br>
But a quick glance at them, it doesn't look like they have variables and functions which I would like to have. I know that calculators don't usually have variables but custom functions are common.<br>
Perhaps this is more of a simple shell / repl than a calculator.

## Why is it called Shal?
Shell Calculator. Didn't want to make another `calc` and this is one of the first things I could come up with.
