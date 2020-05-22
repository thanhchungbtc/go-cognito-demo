export GIN_PORT=4000
export APP_PORT = 3000

dev:
	gin -i main.go

serve:
	go build -o app . && ./app
