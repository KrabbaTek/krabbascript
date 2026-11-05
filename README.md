# KrabbaScript

![GitHub License](https://img.shields.io/github/license/khytryy/krabbascript)
![GitHub top language](https://img.shields.io/github/languages/top/khytryy/krabbascript?logo=go&label=)

KrabbaScript is a simple yet powerful programming language, taking inspiration from OCaml and Golang. 

It is compiled, statically typed and type-safe, requiring a type for every variable declaration.

> [!CAUTION]
> This project is still W.I.P and some features are not finished. Check out our [Discord](https://discord.gg/MQT4YgEYvn) for the latest news and updates

## Getting started

## Building the project

You can build the compiler by running `go build` in the `krabbascript/` directory:

```bash
go build
```

Then, you can run the program by invoking `./krabbascript`.

## Syntax

This is just a brief overview of the syntax, so don't expect an in-depth explanation with jargon and whatnot, but you'll get a sense of Krabbascript's style. Let's begin!

But first, a quick tangent into semicolons:

Technically, KrabbaScript requires them, but you can write your code without them and use the `--format` flag when building to add them back in automatically.

You can also use the standalone `format` command which will only format your file.

My examples will be using them to avoid any confusion, but remember that they aren't required, so don't let semicolons be a detterence for you.

### Variables

Krabbascript has two options for variables. You can have variables, or, you can have values. 

What's the difference? Variables are mutable — they can be mutated and reassigned — values are not. You can declare a variable or value using `var` or `val` respectively, the name of the variable, and then a type. Here's an example:

```kscript
var butterfly: I32 = 1704;
```

If you wanted to create a value, it would be the same approach but with `val`:

```kscript
val butterfly: I32 = 96;
```

If you try to mutate or reassign a value, you'll be hit with an error.

Quick note on variables: mutate them all you want, but you can't touch the type. `butterfly` will always be an I32 for the rest of eternity; Its fate is sealed.

### Control structures

Also like most statically typed and compiled languages, Krabbascript uses curly brackets for its blocks.

However, we won't make you use parens for your conditions — we're not that crazy.

Krabbascript has a pretty hefty arrangement of control structures for you. We have if-elsif-else blocks, while loops, repeat loops, for loops, foreach loops, and when blocks.

#### if-elsif-else blocks

```kscript
if breakfast {
  drinkFrenchPressCoffeeScript!;
} elsif lunch {
  eatSwedishMeatballs!;
} else {
  eatCrab!;
}
```

#### while loops

```kscript
while railsIsUnpopular {
  rantAboutWhyRailsIsGreat!;
}
```

#### repeat loops

In these, you loop now and declare a condition later.

```kscript
repeat {
  rantAboutWhyRailsIsGreat!;
} until railsIsTrending;
```

If you're having trouble understanding these, think about them as reverse while loops.

The example above is just:

```
while !railsIsTrending {
  rantAboutWhyRailsIsGreat!;
}
```

One piece of advice is to avoid negated variable names —
`while not railsIsTrending` is much more painful to read than
`until railsIsTrending`. 

#### when blocks

This is one of the more unique parts of Krabbascript. This is like a switch in nature, but with a slightly different approach.

```kscript
when language {
  == "Elixir" do celebrate!;
  == "Haskell" do complain!;
}
```

#### for loops

```kscript
for 96 times {
  Std.print("Are we there yet?");
}
```

If you want an index, then name your variable within parens like so:

```kscript
for 96 times (index) {
  Std.print(index);
}
```

#### foreach

```kscript
for element in array {
  Std.print(element)
}
```

If you want the index then you can pull a variable down like with for loops:

```kscript
for element in array (index) {
  Std.print(element); Std.print(index)
}
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

Note: while `Bool` is capitalized, you still use normal `true` and `false` — they aren't capitalized like in Python.

### Structs

You can define a struct like this:

```kscript
struct User {
	name: Str,
	birth: U64
}
```

All struct names must be capitalized; if you were to name a struct something like `user`, you would be met with an error.

To access a struct's property, use dot notation. For example:

```kscript
User.name
```

Here's how you can define an instance of a struct:

```kscript
val my_user: User = User{
  name  = “Joe Doe”,
  birth = 1704
};
```

### Functions

Functions must be declared with a return type. You declare a function like this:

```kscript
func giveMeACity() -> Str {
  return "Tampere";
}
```

If you wish to return void, you don't have to type out `-> void` — that would be awful! Instead, you can simply omit the return type and that symbolizes the function returning void.

You can call empty functions using either a bang or empty parens:

```kscript
giveMeACity!;
```

Or

```kscript
giveMeACity();
```

### Comments

Comments are quite simple. They use a single hashtag: `#`. There are no multi-line comments so far.

```kscript
# Finland gained independence from Russia in 1917.
```

### Modules

Like structs, module names must be capitalized. You can have multiple modules in a file. They are declared like this:

```kscript
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

```kscript
import Shapes;
```

## Ecosystem

### KrabKytta

[link](https://github.com/KrabbaTek/KrabKytta)

KrabKytta is a formatter for KrabbaScript — the name roughly translates to Krab Cop in Finnish.

One thing to note about KrabKytta is that it will not change your code — it will just tell you where you broke the rules.

You can either run KrabKytta with a config file to add some
personalization to it, or you can run it without one and it will use sensible default settings.

### Cranium

[link](https://github.com/KrabbaTek/Cranium-cli)

Cranium is an opinionated package manager prioritizing ease-of-use and human satisfaction.

While Cranium lacks features like dependency resolution, it has an amazing UX, offering things like detailed auto-correct and a litany of aliases.

#### Cranium Off the Grid

[link](https://github.com/KrabbaTek/Cranium-off-the-grid)

I forked Cranium and took it off the grid because
I wanted a way to easily create packages that can be called from anywhere on your machine without a full package manager.

### Nihilist

Nihilist is a preprocessor for KrabbaScript aimed to make code more expressive and readable.

The key feature of Nihilist is omitting parens in function calls.

Nihilist also comes with an ASI, which is ran for all processing,
unlike in KrabbaScript.

```
add 2, 2

# => add(2, 2);
```