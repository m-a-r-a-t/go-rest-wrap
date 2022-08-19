package internal

import (
	"github.com/m-a-r-a-t/go-rest-wrap/pkg/router"
)

var R router.Router

func init() {
	R = router.NewRouter("/api")

}
