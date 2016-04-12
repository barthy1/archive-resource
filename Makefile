all: assets/in assets/check assets/out

assets:
	mkdir assets

assets/in: assets in/main.go
	GOARCH=ppc64le GOOS=linux go build -o assets/in in/main.go

assets/out: assets out/main.go
	GOARCH=ppc64le GOOS=linux go build -o assets/out out/main.go

assets/check: assets check/main.go
	GOARCH=ppc64le GOOS=linux go build -o assets/check check/main.go
