package grpc

import (
	"google.golang.org/grpc/codes"

	"github.com/farnasirim/rex"
)

var codeToErr = map[codes.Code]error{
	codes.Unimplemented: rex.ErrNotImplemented,
}

func errFromCode(code codes.Code) error {
	return codeToErr[code]
}
