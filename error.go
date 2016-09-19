package cherry

import (
	"fmt"
	"html/template"
	"net/http"
)

var (
	httpErr map[int]http.HandlerFunc
)

func AddHttpErr(code int, f http.HandlerFunc) {
	httpErr[code] = f
}

var errTpl = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<title>{{.Title}}</title>
</head>
<body>
	<h1>{{.Title}}</h1>
	<p>{{.Content}}</p>
	<hr />
	<p>Powered by <a href="https://github.com/wangybnet/cherry" target="_blank">cherry</a>/{{.CherryVersion}}</p>
</body>
</html>
`

func HttpErr(w http.ResponseWriter, r *http.Request, code int) {
	if httpErr[code] != nil {
		httpErr[code](w, r)
	} else {
		t, _ := template.New("errTpl").Parse(errTpl)
		data := make(map[string]string)
		data["Title"] = fmt.Sprintf("HTTP Status %d", code)
		data["Content"] = ""
		data["CherryVersion"] = VERSION
		w.WriteHeader(code)
		t.Execute(w, data)
	}
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("errTpl").Parse(errTpl)
	data := make(map[string]string)
	data["Title"] = "Bad Request"
	data["Content"] = "The server could not understand this request."
	data["CherryVersion"] = VERSION
	w.WriteHeader(400)
	t.Execute(w, data)
}

func forbidden(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("errTpl").Parse(errTpl)
	data := make(map[string]string)
	data["Title"] = "Forbidden"
	data["Content"] = fmt.Sprintf("You don't have permission to access %s on this server.", r.URL.Path)
	data["CherryVersion"] = VERSION
	w.WriteHeader(403)
	t.Execute(w, data)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("errTpl").Parse(errTpl)
	data := make(map[string]string)
	data["Title"] = "Not Found"
	data["Content"] = fmt.Sprintf("The requested URL %s was not found on this server.", r.URL.Path)
	data["CherryVersion"] = VERSION
	w.WriteHeader(404)
	t.Execute(w, data)
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("errTpl").Parse(errTpl)
	data := make(map[string]string)
	data["Title"] = "Internal Server Error"
	data["Content"] = "The server encountered an internal error and was unable to complete your request."
	data["CherryVersion"] = VERSION
	w.WriteHeader(500)
	t.Execute(w, data)
}

func serviceUnavailable(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("errTpl").Parse(errTpl)
	data := make(map[string]string)
	data["Title"] = "Service Unavailable"
	data["Content"] = "The server is temporarily unable to service your request. Please try again later."
	data["CherryVersion"] = VERSION
	w.WriteHeader(503)
	t.Execute(w, data)
}

func init() {
	httpErr = make(map[int]http.HandlerFunc)
	httpErr[400] = badRequest
	httpErr[403] = forbidden
	httpErr[404] = notFound
	httpErr[500] = internalServerError
	httpErr[503] = serviceUnavailable
}
