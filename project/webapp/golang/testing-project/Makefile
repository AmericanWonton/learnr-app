runtest:
	go test -v
runprogram:
	go run fieldvalidation.go main.go mongoCrudOperations.go pagehandler.go server.go session.go emails.go amazon.go
gomod-exp:
	export GO111MODULE=on
gobuild:
	GOOS=linux GOARCH=amd64 go build -o testlearnrproj
dockerbuild:
	docker build -t testlearnrproj .
dockertag:
	docker tag testlearnrproj americanwonton/testlearnrproj
dockerpush:
	docker push americanwonton/testlearnrproj
dockerbuildandpush:
	docker build -t testlearnrproj .
	docker tag testlearnrproj americanwonton/testlearnrproj
	docker push americanwonton/testlearnrproj
dockerrun:
	docker run -it -p 80:8080 testlearnrproj
dockerrundetached:
	docker run -d -p 80:8080 testlearnrproj
dockerrunitvolume:
	docker run -it -p 80:8080 -v photo-images:/static/images testlearnrproj
dockerrundetvolume:
	docker run -d -p 80:8080 -v photo-images:/static/images testlearnrproj
dockertagimage:
	docker tag testlearnrproj americanwonton/testlearnrproj
dockerimagepush:
	docker push americanwonton/testlearnrproj
dockerallpush:
	docker tag testlearnrproj americanwonton/testlearnrproj
	docker push americanwonton/testlearnrproj
dockerseeshell:
	docker run -it testlearnrproj sh