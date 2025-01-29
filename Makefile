run:
	go run ./cmd/api

wire:
	 cd pkg/di && wire
	 
swag :
	swag init -g cmd/api/main.go -o ./cmd/api/docs