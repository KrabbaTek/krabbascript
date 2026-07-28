# 02-Getting started

## Installing KrabbaScript
To install KrabbaScript, you need the following:
- GNU make
- GCC
- Git
- A good mindset :)

Once you got all of them, you can start by cloning the repo. Run this in your UNIX-like terminal™ and then cd into the directory:
```git clone https://github.com/khytryy/krabbascript```

Now build the source code, run:
```make```

This will build the executable at the project root called kscript.

Now run the following to install KrabbaScript (or to copy it to /usr/bin):
```sudo make install```

You're all set now! Now create a new project in whatever location you would with the following command:
```kscript init xyz```

Replace ``xyz`` with the name of your project. Once you will open the project and ls it, you will see the following structure:

- src/main.krs
- manifest.toml

Now let me explain what each of these files and directories have for a role.
``src`` folder is where your code lives. It contains .krs files, KrabbaScript source files.

``manifest.toml`` tells the compiler stuff like wheres the source folder located, where to build our binary, the name of our project, its version and more info like compile flags.

Open up the project root in your code editor and check out ``src/main.krs``, the code will look like this:

```
from <stdlib.krs> import *;

func main() -> i32 {
    std.print("Hello, world!");

    return 0;
}
```

Let's walk through what each statement means.

## from <xyz> import abc;
From brings in the specified module and it's dependencies into the scope.

Sometimes in your code you need to import some other libraries, like the standard library. To import a library, you do ``from <xyz.krs> import abc;``, where xyz is the file and abc is the module you are trying to import. The compiler then imports that module and it's dependencies into your code. You can also bring in everything in that file into the scope, by replacing abc with *.

## func xyz() -> type { ... }
This is how you define a function, where xyz is the function identifier and type is the return type. There are a lot of types in KrabbaScript, one of them is u32 (unsigned 32-bit number). Most of the number types follow this syntax:

- i/u/f: integer (signed), unsigned, float
- 8-64: bit-width

Keep in mind that theres only f32 and f64 (or a double in C) for floats.

You can also define a function without a return type (void). Simply write ``func xyz() {}`` where xyz is the function name.

## std.print(str);
This is a function part of the std module, the standard library. Every function from a module must start with the module name, a dot and the actual function you're trying to call.

## return x;
The return statement returns a value, note that because KrabbaScript is statically typed, the value you're trying to return must match the functions return type.

For functions with no return type, you dont have to return anything, but doing ``return;`` will return from the function and jump back to the code where it was called and advanced to the the next statement.

Now compile your epic hello world program and run it with ``kscript build .`` (if you're in another directory just do ``kscript build /path/to/project/``) and ``./xyz/abc`` where xyz is the path to the binary folder in your project and abc is the programs name.

Next > [Variables](/DOCUMENTATION/03-Variables/README.md)