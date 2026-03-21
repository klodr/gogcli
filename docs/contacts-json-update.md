# Update Contacts From JSON

`gog contacts update` supports JSON input via `--from-file`, so you can update People API fields without adding new CLI flags.

## Usage

Update from a file:

```bash
gog contacts get people/c123456 --json > contact.json

# Edit contact.json (see notes below)
gog contacts update people/c123456 --from-file contact.json
```

Update from stdin:

```bash
gog contacts get people/c123456 --json | \
  jq '(.contact.urls //= []) | (.contact.urls += [{"value":"https://example.com","type":"profile"}])' | \
  gog contacts update people/c123456 --from-file -
```

## Input Formats

The command accepts:

- Wrapped (from `gog contacts get --json`): `{"contact": { ...person... }}`
- Direct Person object: `{ ...person... }`

## What Can Be Updated

`--from-file` updates only fields that the People API allows via `people.updateContact` `updatePersonFields`.

Practical rule: include only fields you want to change, at the top level of the JSON object (for example `urls`, `biographies`, `names`, `emailAddresses`, `phoneNumbers`, `addresses`, `organizations`, `genders`, ...).

The following fields are also available as dedicated CLI flags on `contacts create` and `contacts update`, which is simpler than `--from-file` for single-field edits:

| Flag | People API field | Notes |
|---|---|---|
| `--given` / `--family` | `names` | |
| `--email` | `emailAddresses` | |
| `--phone` | `phoneNumbers` | |
| `--org` / `--title` | `organizations` | |
| `--address` | `addresses` | |
| `--url` | `urls` | |
| `--note` | `biographies` | |
| `--birthday` | `birthdays` | YYYY-MM-DD |
| `--notes` | `biographies` | plain text |
| `--relation` | `relations` | `type=person` |
| `--custom` | `userDefined` | `key=value` |
| `--gender` | `genders` | e.g. `male`, `female`, `unspecified`, or any custom value |

If the JSON contains unsupported fields (for `updateContact`), gog errors instead of silently ignoring them.

Notes:

- Some fields are “singleton” for contact sources. Don’t include more than one value for `biographies`, `birthdays`, `genders`, or `names`.
- If you update `memberships`, the Person must include contact group memberships or the API will error.

## Clearing Fields

Clearing list fields is supported by including the key with an empty value:

- Use `[]` to clear a list field (example: `"urls": []`)
- Use `null` to clear a list field (example: `"biographies": null`)

## Concurrency (ETags)

To avoid overwriting concurrent contact edits, gog compares the JSON etag with the current contact etag:

- If they mismatch, update fails with an etag error.
- Use `--ignore-etag` to apply your JSON changes to the latest version anyway.
