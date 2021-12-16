# GOQU CRUD Generator

[![Go](https://github.com/funvit/goqu-crud-gen/actions/workflows/go.yml/badge.svg)](https://github.com/funvit/goqu-crud-gen/actions/workflows/go.yml)

> Work in progress!

Generates a basic repository by db model definition.

Expected to be used by [goqu](https://github.com/doug-martin/goqu)
and [sqlx](https://github.com/jmoiron/sqlx) users.

# What is generated?

> TODO: list and describe generated structs and methods.

Generator creates a repository definition file containing the main CRUD methods:

- `Create`
- `Get`
- `Update`
- `Delete`, `DeleteMany`

... and some special methods:

- `WithTran` currently must be used for wrapping each repo method call.
- `each`, `iter` for user future use in custom repo methods.

# Install

```bash
$ go get github.com/funvit/goqu-crud-gen/cmd/...
```

# Usage

Define a db-model.

Rules:

- model must have one field marked as primary key via adding option  `primary`
  for tag `db`
    - supported field types:
        - standard (`int64`, `string`...)
        - any other which implements `Scanner` and `Valuer` interfaces (ex:
          github.com/google/uuid)
- if model primary key field value is database-side generated `int64` -
  use `auto` option for tag `db`
- generated file will be placed near model definition

See examples for mode info.

# Examples

See [./examples](./examples) folder.

# TODO

- [ ] tests
- [ ] tests with mysql in docker
- [ ] string field maxlen rule by annotation?
- [x] ~~new flag allowed customising `WithTran` method
  name.~~ (`-rename-with-tran`)
