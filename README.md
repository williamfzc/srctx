# srctx: source context

A library for extracting and analyzing definition/reference graphs from your codebase. Powered by tree-sitter and LSIF/SCIP.

| Name           | Status                                                                                                                                            |
|----------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| Latest Version | ![GitHub release (latest by date)](https://img.shields.io/github/v/release/williamfzc/srctx)                                                      |
| Unit Tests     | [![Go](https://github.com/williamfzc/srctx/actions/workflows/ci.yml/badge.svg)](https://github.com/williamfzc/srctx/actions/workflows/ci.yml)     |
| Code Coverage  | [![codecov](https://codecov.io/github/williamfzc/srctx/branch/main/graph/badge.svg?token=1DuAXh12Ys)](https://codecov.io/github/williamfzc/srctx) |
| Code Style     | [![Go Report Card](https://goreportcard.com/badge/github.com/williamfzc/srctx)](https://goreportcard.com/report/github.com/williamfzc/srctx)      |

## About this tool

This library processes your code into precise function-level graphs, seamlessly integrated with Git, and then you can apply some analysis to them.

<img width="1389" alt="image" src="https://github.com/williamfzc/srctx/assets/13421694/e48a51c7-e95b-4da7-994f-f5fe1f461477">

With this lib developers can know exactly what happened in every lines of your code. Such as definition, reference. And understand the actual impacts of your git commits.

Some "dangerous" line changes can be found automatically.

<img width="843" alt="image" src="https://github.com/williamfzc/srctx/assets/13421694/e6e48e67-35b1-4c52-aa6a-99b1ba2f02db">

You can see a dangerous change in file `cmd/srctx/diff/cmd.go#L29-#143`, .

We hope to utilize the powerful indexing capabilities of LSIF to quantify and evaluate the impact of text changes on the
repository, reducing the mental burden on developers.

# Usage

## Quick Start

We provide a one-click script for quickly deploying srctx anywhere. Common parameters include:

- SRCTX_LANG: Required, specifies the language, such as GOLANG/JAVA/KOTLIN.
- SRCTX_BUILD_CMD: Optional, specifies the compilation command.

### For Golang

```bash
curl https://raw.githubusercontent.com/williamfzc/srctx/main/scripts/quickstart.sh \
| SRCTX_LANG=GOLANG bash
```

### For Java/Kotlin

As there is no unique compilation toolchain for Java (it could be Maven or Gradle, for example), 
so at the most time, you also need to specify the compilation command to obtain the invocation information.

You should replace the `SRCTX_BUILD_CMD` with your own one.

Java:

```bash
curl https://raw.githubusercontent.com/williamfzc/srctx/main/scripts/quickstart.sh \
| SRCTX_LANG=JAVA SRCTX_BUILD_CMD="clean package -DskipTests" bash
```

Kotlin:

Change the `SRCTX_LANG=JAVA` to `SRCTX_LANG=KOTLIN`.

## In Production

In proudction, it is generally recommended to separate the indexing process from the analysis process, rather than using a one-click script to complete the entire process. This can make the entire process easier to maintain.

### 1. Generate LSIF file

Tools can be found in https://lsif.dev/ .

You will get a `dump.lsif` file after that.

### 2. Run `srctx`

Download our prebuilt binaries from [release page](https://github.com/williamfzc/srctx/releases).

For example, diff from `HEAD~1` to `HEAD`:

```bash
./srctx diff \
  --before HEAD~1 \
  --after HEAD \
  --lsif dump.lsif \
  --outputCsv output.csv \
  --outputDot output.dot \
  --outputHtml output.html
```

See details with `./srctx diff --help`.

### Prefer a real world sample?

[Our CI](https://github.com/williamfzc/srctx/blob/main/.github/workflows/ci.yml) is a good start.

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
