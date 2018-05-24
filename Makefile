.PHONY:	up build destroy assets

SERVICE="inspectr-backend"
CONFIG="./configs/config.dev.yml"

up:
	docker-compose run --rm ${SERVICE} go run main.go --config ${CONFIG} migrate up
	docker-compose up -d redis postgres
	docker-compose up ${SERVICE}

build:
	docker-compose build --pull ${SERVICE}

assets:
	docker-compose run --rm ${SERVICE} go-bindata -pkg assets -o assets/assets.go \
		plugins/api/schema.graphql \
		plugins/api/static/

destroy:
	docker-compose stop
	docker-compose rm -f
