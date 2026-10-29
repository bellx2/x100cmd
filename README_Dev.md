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

OTPサーバーのURLはソースに含めず、ビルド時に埋め込みます。未指定の場合は`--url`を指定しない限りOTP認証（`otp`、`bank`）が使えません。

```bash
go build -ldflags "-X github.com/bellx2/x100cmd/cmd.otpURL=https://example.com/v1/otp/:otp" -o x100cmd main.go
```

リリース（GoReleaser）では環境変数`OTP_URL`の値を埋め込みます。GitHub Actionsではリポジトリの Secrets に`OTP_URL`を登録してください。
