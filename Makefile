run: build
	@./bin/app

build:
	@go build -o bin/app .

rundev:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "assets" \
	--build.include_ext "js,css"


gen-css:
	tailwindcss -i static/style-input.css -o static/style.css --watch

fmt-templ:
	@templ fmt views/*.templ

fmt-css:
	@tailwindcss -i static/style-input.css -o static/style.css -c tailwind.config.js

watch-templ:
	@templ generate -watch -path views/

watch-css:
	@tailwindcss -i static/style-input.css -o static/style.css --watch -c tailwind.config.js

fmt-go:
	@go fmt ./...