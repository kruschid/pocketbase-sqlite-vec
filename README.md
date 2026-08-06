# PocketBase + sqlite-vec

A small, multi-architecture Docker image that runs PocketBase with
[`sqlite-vec`](https://github.com/asg017/sqlite-vec) available on every SQLite
connection.

`sqlite-vec` is a SQLite extension, not a replacement storage engine. This
project swaps PocketBase's default pure-Go SQLite driver for the CGO-based
`mattn/go-sqlite3` driver and statically registers `sqlite-vec`. PocketBase's
normal collections, API, dashboard, JavaScript hooks, and migrations continue
to work as usual.

## Quick start

```sh
docker compose up --build
```

Open <http://127.0.0.1:8090/_/> and create the first superuser. Data is stored
in the `pocketbase-data` Docker volume.

To use a published image after pushing this repository to GitHub:

```sh
docker run --rm -p 8090:8090 \
  -v pocketbase-data:/pb/pb_data \
  ghcr.io/kruschid/pocketbase-sqlite-vec:main
```

The container uses fixed mount locations and runs as the non-root user and
group `1000:1000`:

- `/pb/pb_data` for database and uploaded data
- `/pb/pb_hooks` for JavaScript hooks
- `/pb/pb_migrations` for JavaScript migrations
- `/pb/pb_public` for static files

Make sure bind-mounted writable directories are accessible to that user.

## Create and query a vector table

Vector virtual tables are separate from PocketBase collections. Create them in
a JavaScript migration or hook, then link their rows to PocketBase records with
an auxiliary text column. A complete migration is in
[`examples/pb_migrations/1750000000_create_vec_items.js`](examples/pb_migrations/1750000000_create_vec_items.js).

```sql
CREATE VIRTUAL TABLE vec_items USING vec0(
  embedding float[3],
  +record_id TEXT
);

INSERT INTO vec_items(rowid, embedding, record_id)
VALUES (1, '[1, 0, 0]', 'a-pocketbase-record-id');

SELECT rowid, record_id, distance
FROM vec_items
WHERE embedding MATCH '[0.9, 0.1, 0]'
ORDER BY distance
LIMIT 10;
```

Copy the example into `pb_migrations/` and uncomment the migration mount in
`compose.yaml` to try it. The `vec_version()` SQL function reports the embedded
extension version.

## Development

Requirements: Go 1.25+, a C compiler, and Docker.

```sh
make test
make build
```

The executable is container-oriented and always stores data under
`/pb/pb_data`; use Docker Compose to run it locally.

The test suite verifies both `vec_version()` and a real nearest-neighbor query.
The Docker build uses the `no_default_driver` build tag so the unused pure-Go
SQLite implementation is not included.

## Publishing to GHCR

The publish workflow builds `linux/amd64` and `linux/arm64` images. On pushes to
`main` it publishes the `main` tag; Git tags such as `v0.1.0` also publish
semantic-version tags.

## Versions

- PocketBase: `0.39.10`
- sqlite-vec Go bindings: `0.1.6`
- go-sqlite3: `1.14.49`

Dependencies are pinned for reproducible builds and Dependabot is configured
to propose updates. Note that both PocketBase and sqlite-vec are pre-v1 and may
introduce breaking changes.
