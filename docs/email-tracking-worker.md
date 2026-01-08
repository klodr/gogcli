---
summary: "Tracking Worker internals (routes, schema, keys)"
read_when:
  - Changing Worker endpoints or schema
  - Debugging tracking id collisions / open queries
---

# Email tracking worker

Location:
- Worker source: `internal/tracking/worker/src/`
- Schema: `internal/tracking/worker/schema.sql`

## Bindings / config

Expected bindings:
- D1 database binding: `DB`
- Secrets: `TRACKING_KEY`, `ADMIN_KEY`

`wrangler.toml` is the local template; deployments set the real D1 database id.

## Routes (high-level)

- Pixel:
  - `GET /t/<tracking_id>.png`
  - Validates/decrypts `tracking_id`, stores an open row, returns a transparent PNG.

- Admin:
  - `GET /opens?to=<email>&since=<...>`
  - `GET /opens/<tracking_id>`
  - Auth: `Authorization: Bearer <ADMIN_KEY>` (or equivalent, per implementation).

## Schema notes

- Primary key uses `tracking_id` (the encrypted blob) to avoid collisions.
- `opened_at` stored as an ISO string for consistent ordering/comparison.

## Local dev

```sh
cd internal/tracking/worker
pnpm install
pnpm dev
```

## Tests

```sh
cd internal/tracking/worker
pnpm lint
pnpm build
pnpm test
```

