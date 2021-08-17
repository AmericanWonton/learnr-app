runtest:
	go test -v
runprogram:
	go run fieldvalidation.go main.go mongoCrudOperations.go pagehandler.go server.go session.go emails.go amazon.go
gomod-exp:
	export GO111MODULE=on
gobuild:
	GOOS=linux GOARCH=amd64 go build -o prodlearnrproj
dockerbuild:
	docker build -t prodlearnrproj .
dockertag:
	docker tag prodlearnrproj americanwonton/prodlearnrproj
dockerpush:
	docker push americanwonton/prodlearnrproj
dockerbuildandpush:
	docker build -t prodlearnrproj .
	docker tag prodlearnrproj americanwonton/prodlearnrproj
	docker push americanwonton/prodlearnrproj
dockerrun:
	docker run -it -p 80:8080 prodlearnrproj
dockerrundetached:
	docker run -d -p 80:8080 prodlearnrproj
dockerrunitvolume:
	docker run -it -p 80:8080 -v photo-images:/static/images prodlearnrproj
dockerrundetvolume:
	docker run -d -p 80:8080 -v photo-images:/static/images prodlearnrproj
dockertagimage:
	docker tag prodlearnrproj americanwonton/prodlearnrproj
dockerimagepush:
	docker push americanwonton/prodlearnrproj
dockerallpush:
	docker tag prodlearnrproj americanwonton/prodlearnrproj
	docker push americanwonton/prodlearnrproj
dockerseeshell:
	docker run -it prodlearnrproj sh