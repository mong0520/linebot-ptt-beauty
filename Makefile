APP=linebot-ptt-beauty
NAMESPACE := mong0520

build:
	docker build -t ${NAMESPACE}/${APP} .

dev:
	docker-compose up -d

down:
	docker-compose down

push:
	@docker tag ${APP} ${NAMESPACE}/${APP}
	@docker push ${NAMESPACE}/${APP}
	@heroku container:push web

release:
	@heroku container:release web