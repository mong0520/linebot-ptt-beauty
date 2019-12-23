APP=linebot-ptt-beauty
NAMESPACE := mong0520

build:
	docker build -t ${NAMESPACE}/${APP} .

dev:
	sudo docker-compose up -d

push:
	@docker tag ${APP} ${NAMESPACE}/${APP}
	@docker push ${NAMESPACE}/${APP}
	@heroku container:push web

release:
	@heroku container:release web