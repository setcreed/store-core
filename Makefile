.PHONY: run build image push clean

tag = v0.1
releaseName = store-core
dockerhubUser = setcreed

ALL: run

run: build
	./store-core --configfile ./config.yaml

build:
	go build -o $(releaseName) ./cmd/

image:
	docker build -t $(dockerhubUser)/$(releaseName):$(tag) .

push: image
	docker push $(dockerhubUser)/$(releaseName):$(tag)

clean:
	-rm -f ./$(releaseName)
