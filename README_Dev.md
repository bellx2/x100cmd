# Development Memo

## Setup

Use `mise` to setup the development environment.
For more information, visit https://mise.jdx.dev/.

```bash
curl https://mise.run | sh
```

## Run

```bash
go mod tidy
go run main.go
```

## Build

```bash
go build -o x100cmd main.go
```
