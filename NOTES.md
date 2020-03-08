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

Requests Are Handled Concurrently
- all incoming HTTP requests are served in their own goroutines.
- need to guard against race conditions when accessing shared resources from handlers.

---

### Chapter 3
Managing Configuration Settings
- Command-line Flags: common and idiomatic way to manage configuration settings
- provide an explicit and documented interface between the application and its operating configuration
- alternatively, store configuration settings in environment variables and access them directly from application by using the `os.Getenv()` function
  ```
  addr := os.Getenv("SNIPPETBOX_ADDR")
  ```

Leveled Logging
- use `log.New` to define custom logger with arguments (1) destination, (2) prefix, (3) flags.
- custom loggers created by `log.New()` are concurrency-safe.
- initialize a new struct `http.Server` to use the custom errorLog logger.
- recommended to log output to standard streams and redirect the output to a file at runtime.

Dependency Injection
- how to make any dependency available to handlers?
- inject dependencies into handlers to make code more explicit, less error-prone and easier to unit test.
- one approach can be to put dependencies into a custom `application` struct, and define handle functions as methods against `application`

---

### Chapter 4

Database-Driven Responses, using PostgreSQL

- Create `snippetbox` database, and `snippets` table with index on `created_at` field

```sql
$ psql -h localhost -p 5432 -d postgres

CREATE DATABASE snippetbox ENCODING UTF8;
psql> \c snippetbox
-- You are now connected to database "snippetbox" as user "mhan".

CREATE TABLE snippets(
  id SERIAL NOT NULL PRIMARY KEY,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  expires_at TIMESTAMP NOT NULL
);
CREATE INDEX idx_snippets_created_at ON snippets(created_at);

psql> \d snippets
--                                         Table "public.snippets"
--    Column   |            Type             | Collation | Nullable |               Default
-- ------------+-----------------------------+-----------+----------+--------------------------------------
--  id         | integer                     |           | not null | nextval('snippets_id_seq'::regclass)
--  title      | character varying(100)      |           | not null |
--  content    | text                        |           | not null |
--  created_at | timestamp without time zone |           | not null |
--  expires_at | timestamp without time zone |           | not null |
-- Indexes:
--     "snippets_pkey" PRIMARY KEY, btree (id)
--     "idx_snippets_created_at" btree (created_at)
```

- Seed example data
```sql
INSERT INTO snippets (title, content, created_at, expires_at) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    now() at time zone 'utc',
    now() at time zone 'utc' + 365 * INTERVAL '1 day'
);

INSERT INTO snippets (title, content, created_at, expires_at) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    now() at time zone 'utc',
    now() at time zone 'utc' + 365 * INTERVAL '1 day'
);

INSERT INTO snippets (title, content, created_at, expires_at) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    now() at time zone 'utc',
    now() at time zone 'utc' + 7 * INTERVAL '1 day'
);
```

- Create new user `web` and grant access to `snippets` table
```sql
CREATE USER web WITH PASSWORD 'pass';
GRANT SELECT,INSERT,UPDATE ON snippets TO web;
psql> \dp
--                                    Access privileges
--  Schema |      Name       |   Type   | Access privileges | Column privileges | Policies
-- --------+-----------------+----------+-------------------+-------------------+----------
--  public | snippets        | table    | mhan=arwdDxt/mhan+|                   |
--         |                 |          | web=arw/mhan      |                   |
--  public | snippets_id_seq | sequence |                   |                   |
-- (2 rows)
```

- Verify permissions of user `web`
```sql
$ psql -U web -d snippetbox
DROP TABLE snippets;
-- ERROR:  must be owner of table snippets
```
