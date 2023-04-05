module github.com/williamfzc/srctx

go 1.19

require (
	github.com/alecthomas/chroma/v2 v2.7.0
	github.com/bluekeyes/go-gitdiff v0.7.1
	github.com/dominikbraun/graph v0.16.2
	github.com/opensibyl/sibyl2 v0.15.4
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dlclark/regexp2 v1.4.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/smacker/go-tree-sitter v0.0.0-20230113054119-af7e2ef5fed6 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.23.0 // indirect
	golang.org/x/exp v0.0.0-20220929160808-de9c53c655b9 // indirect
	golang.org/x/sys v0.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/opensibyl/sibyl2 => ../sibyl2

replace github.com/sourcegraph/lsif-go => github.com/williamfzc/lsif-go v0.0.0-20230405041046-51041285c704
