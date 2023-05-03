package example

import (
	log "github.com/sirupsen/logrus"
	"github.com/williamfzc/srctx/parser"
)

func apiDesc() {
	yourLsif := "../parser/lsif/testdata/dump.lsif.zip"

	sourceContext, err := parser.FromLsifFile(yourLsif, "..")
	if err != nil {
		panic(err)
	}

	// all files?
	files := sourceContext.Files()
	log.Infof("files in lsif: %d", len(files))

	// search definition in a specific file
	defs, err := sourceContext.DefsByFileName(files[0])
	if err != nil {
		panic(err)
	}
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
}
