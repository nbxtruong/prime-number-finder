# Prime Number Finder

```bash
go run .
```

Or with parameters
```bash
go run . -workers=4 -max=100000
```

![Service log](images/log.png)

```bash
go test -coverprofile=coverage.out && go tool cover -func=coverage.out
```

![Unit test](images/test.png)