docker-build:
	docker build -t temporal-website .
	docker image save temporal-website --output temporal-website_docker_image.tar
	gzip -9 temporal-website_docker_image.tar