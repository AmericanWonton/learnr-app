runtest:
	go test -v
runprogram:
	go run server.go main.go learnrCrudOperations.go crudOperations.go
gomod-exp:
	export GO111MODULE=on
gobuild:
	GOOS=linux GOARCH=amd64 go build -o testcrudproj
dockerbuild:
	docker build -t testcrudproj .
dockertag:
	docker tag testcrudproj americanwonton/testcrudproj
dockerpush:
	docker push americanwonton/testcrudproj
dockerbuild:
	docker build -t testcrudproj .
dockerbuildandpush:
	docker build -t testcrudproj .
	docker tag testcrudproj americanwonton/testcrudproj
	docker push americanwonton/testcrudproj
dockerrun:
	docker run -it -p 80:4000 testcrudproj
dockerrundetached:
	docker run -d -p 80:4000 testcrudproj
dockerrunitvolume:
	docker run -it -p 80:4000 -v photo-images:/static/images testcrudproj
dockerrundetvolume:
	docker run -d -p 80:4000 -v photo-images:/static/images testcrudproj
dockertagimage:
	docker tag testcrudproj americanwonton/testcrudproj
dockerimagepush:
	docker push americanwonton/testcrudproj
dockerallpush:
	docker tag testcrudproj americanwonton/testcrudproj
	docker push americanwonton/testcrudproj
dockerseeshell:
	docker run -it testcrudproj sh