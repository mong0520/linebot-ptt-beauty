APP=linebot-ptt-beauty

build:
	docker build -t mong0520/${APP} .

dev:
	sudo docker-compose up -d app

push:
	@docker tag ${APP} mong0520/${APP}
	@docker push mong0520/${APP}
	@heroku container:push web

release:
	@heroku container:release web