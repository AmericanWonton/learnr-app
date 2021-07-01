runtest:
	go test -v
runprogram:
	go run server.go main.go learnrCrudOperations.go crudOperations.go
gomod-exp:
	export GO111MODULE=on
gobuild:
	GOOS=linux GOARCH=amd64 go build -o prodcrudproj
dockerbuild:
	docker build -t prodcrudproj .
dockertag:
	docker tag prodcrudproj americanwonton/prodcrudproj
dockerpush:
	docker push americanwonton/prodcrudproj
dockerbuild:
	docker build -t prodcrudproj .
dockerbuildandpush:
	docker build -t prodcrudproj .
	docker tag prodcrudproj americanwonton/prodcrudproj
	docker push americanwonton/prodcrudproj
dockerrun:
	docker run -it -p 80:4000 prodcrudproj
dockerrundetached:
	docker run -d -p 80:4000 prodcrudproj
dockerrunitvolume:
	docker run -it -p 80:4000 -v photo-images:/static/images prodcrudproj
dockerrundetvolume:
	docker run -d -p 80:4000 -v photo-images:/static/images prodcrudproj
dockertagimage:
	docker tag prodcrudproj americanwonton/prodcrudproj
dockerimagepush:
	docker push americanwonton/prodcrudproj
dockerallpush:
	docker tag prodcrudproj americanwonton/prodcrudproj
	docker push americanwonton/prodcrudproj
dockerseeshell:
	docker run -it prodcrudproj sh