# Krabbascript
![GitHub License](https://img.shields.io/github/license/khytryy/krabbascript)
![GitHub top language](https://img.shields.io/github/languages/top/khytryy/krabbascript?logo=c&label=)

<div style="center">
    <img width="45%" height="45%" src=".github/logo/kscript_logo.png">
</div>

KrabbaScript is a statically typed, AOT compiled, high level programming language designed to be easy to write and read.

> [!CAUTION]
> This project is still in its early stages, so expect a lot of unfinished stuff. Check-out our [Discord](https://discord.gg/MQT4YgEYvn) for news and progress on Krabbascript!

Krabbascript mixes C and Python together to form simple, and easy to learn syntax.

## Getting started

### Prerequisites

#### Git (optional)

Git is for cloning Krabbascript's source code.
You can check if Git is installed with:

```bash
git --version
```

You can then clone the repo with

```bash
git clone https://github.com/KrabbaTek/krabbascript
```

#### Make

To build Krabbascript, you have to use GNU Make or compatible.

#### C compiler

Krabbascript is written in C, and requires a C compiler do build. I genuinely recommend that you use gcc to build the compiler, because its set as default in the Makefile. If you cant use gcc for some reason, you can edit CC variable to the compiler you want to build the compiler with.

### Building Krabbascript

```bash
make && sudo make install
```

This creates the final binary in the project's directory and copies the compiler to `/usr/local/bin/kscript`.

## Libraries used
- [Unity](https://github.com/ThrowTheSwitch/Unity) by ThrowTheSwitch

Licenses can be found in `3rd_party/xyz/LICENSE.md`