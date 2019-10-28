all: fetch_strongswan docker_build

fetch_strongswan:
	./fetch-strongswan

docker_build:
	docker build --build-arg http_proxy=${http_proxy} --build-arg https_proxy=${https_proxy} --build-arg no_proxy=${no_proxy} -t nuage/strongswan:${version} .
	docker save nuage/strongswan:${version} > nuage-strongswan-docker-${version}.tar
	docker rmi nuage/strongswan:${version} 
