<h1 align="center">
  <img src="https://github.com/williamfzc/srctx/assets/13421694/e99b5cf6-07d7-49fb-a70a-862deab83e49" width="400" height="300">
</h1>

<h3 align="center">srctx: source context</h3>
<p align="center">
    <em>A library for extracting and analyzing definition/reference graphs from your codebase. Powered by tree-sitter and LSIF/SCIP.</em>
</p>

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

In addition, as a library, it also provides convenient APIs for secondary development, allowing you to freely access the content in the graph.

```golang
src := filepath.Dir(filepath.Dir(curFile))
lsif := "../dump.lsif"
lang := core.LangGo

funcGraph, _ := function.CreateFuncGraphFromDirWithLSIF(src, lsif, lang)

functions := funcGraph.GetFunctionsByFile("cmd/srctx/main.go")
for _, each := range functions {
    // about this function
    log.Infof("func: %v", each.Id())
    log.Infof("decl location: %v", each.FuncPos.Repr())
    log.Infof("func name: %v", each.Name)

    // context of this function
    outVs := funcGraph.DirectReferencedIds(each)
    log.Infof("this function reach %v other functions", len(outVs))
    for _, eachOutV := range outVs {
        outV, _ := funcGraph.GetById(eachOutV)
        log.Infof("%v directly reached by %v", each.Name, outV.Name)
    }
}
```

> Currently, srctx is still in an active development phase. 
> If you're interested in its iteration direction and vision, you can check out [our roadmap page](https://github.com/williamfzc/srctx/issues/31).

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

This API allows developers accessing the data of FuncGraph.

- Start up example: [example/api_test.go](example/api_test.go)
- Real world example: [cmd/srctx/diff/cmd.go](cmd/srctx/diff/cmd.go)

### Low level API

Low level API allows developers consuming LSIF file directly.

See [example/api_base_test.go](example/api_base_test.go) for details.

# Correctness / Accuracy

<img width="1159" alt="image" src="https://github.com/williamfzc/srctx/assets/13421694/6cfa72c2-787a-4ae6-8cef-e77c1985d307">

We wanted it to provide detection capabilities as accurate as an IDE.

# Roadmap

See [Roadmap Issue](https://github.com/williamfzc/srctx/issues/31).

# Contribution

[Issues and PRs](https://github.com/williamfzc/srctx/issues) are always welcome.

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
