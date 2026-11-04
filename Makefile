include .env.local
MIGRATIONS_DIR = ./migrations

##Command: Description: Usage
##--------------:--------------:--------------

## help: print this help message
help:
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

live/templ:
	go tool templ generate --watch --proxy="http://localhost:6969" --open-browser=false -path="./views" -cmd="echo Initial generation done"

live/server:
	go run github.com/air-verse/air@v1.63.0 \
	--build.cmd "go build -tags dev -o tmp/bin/main" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

live/tailwind:
	npx --yes @tailwindcss/cli -i ./assets/css/input.css -o ./static/style.css --minify --watch=always

live/esbuild:
	npx esbuild --bundle assets/js/index.ts --outdir=static --global-name=bundle --watch=forever

live/sync_assets:
	go run github.com/air-verse/air@v1.63.0 \
	--build.cmd "go tool templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "static" \
	--build.include_ext "js,css"

## live: run the server on watch mode
live: 
	make -j5 live/templ live/server live/tailwind live/esbuild live/sync_assets

## migrate_create: Creates new migration file with the current timestamp. : make migrate_create name=[name]
migrate_create:
	@go tool goose postgres $(DSN) -dir ${MIGRATIONS_DIR} create $(name) sql

## migrate_up: Migrate the DB to the most recent version available.
migrate_up:
	@go tool goose postgres $(DSN) -dir ${MIGRATIONS_DIR} up

## migrate_down: Roll back the version by 1.
migrate_down:
	@go tool goose postgres $(DSN) -dir ${MIGRATIONS_DIR} down

## migrate_status: Dump the migration status for the current DB.
migrate_status:
	@go tool goose postgres $(DSN) -dir ${MIGRATIONS_DIR} status
