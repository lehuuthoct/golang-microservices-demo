build-leader-service:
	GOARCH=amd64 GOOS=linux go build -o ./leader/leader-service-linux-amd64 ./leader
	GOARCH=amd64 GOOS=linux go build -o ./system-checker/system-checker-linux-amd64 ./system-checker
	cp ./system-checker/system-checker-linux-amd64 ./leader/

# build leader docker service
build-docker-leader:
	make build-leader-service

	#	repository name must be lowercase
	docker build -t lht/leader-service leader

	make show-docker-images

# run & remove container when it exists
run-docker-leader:
	docker run --rm lht/leader-service

test-all:
	go test ./...

# prerequisites:
# - install virtualbox (Mac):
# -- brew cask install virtualbox

# - create second machine lhtvm1:
# -- docker-machine create --driver virtualbox lhtvm1

# -- machine 1 (swarm manager: default): 192.168.64.3:2376
# -- machine 2 (worker: lhtvm1): 192.168.99.100:2376

# 	curl localhost:9001/leader/1
#	curl localhost:9001/health
deploy-service-leader:
	docker service rm lht-service; docker service create --name=lht-service --replicas=1 --network=lht_network -p=9001:9001 lht/leader-service
	docker service ls

build-deploy-service-leader:
	make build-docker-leader
	make deploy-service-leader

scale-leader-service:
	docker service scale lht-service=3

# deploy a spring boot service from docker hub
deploy-quote-service:
	docker service create --name=quotes-service --replicas=1 --network=lht_network eriklupander/quotes-service

show-docker-images:
	docker image ls
