# srctx: source context

A lib for source context analysis based on [LSIF](https://code.visualstudio.com/blogs/2019/02/19/lsif).

| Name           | Status                                                                                                                                            |
|----------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| Latest Version | ![GitHub release (latest by date)](https://img.shields.io/github/v/release/williamfzc/srctx)                                                      |
| Unit Tests     | [![Go](https://github.com/williamfzc/srctx/actions/workflows/ci.yml/badge.svg)](https://github.com/williamfzc/srctx/actions/workflows/ci.yml)     |
| Code Coverage  | [![codecov](https://codecov.io/github/williamfzc/srctx/branch/main/graph/badge.svg?token=1DuAXh12Ys)](https://codecov.io/github/williamfzc/srctx) |
| Code Style     | [![Go Report Card](https://goreportcard.com/badge/github.com/williamfzc/srctx)](https://goreportcard.com/report/github.com/williamfzc/srctx)      |

## About

This lib processes LSIF file into graphs then you can apply some analysis on them.

This lib originally was designed for monitoring the influence of each commits.

With the raw diff, we can only get something like:

```text
@@ -19,6 +20,7 @@ func AddDiffCmd(app *cli.App) {
        var lsifZip string
        var outputJson string
        var outputCsv string
+       var outputDot string
 
        diffCmd := &cli.Command{
                Name:  "diff",
:

```

If we need to know the impact of these line changes on the entire repository, we can only rely on manual reading of the code to evaluate.

With this lib developers can know exactly what happened in every lines of your code. Such as definition, reference. For example, get all the references happened in a specific line:

```go
vertices, _ := sourceContext.RefsByLine(path, eachLine)
log.Debugf("path %s line %d affected %d vertexes", path, eachLine, len(vertices))
```

Some "dangerous" line changes can be found automatically.

<img width="560" alt="image" src="https://user-images.githubusercontent.com/13421694/230306221-75454e61-7be0-439c-976e-b7f94426c3b9.png">

We hope to utilize the powerful indexing capabilities of LSIF to quantify and evaluate the impact of text changes on the repository, reducing the mental burden on developers.

# Usage

## Usage as Lib

Please see [cmd/srctx/cmd_diff.go](cmd/srctx/cmd_diff.go).

## Usage as Cli

### 1. Generate LSIF file

https://lsif.dev/

LSIF is a standard format for persisted code analyzer output.
Today, several companies are working to support its growth, including Sourcegraph and GitHub/Microsoft.
You can easily find an existed tool for generating LSIF file for your repo.

You will get a `dump.lsif` file after that.

### 2. Run `srctx`

Download our prebuilt binaries from [release page](https://github.com/williamfzc/srctx/releases).

```bash
./srctx diff --lsif dump.lsif --outputCsv output.csv
```

<img width="658" alt="image" src="https://user-images.githubusercontent.com/13421694/230318698-35cdf294-67b0-4eda-8da8-e53602e691ae.png">

You can see every edited lines and their impacts. Currently we provided:

- total referenced
- referenced out of current files
- referenced out of current dir

# Contribution

Issues and PRs are always welcome.

# License

[Apache 2.0](LICENSE)
