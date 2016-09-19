package cherry

import (
	"net/http"
)

type Business interface {
	Init(resRoot string)
	Handle(w http.ResponseWriter, r *http.Request)
}
