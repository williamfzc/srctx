# srctx: source context

A library for extracting and analyzing definition/reference graphs from your codebase. Powered by tree-sitter and LSIF/SCIP.

| Name           | Status                                                                                                                                            |
|----------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| Latest Version | ![GitHub release (latest by date)](https://img.shields.io/github/v/release/williamfzc/srctx)                                                      |
| Unit Tests     | [![Go](https://github.com/williamfzc/srctx/actions/workflows/ci.yml/badge.svg)](https://github.com/williamfzc/srctx/actions/workflows/ci.yml)     |
| Code Coverage  | [![codecov](https://codecov.io/github/williamfzc/srctx/branch/main/graph/badge.svg?token=1DuAXh12Ys)](https://codecov.io/github/williamfzc/srctx) |
| Code Style     | [![Go Report Card](https://goreportcard.com/badge/github.com/williamfzc/srctx)](https://goreportcard.com/report/github.com/williamfzc/srctx)      |

## About this tool

This library processes your code into precise function-level graphs, just like an IDE, and then you can apply some
analysis on them.

<img width="1064" alt="image" src="https://github.com/williamfzc/srctx/assets/13421694/aa3a37ed-b7b8-4703-90f7-d61b34065994">

This lib originally was designed for monitoring the impacts of each commits.

If we need to know the impact of git line changes on the entire repository, we can only rely on manual reading of the
code to evaluate.

With this lib developers can know exactly what happened in every lines of your code. Such as definition, reference.

```bash
./srctx diff --outputHtml output.html
```

Some "dangerous" line changes can be found automatically.

<img width="843" alt="image" src="https://github.com/williamfzc/srctx/assets/13421694/e6e48e67-35b1-4c52-aa6a-99b1ba2f02db">

You can see a dangerous change in file `cmd/srctx/diff/cmd.go#L29-#143`, .

![](https://user-images.githubusercontent.com/13421694/236666915-5d403e4a-9cc1-4364-afbe-363cf82e5e49.png)

Or you prefer a text report?

We hope to utilize the powerful indexing capabilities of LSIF to quantify and evaluate the impact of text changes on the
repository, reducing the mental burden on developers.

# Usage

## Usage as Cli (Recommendation)

### For Golang

We have embedded the lsif-go indexer in our prebuilt binary files. So all you need is:

```bash
wget https://github.com/williamfzc/srctx/releases/download/v0.6.0/srctx-linux-amd64
chmod +x srctx-linux-amd64
./srctx-linux-amd64 diff --withIndex --before HEAD~1 --after HEAD --lsif dump.lsif --outputCsv output.csv --outputDot output.dot
```

### For Python

```bash
pip3 install --upgrade git+https://github.com/sourcegraph/lsif-py.git
lsif-py . --file ./dump.lsif

wget https://github.com/williamfzc/srctx/releases/download/v0.6.0/srctx-linux-amd64
chmod +x srctx-linux-amd64
./srctx-linux-amd64 diff --before HEAD~1 --after HEAD --lsif dump.lsif --outputCsv output.csv --outputDot output.dot
```

### For Java/Kotlin

Linux only. If you're using other platforms, please see [scip-java](https://sourcegraph.github.io/scip-java/docs/getting-started.html#run-scip-java-index) for details.

```bash
wget https://github.com/williamfzc/srctx/releases/download/v0.6.0/srctx-linux-amd64-full.zip
unzip srctx-linux-amd64-full.zip

# https://sourcegraph.github.io/scip-java/docs/getting-started.html#run-scip-java-index
./scip-java index -- clean assembleDebug
```

This bash will create a scip file for you. Then:

```bash
wget https://github.com/williamfzc/srctx/releases/download/v0.6.0/srctx-linux-amd64
chmod +x srctx-linux-amd64
./srctx-linux-amd64 diff --before HEAD~1 --after HEAD --scip index.scip --outputCsv output.csv --outputDot output.dot
```

### For JavaScript

https://github.com/sourcegraph/scip-typescript

### For Other Languages

#### 1. Generate LSIF file

Tools can be found in https://lsif.dev/ .

You will get a `dump.lsif` file after that.

#### 2. Run `srctx`

Download our prebuilt binaries from [release page](https://github.com/williamfzc/srctx/releases).

For example, diff from `HEAD~1` to `HEAD`:

```bash
./srctx diff --before HEAD~1 --after HEAD --lsif dump.lsif --outputCsv output.csv --outputDot output.dot
```

See details with `./srctx diff --help`.

## Usage as Github Action (Recommendation)

Because LSIF files require dev env heavily, it's really hard to provide a universal solution in a single binary file for
all the repos.

<img width="697" alt="image" src="https://user-images.githubusercontent.com/13421694/236666915-5d403e4a-9cc1-4364-afbe-363cf82e5e49.png">

We are currently working on [diffctx](https://github.com/williamfzc/diffctx), which will provide a GitHub Actions plugin
that allows users to use it directly in a Pull Request.

## Usage as Lib

### API

Our built-in diff implementation is a good example. [cmd/srctx/diff/cmd.go](cmd/srctx/diff/cmd.go)

### Low level API

Low level API allows developers consuming LSIF file directly.

```golang
yourLsif := "../parser/lsif/testdata/dump.lsif.zip"
sourceContext, _ := parser.FromLsifFile(yourLsif)

// all files?
files := sourceContext.Files()
log.Infof("files in lsif: %d", len(files))

// search definition in a specific file
defs, _ := sourceContext.DefsByFileName(files[0])
log.Infof("there are %d def happend in %s", len(defs), files[0])

for _, eachDef := range defs {
    log.Infof("happened in %d:%d", eachDef.LineNumber(), eachDef.Range.Character)
}

// or specific line?
_, _ = sourceContext.DefsByLine(files[0], 1)

// get all the references of a definition
refs, err := sourceContext.RefsByDefId(defs[0].Id())
if err != nil {
    panic(err)
}
log.Infof("there are %d refs", len(refs))

for _, eachRef := range refs {
    log.Infof("happened in file %s %d:%d",
    sourceContext.FileName(eachRef.FileId),
    eachRef.LineNumber(),
    eachRef.Range.Character)
}
```

Or see [cmd/srctx/cmd_diff.go](cmd/srctx/cmd_diff.go) for a real example with git diff.

# Correctness / Accuracy

<img width="1159" alt="image" src="https://github.com/williamfzc/srctx/assets/13421694/6cfa72c2-787a-4ae6-8cef-e77c1985d307">

We wanted it to provide detection capabilities as accurate as an IDE.

# Roadmap

See [Roadmap Issue](https://github.com/williamfzc/srctx/issues/31).

# Contribution

Issues and PRs are always welcome.

# References

LSIF is a standard format for persisted code analyzer output.
Today, several companies are working to support its growth, including Sourcegraph and GitHub/Microsoft.
The LSIF defines a standard format for language servers or other programming tools to emit their knowledge about a code
workspace.

- https://lsif.dev/
- https://microsoft.github.io/language-server-protocol/overviews/lsif/overview/
- https://code.visualstudio.com/blogs/2019/02/19/lsif#_how-to-get-started

# Thanks

- SCIP/LSIF toolchains from https://github.com/sourcegraph
- LSIF from Microsoft
- LSIF parser from GitLab
- IDE support from JetBrains

# License

[Apache 2.0](LICENSE)
