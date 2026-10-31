# KrabbaScript

![GitHub License](https://img.shields.io/github/license/khytryy/krabbascript)
![GitHub top language](https://img.shields.io/github/languages/top/khytryy/krabbascript?logo=go&label=)

KrabbaScript is a simple yet powerful programming language, taking inspiration from C, Lua and Python. It is compiled, statically typed and type-safe, requiring a type for every variable declaration.

> [!CAUTION]
> This project is still W.I.P and some stuff are not finished. Check out our [Discord](https://discord.gg/MQT4YgEYvn) for news and updates

## Getting started

## Building the project
You can simply build the compiler by running `go build`

```bash
go build
```

## Syntax

This is a brief overview of the docs. Don't expect an in-depth explanation with jargon and whatnot, but you'll get a sense of Krabbascript's style. Let's begin!

As with most statically typed and compiled languages, Krabbascript enforces semicolons to separate expressions.

### Variables

Krabbascript has two options for variables. You can have variables, or, you can have values. What's the difference? Variables are mutable -- they can be mutated and reassigned -- values are not. You can declare a variable or value using `var` or `val` respectively, the name of the variable, and then a type. Here's an example:

```
var butterfly: I32 = 1704;
```

If you wanted to create a value, it would be the same approach but with `val`:

```
val butterfly: I32 = 96;
```

If you try to mutate or reassign a value, you'll be hit with an error.

Quick note on variables: mutate them all you want, but you can't touch the type. `butterfly` will always be an I32 for the rest of eternity. It's fate is sealed.

### Control structures

Also like most statically typed and compiled languages, Krabbascript uses curly brackets for its blocks. 
However, we won't make you use parens for your conditions. Your welcome.

Krabbascript has a solid arrangement of control structures for you. We have if-elsif-else blocks, while loops, repeat loops, for loops and when blocks.

#### if-elsif-else blocks

```
if breakfast {
  drinkFrenchPressCoffeeScript();
} elsif lunch {
  eatSwedishMeatballs();
} else {
  eatCrab();
}
```

#### while loops (TODO)

```
while railsIsUnpopular {
  rantAboutWhyRailsIsGreat();
}
```

#### repeat loops (TODO)

In these, you loop now and declare a condition later.

```
repeat {
  rantAboutWhyRailsIsGreat();
} until railsIsTrending;
```

#### when blocks

This is one of the more unique parts of Krabbascript. This is like a switch in nature, but with a slightly different approach.

```
when language {
  == "Elixir" do celebrate();
  == "Haskell" do complain();
}
```

#### for loops (TODO)

Krabbascript's for loops are closer to foreach loops, and the design was heavily inspired by Lua.

```
for index, language in languages {
    when language {
    == "Elixir" {
      celebrate();
    }
    == "Haskell" {
      complain();
    }
  }
  Std.print(index);
}
```

If you don't want that pesky index variable, feel free to use an underscore in its place. Like:

```
for _, variable in arrays
```

### Types

Krabbascript provides 13 types. Remember that all types are capitalized like nouns.

- Str
- I64
- I32
- I16
- I8
- U64
- U32
- U16
- U8
- Any
- Bool
- F64
- F32

### Structs

You can define a struct like this:

```
struct User {
	name: Str,
	birth: U64
}
```

All struct names must be capitalized; if you were to name a struct something like `user`, you would be met with an error.

To access a struct's property, use dot notation. For example:

```
User.name
```

Here's how you can make an instance of a struct:

```
val my_user: User = User{
      name  = “Joe Doe”,
      birth = 1704
};
```

### Functions

Functions must be declared with a return type. You declare a function like this:

```
func giveMeACity() -> Str {
  return "Tampere";
}
```

If you wish to return void, you don't have to type out `-> void` -- that would be awful! You can simply omit the return type and that symbolizes the function returning void.

You can call a function like in any other language:

```
giveMeACity();
```

### Comments

Comments are quite simple. They use a single hashtag: `#`. There are no multi-line comments so far.

```
# Finland gained independence from Russia in 1917.
```

### Modules

Like structs, module names must be capitalized. You can have multiple modules in a file. They are declared like this:

```
mod Shapes {
	struct Rectangle {
	  x: I32;
	  y: I32;
	  w: I32;
	  h: I32;
  }

	func rect_area(rect: Rectangle) -> I32 {
		return rect.w * rect.h;
  }
}
```

### Importing

This is how you import a module:

```
import Shapes;
```

## Conclusion

The website is a work in progress, but it will host more details regarding Krabbascript.
