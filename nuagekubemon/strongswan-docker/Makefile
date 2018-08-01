all: fetch_strongswan docker_build

fetch_strongswan:
	./fetch-strongswan

docker_build:
	docker build -t nuage/strongswan:${version} .
	docker save nuage/strongswan:${version} > nuage-strongswan-docker-${version}.tar
	docker rmi nuage/strongswan:${version} 
