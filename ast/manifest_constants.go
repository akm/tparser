package ast

import (
	"strings"

	"github.com/akm/tparser/ext"
)

func IsManifestConstant(w string) bool {
	return manifestConstants.Include(strings.ToUpper(w))
}

var manifestConstants = ext.Strings{
	"FALSE",
	"TRUE",
}.Set()
