runtest:
	go test -v
runprogram:
	go run main.go mongoCrudOperations.go server.go texting.go amazon.go
gomod-exp:
	export GO111MODULE=on
gobuild:
	GOOS=linux GOARCH=amd64 go build -o prodtextproj
dockerbuild:
	docker build -t prodtextproj .
dockertag:
	docker tag prodtextproj americanwonton/prodtextproj
dockerpush:
	docker push americanwonton/prodtextproj
dockerbuildandpush:
	docker build -t prodtextproj .
	docker tag prodtextproj americanwonton/prodtextproj
	docker push americanwonton/prodtextproj
dockerrun:
	docker run -it -p 80:3000 prodtextproj
dockerrundetached:
	docker run -d -p 80:3000 prodtextproj
dockerrunitvolume:
	docker run -it -p 80:3000 -v photo-images:/static/images prodtextproj
dockerrundetvolume:
	docker run -d -p 80:3000 -v photo-images:/static/images prodtextproj
dockertagimage:
	docker tag prodtextproj americanwonton/prodtextproj
dockerimagepush:
	docker push americanwonton/prodtextproj
dockerallpush:
	docker tag prodtextproj americanwonton/prodtextproj
	docker push americanwonton/prodtextproj
dockerseeshell:
	docker run -it prodtextproj sh