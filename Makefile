test:
	go test github.com/kevinjqiu/phantomail/pkg/...

install:
	- go get -v github.com/Masterminds/glide
	- glide install
