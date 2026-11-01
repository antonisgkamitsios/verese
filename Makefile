live/templ:
	go tool templ generate --watch --proxy="http://localhost:6969" --open-browser=false -path="./views" -cmd="echo Initial generation done"
live/server:
	go run github.com/air-verse/air@v1.63.0 \
	--build.cmd "go build -o tmp/bin/main" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true
live/tailwind:
	npx --yes @tailwindcss/cli -i ./assets/css/input.css -o ./static/output.css --minify --watch=always
live/sync_assets:
	go run github.com/air-verse/air@v1.63.0 \
	--build.cmd "go tool templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "static" \
	--build.include_ext "js,css"
live: 
	make -j4 live/templ live/server live/tailwind live/sync_assets