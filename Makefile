GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /amazonian/)
default: build

workdir:
	mkdir -p workdir

build: workdir/amazonian

workdir/amazonian: $(GOFILES)
	go build -o workdir/amazonian .

dependencies: 
	@go get github.com/aws/aws-sdk-go/service/cloudformation
	@go get github.com/aws/aws-sdk-go/service/elbv2
	# @go get github.com/kevinburke/go-bindata
	aws s3 cp s3://ecs.bucket.template/ecs/ecs.yml ias/cloudformation 
	aws s3 cp s3://ecs.bucket.template/ecstenant/containertemplate.yml ias/cloudformation 
	aws s3 cp s3://ecs.bucket.template/vpc/vpc.yml ias/cloudformation 

	# ./go-bindata -o assets/myfile.go ias/...
bindata:
	./go-bindata -o assets/bindata.go ias/...
test: test-all

test-all:
	@go test -v ./...

test-min:
	@go test ./...

release:
	aws s3 cp workdir/amazonian s3://amazonian.package.release/latest/amazonian
	aws s3 cp workdir/amazonian s3://amazonian.package.release/$(VERSION)/amazonian
