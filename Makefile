generage:
	rm -f gen_*.go
	go run cmd/generate/main.go
test:
	go test -v
