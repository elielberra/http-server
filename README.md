# http-server

## Introduction
This is an exploratory project to investigate the internal mechanics of HTTP servers. The project originated from a desire to understand how servers function under the hood, moving beyond the abstraction of high-level frameworks. 

Two fundamental rules guided the development:
* The standard `net/http` library was not used for the server implementation.
* AI was not used to generate the source code.

The development process involved researching documentation and analyzing the Go `net` package source code to replicate simplified logic. Instead of relying on high-level libraries, this implementation utilizes direct system calls via `golang.org/x/sys/unix` to manage socket creation, address binding, and connection listening.

## Limitations
This server is a minimalist educational tool and is not intended for production use. It features limited functionality focused on core networking concepts:
* **Supported Methods:** Only `GET` and `POST` methods are accepted.
* **Parameter Parsing:** Query parameters are not parsed.
* **Response Logic:** The server returns a basic response echoing the method, path, and body of the request.
* **Parsing:** The HTTP parser is rudimentary and does not account for complex edge cases or all standard headers.

## Development and Execution
The project uses `air` for live reloading during development. Use the following commands to install the tool and start the server:

```bash
# Install air for development mode
go install [github.com/air-verse/air@latest](https://github.com/air-verse/air@latest)
alias air='$(go env GOPATH)/bin/air'

# Start the server
air
```

### Debugging
A `launch.json` configuration file is included for debugging with VS Code. This allows for setting breakpoints and inspecting the execution flow of the server, specifically for monitoring socket syscalls and request parsing.
