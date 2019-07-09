test:
	go test github.com/kevinjqiu/phantomail/pkg/...

deps:
	go get -v github.com/Masterminds/glide
	glide install

build: deps
	mkdir -p build
	patch vendor/github.com/caddyserver/caddy/caddy/caddymain/run.go < patch.diff
	cd vendor/github.com/caddyserver/caddy/caddy/ && go build
	cp vendor/github.com/caddyserver/caddy/caddy/caddy build/
