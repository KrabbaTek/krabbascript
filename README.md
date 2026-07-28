# KrabbaScript
![GitHub License](https://img.shields.io/github/license/khytryy/Krabbascript)
![GitHub top language](https://img.shields.io/github/languages/top/khytryy/Krabbascript?logo=c&label=)

<div style="center">
    <img width="45%" height="45%" src=".github/logo/kscript_logo.png">
</div>

KrabbaScript is a statically typed, AOT compiled, high level programming language designed to be easy to write and read.

> [!CAUTION]
> This project is still in its early stages, so expect a lot of unfinished stuff. Check-out our [Discord](https://discord.gg/MQT4YgEYvn) for news and progress on KrabbaScript!

KrabbaScript mixes C and Python together to form simple, and easy to learn syntax.

## Getting started

### Prerequisites

#### Git (optional)

Git is for cloning KrabbaScript's source code.
You can check if Git is installed with:

```bash
git --version
```

You can then clone the repo with:

```bash
git clone https://github.com/KrabbaTek/Krabbascript
```

#### Make

To build KrabbaScript, you have to use [premake5](https://premake.github.io/download).

#### C compiler

KrabbaScript is written in C, and requires a C compiler (depends on your platform and solution file) to build.

### Building KrabbaScript
If you run `premake5 --help`, premake will show you available generators like `vs2026` or `xcode4`.

#### Windows (W.I.P)

To build for Windows, you will need Visual Studio 2008 or higher. Latest recommended.
By running:

```bash
premake5 vs<EDITION>
```

where `EDITION` is the edition you have, premake will generate you a solution file for the project inside `build/`. From there build it within Visual Studio.

#### Linux

To build for Linux, you will need GNU make or compatible.
By running:

```bash
premake5 gmake
```

premake will generate you a Makefile inside `build/`. By default, when running `make` it will build the compiler __and__ the tests for it. You can avoid this by doing `make KrabbaScript` instead.

#### macOS

To build for macOS, you'll have to research yourself. Yeah sorry, I don't have a mac to tell you how to build the project with XCode. But what I __can__ tell you is that by running:

```bash
premake5 xcode4
```

premake will generate you an XCode 4 solution. Now go figure it out yourself! (And maybe write down how to do it here?)


## Libraries used
- [Unity](https://github.com/ThrowTheSwitch/Unity) by ThrowTheSwitch

Licenses can be found in `3rd_party/xyz/LICENSE.md`
