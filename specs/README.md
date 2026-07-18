# Local specification mirror

`manifest.json` pins the exact dated XHTML representations used for XSD 1.1
research. Each digest is lowercase SHA-256 of the response body bytes before
parsing or newline conversion. Representation URLs must use the official dated
`https://www.w3.org/TR/YYYY/REC-name-YYYYMMDD/file.html` or corresponding
`NOTE-name` shape without credentials, a port, query, or fragment; mutable W3C
aliases are rejected.

From the repository root, synchronize every document with:

```sh
go run ./cmd/specsync
```

Pass manifest IDs to synchronize a subset. Selection does not change manifest
order:

```sh
go run ./cmd/specsync xsd11-datatypes
```

The command rejects redirects, non-2xx responses, malformed media types, media
types other than `application/xhtml+xml` and `text/html`, and documents over
64 MiB. It never retrieves external DTDs or other resources. It hashes the
exact acquired bytes, then writes verified XHTML to `specs/cache/` and
generated Markdown to `specs/markdown/`. Both directories are tool-owned,
reproducible, and ignored by Git; edits there are discarded by the next
synchronization. Documents are processed in manifest order and processing
stops at the first failure, leaving earlier completed documents in place.

Markdown starts with a generated-code header, source URL, and digest. Source
headings retain their source level, and fragment IDs become explicit HTML
anchors. A repeated exact ID receives a deterministic generated ID and an
adjacent reversible mapping comment; case-distinct IDs remain distinct.
Fragment-only links remain local, and other relative links resolve against the
representation URL. The complete XHTML body, including navigation and
boilerplate, is converted; unsupported presentation-only containers retain
their content without inventing normative meaning.

Named HTML entities are decoded before the XHTML parse; the five XML entities
retain their XML meaning. Literal `pre` content retains decoded text and
indentation, normalizes CRLF and bare CR line endings to LF, and uses a fence
longer than any backtick run in the snippet. Semantic XHTML classes used for
notes, examples, constraints, schema components, and grammar are retained in
adjacent comments. Output is deterministic UTF-8 with LF line endings, one
final newline, and no timestamp.

To update a pin, change its dated representation URL, acquire the response
without redirects, calculate SHA-256 over the unmodified response body, update
the digest, and run the command. Review both the URL/digest change and the
resulting local conversion before relying on it. Do not replace dated URLs with
mutable `/TR/.../` aliases and do not rewrite or summarize normative prose.
