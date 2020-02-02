## NOTES

### Chapter 2
A basic web application in `go` has these components:
- handlers (like controllers in MVC frameworks): responsible for executing application logic and writing HTTP response headers and bodies.
- router (i.e. `servemux`): stores a mapping between URL patterns and corresponding handlers.
- web server (built-in).

Fixed Path and Subtree Patterns
- fixed path patterns are only matched when the request URL path exactly matches the fixed path, e.g. `/snippet`, `/snippet/create`.
- subtree path patterns (`.../`) are matched whenever the start of a request URL path matches the subtree path, e.g. `/`, `/static/`.

ServeMux
- `net/http` exposes a default global variable `DefaultServeMux`, which any package can access and register routes, which may be served by malicious handler functions. Security-wise, it may be more prudent to use locally-scoped servemux.
- longer URL patterns take precedence over shorter ones, so URL patterns can be registered in any order.
- request URL paths are automatically sanitised and redirected, e.g. `/../`, `//`.
- does not support routing based on HTTP request method, semantic URLs with variables, regexp-based patterns.

Customizing HTTP Headers
- `w.writeHeader()` can only be called once per response, and cannot be changed after the status code is written.
- first call to `w.Write()` will automatically send `200 OK` status code, unless `w.WriteHeader()` has been called first.
- any changes to the response header map need to be done before calling `w.WriteHeader()` or `w.Write()`.
- use `http.Error` to send non-`200` status code and plain-text response body.
- manipulate the reponse Header Map with `w.Header().Set()`, `.Add()`, `.Delete()`, `.Get()`.
- when sending HTTP response, Go automatically sets three system-generated headers: `Date`, `Content-Length`, `Content-Type`:
  - Go will attempt to sniff the response body with `http.DetectContentType()`.
  - If detection fails, Go will set default response header `Content-Type: application/octet-stream`.
  - `http.DetectContentType()` cannot distinguish JSON from plain text, so for JSON response body, use `w.Header().Set("Content-Type", "application/json")`.

URL Query Strings
- extract query parameters from URL using ` r.URL.Query().Get()` with the desired key - returns as type `string`
- use `strconv.Atoi()` to cast `string` to `integer`
- `http.ResponseWriter` implements the `io.Writer` interface
  ```
  func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
  ```

HTML Templating and Composition
- package `html/template` provides functions for safely parsing and rendering HTML templates.
- use action `{{define "<template-name>"}}...{{end}}` to define distinct named templates.
- use action `{{template "<template-name>" .}}` to invoke named templates, passing in the current context.
- benefits of template composition with layouts and partials:
  - cleanly define the page-specific content in individual files.
  - control which `layout` template the page uses.
  - use `partials` to share and reuse code in different pages or layouts.

Serving Static Files
- use `http.FileServer` handler to serve files over HTTP from a specific directory.
- use `http.StripPrefix()` to strip leading characters from URL path before passing to `http.FileServer`.
