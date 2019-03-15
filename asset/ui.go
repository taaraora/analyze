// +build dev

package asset

import (
	"github.com/shurcooL/httpfs/union"
	"net/http"
)


var Assets http.FileSystem = union.New(map[string]http.FileSystem{
	"/ui":      http.Dir("./ui"),
	"/swagger": http.Dir("./swagger"),
})
