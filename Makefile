build:
	go build -v -o altv cmd/altv/*.go
completion:
	go run cmd/altv/*.go completion bash > /etc/bash_completion.d/altv
