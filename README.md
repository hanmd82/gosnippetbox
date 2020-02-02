# SnippetBox
A simple web application written in golang, for pasting and sharing snippets of text.

Based on https://lets-go.alexedwards.net/.

---

## Directory Structure

- `cmd`: application-specific code for executable applications.
- `pkg`: non-application-specific code, such as validation helpers and database models.
- `ui`: user-interface assets, such as HTML templates and static files, e.g CSS, images.
  - `ui/html`: templates with filename convention `<name>.<role>.tmpl`, where `role` can be `page`, `partial` or `layout`.

Reference: https://peter.bourgon.org/go-best-practices-2016/#repository-structure
