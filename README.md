# Tapestry Bindings for Go

[![API tests](https://github.com/Access-Labs-Inc/tapestry-go/actions/workflows/go.yml/badge.svg)](https://github.com/Access-Labs-Inc/tapestry-go/actions/workflows/go.yml)

Bindings for the Tapestry API.

Tapestry documentation: <https://docs.usetapestry.dev/documentation/what-is-tapestry>

Tapestry API reference: <https://tapestry.apidocumentation.com/reference>

## Completness

All current endpoints are implemented.

API tests cover endpoints except for:

- `GET /api/v1/profiles/__ID__/following-who-follow`
- `GET /api/v1/profiles/suggested/__ADDRESS__`
