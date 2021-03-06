# protostub [![](https://travis-ci.org/arachnys/protostub.svg?branch=master)](https://travis-ci.org/arachnys/protostub) [![](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/arachnys/protostub)

A tool for generating Mypy type stubs from a Protocol Buffer definition.

## Usage
You can download a binary from the
[releases](https://github.com/arachnys/protostub/releases) page, or you can use 
[Docker](https://github.com/arachnys/protostub#docker).

Assuming that you have a `.proto` file called `foo.proto` and you want to 
generate a `.pyi` file from it, usage is as such:

```
protostub generate --proto foo.proto
```

In this case, protostub assumes you want to call your stub `foo_pb2.pyi`. If 
this is not the case, or perhaps you want to specify a different directory, the 
output can be specified like so:

```
protostub generate --proto foo.proto --mypy bar.pyi
```

Help is included in the tool:

```
Generate Mypy type stubs from Protobuf definitions

Usage:
  protostub [command]

Available Commands:
  generate    Generate a stub from a given proto file
  help        Help about any command

Flags:
  -h, --help   help for protostub

Use "protostub [command] --help" for more information about a command.
```

### Using the stubs
If you've tried using mypy in your project, you've probably noticed that it does
not do anything for protobuf types. This is because the generated python contains
nothing that Mypy can use - classes are not defined in the "normal" way.

We currently use a script to generate a list of files for Mypy to check. If the 
`.py` file has a corresponding `.pyi` file, then the `.py` is ignored. Otherwise
it is used. This allows us to override any Python with a type stub, meaning we
can override Python generated by `protoc` with something generated by `protostub`,
and get functioning type checks!

An alternative would be to invoke Mypy on a module rather than on files, and 
provide the stubs in `MYPYPATH`. However, I have not tested this, so I don't know
how well it works.


### Docker
A docker image is also provided! The easiest way for you to use this is as follows 
(assuming there is a file called foo.proto in the current directory):

```
docker run -v $(pwd):/protostub protostub:latest generate -p foo.proto
```

After doing this, you should see the help text.

## Building

### Requirements
- Go
- make (optional)


### With go get
If you already have Go all setup in your `PATH`, then it is as simple as:

```
go get github.com/arachnys/protostub/cmd/protostub
```

### With Make
This approach might be best if you're less familiar with Go, and want it to 
*just work*. It requires no messing with `$GOPATH`.

```
git clone https://github.com/arachnys/protostub
cd protostub
make
```

The protostub binary should be in the `bin` folder.

## How is this different to `mypy-protobuf`
Protostub was created and used internally at Arachnys before mypy-protobuf was
released as an open source project. We use it as part of our CI to try and catch
issues before they make it into production. As such, it's very important that
everything is dockerized - hence it's on Docker hub, so there's no dependency on
Python or protoc.

We also have out of the box Python 3 support, as that is what our codebase
required.

Our version is also completely standalone, and does not function as a plugin for
protoc. While this means you only need one binary (statically linked with no
dependencies), it does mean that it may not be quite as "correct" in terms of
parsing protobuffers, or generating python code.

## License
```
The MIT License (MIT)
Copyright © 2018 Arachnys Information Services Ltd and individual contributors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the “Software”), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```
