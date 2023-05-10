# srctx: source context

A library for converting your codebase into graph data structures. Powered by tree-sitter and LSIF.

| Name           | Status                                                                                                                                            |
|----------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| Latest Version | ![GitHub release (latest by date)](https://img.shields.io/github/v/release/williamfzc/srctx)                                                      |
| Unit Tests     | [![Go](https://github.com/williamfzc/srctx/actions/workflows/ci.yml/badge.svg)](https://github.com/williamfzc/srctx/actions/workflows/ci.yml)     |
| Code Coverage  | [![codecov](https://codecov.io/github/williamfzc/srctx/branch/main/graph/badge.svg?token=1DuAXh12Ys)](https://codecov.io/github/williamfzc/srctx) |
| Code Style     | [![Go Report Card](https://goreportcard.com/badge/github.com/williamfzc/srctx)](https://goreportcard.com/report/github.com/williamfzc/srctx)      |

## About this tool

This lib processes your code into function level graphs then you can apply some analysis on them.

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

If we need to know the impact of these line changes on the entire repository, we can only rely on manual reading of the
code to evaluate.

With this lib developers can know exactly what happened in every lines of your code. Such as definition, reference.

```bash
./srctx diff --outputDot diff.dot
```

![display](https://user-images.githubusercontent.com/13421694/236665125-4968558b-8601-43d0-9618-97e146f93749.svg)

Some "dangerous" line changes can be found automatically.

We hope to utilize the powerful indexing capabilities of LSIF to quantify and evaluate the impact of text changes on the
repository, reducing the mental burden on developers.

# Usage

## Out-of-box production (Recommendation)

Because LSIF files require dev env heavily, it's really hard to provide a universal solution in a single binary file for
all the repos.

<img width="697" alt="image" src="https://user-images.githubusercontent.com/13421694/236666915-5d403e4a-9cc1-4364-afbe-363cf82e5e49.png">

We are currently working on [diffctx](https://github.com/williamfzc/diffctx), which will provide a GitHub Actions plugin
that allows users to use it directly in a Pull Request.

## Usage as Cli (Recommendation)

A short golang example:

```bash
# indexing
curl -L  https://github.com/sourcegraph/lsif-go/releases/download/v1.9.3/src_linux_amd64 -o /usr/local/bin/lsif-go
chmod +x /usr/local/bin/lsif-go
lsif-go -v

# diff
wget https://github.com/williamfzc/srctx/releases/download/v0.4.2/srctx-linux-amd64
chmod +x srctx-linux-amd64
./srctx-linux-amd64 diff --before HEAD~1 --after HEAD --lsif dump.lsif --outputCsv output.csv --outputDot output.dot
```

It will produce:

- JSON for processing by other programs
- CSV for reading
- graphviz DOT file

### 1. Generate LSIF file

Some official tools we used in `diffctx`:

| Language    | Link                                                                                                | Command                                                                  |
|-------------|-----------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------|
| Golang      | [lsif-go](https://github.com/sourcegraph/lsif-go)                                                   | `lsif-go -v`                                                             |
| Python      | [lsif-py](https://github.com/sourcegraph/lsif-py)                                                   | `lsif-py .`                                                              |
| Java/Kotlin | [scip-java](https://github.com/sourcegraph/scip-java) / [scip](https://github.com/sourcegraph/scip) | see [index.py](https://github.com/williamfzc/diffctx/blob/main/index.py) |
| JavaScript  | [lsif-node](https://github.com/microsoft/lsif-node)                                                 | see [website](https://github.com/microsoft/lsif-node)                    |

You will get a `dump.lsif` file after that.

### 2. Run `srctx`

Download our prebuilt binaries from [release page](https://github.com/williamfzc/srctx/releases).

For example, diff from `HEAD~1` to `HEAD`:

```bash
./srctx diff --before HEAD~1 --after HEAD --lsif dump.lsif --outputCsv output.csv --outputDot output.dot
```

See details with `./srctx diff --help`.

## Usage as Lib

### API

Our built-in diff implementation is a good example. [cmd/srctx/diff/cmd.go](cmd/srctx/diff/cmd.go)

### Low level API

Low level API allows developers consuming LSIF file directly.

```go
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

# Roadmap

- Simpler installation
- Full support LSIF
- Better API

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
