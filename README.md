![tests](https://github.com/nil-nil/ticket/actions/workflows/go-test.yml/badge.svg)

<div align="center">
  <h3 align="center">Ticket</h3>
  <p align="center">
    A Ticket system.
  </p>
</div

# Development

## Running the dev environment locally

Use the vscode tasks.

You need [Air](https://github.com/cosmtrek/air) to run the `api` task.

You also need Go to build the API server.

Air will rebuild the project every time there is a change.

## API

The API is defined in OpenAPI files in `./apispec`. In there you will also find the `openapi-codegen.conf.yaml` configuration file for [Deepmap's OpenAPI Code Generator](https://github.com/deepmap/oapi-codegen). This config instructs it to use Echo as a webserver, and to use strict mode (generating RPC style handlers to reduce boilerplate), and sets the output file,

You can trigger the codegen using the `openapi` task.
