# SGroups API

Protocol Buffers definitions and generated code for the **SGroups** API.

Generated artifacts are written to `pkg/api/` (Go stubs and OpenAPI/Swagger JSON).

You do **not** need to install `protoc` or `protoc-gen-*` plugins manually.

## Quick Start

Generate Go code + Swagger/OpenAPI JSON:

`make generate-api`

Run tests:

`make test`

## Swagger/OpenAPI

Swagger JSON files are generated into `pkg/api/**.swagger.json` and embedded by the Go package under `pkg/`.
