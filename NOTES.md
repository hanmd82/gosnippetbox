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
