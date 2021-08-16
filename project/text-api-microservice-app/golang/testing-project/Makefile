runtest:
	go test -v
runprogram:
	go run main.go mongoCrudOperations.go server.go texting.go amazon.go
gomod-exp:
	export GO111MODULE=on
gobuild:
	GOOS=linux GOARCH=amd64 go build -o testtextproj
dockerbuild:
	docker build -t testtextproj .
dockertag:
	docker tag testtextproj americanwonton/testtextproj
dockerpush:
	docker push americanwonton/testtextproj
dockerbuildandpush:
	docker build -t testtextproj .
	docker tag testtextproj americanwonton/testtextproj
	docker push americanwonton/testtextproj
dockerrun:
	docker run -it -p 80:3000 testtextproj
dockerrundetached:
	docker run -d -p 80:3000 testtextproj
dockerrunitvolume:
	docker run -it -p 80:3000 -v photo-images:/static/images testtextproj
dockerrundetvolume:
	docker run -d -p 80:3000 -v photo-images:/static/images testtextproj
dockertagimage:
	docker tag testtextproj americanwonton/testtextproj
dockerimagepush:
	docker push americanwonton/testtextproj
dockerallpush:
	docker tag testtextproj americanwonton/testtextproj
	docker push americanwonton/testtextproj
dockerseeshell:
	docker run -it testtextproj sh