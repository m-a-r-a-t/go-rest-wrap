package internal

import (
	"go_http_test/pkg/router"
)

var R router.Router

func init() {
	R = router.NewRouter("/api")

}
