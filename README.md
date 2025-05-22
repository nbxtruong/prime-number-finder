# Prime Number Finder

```bash
go run .
```

With parameters
```bash
go run . -workers=4 -max=100000
```

```bash
go test -coverprofile=coverage.out && go tool cover -func=coverage.out
```