package web

import (
	"embed"
	"io/fs"
)

//go:embed static
var embeddedStatic embed.FS

var AssetFS fs.FS

func init() {
	var err error
	AssetFS, err = fs.Sub(embeddedStatic, "static")
	if err != nil {
		panic(err)
	}
}
