.PHONY: all go-iam release-version

all: clean dist

version := $(shell cat VERSION)

#util/version.go:
release-version:
	git rev-parse HEAD|awk 'BEGIN {print "package util"} {print "const BuildGitVersion=\""$$0"\""} END{}' > util/version.go
	date +'%Y%m%d%H'| awk 'BEGIN{} {print "const BuildGitDate=\""$$0"\""} END{}' >> util/version.go

dist: go-iam
	mkdir -p build/go-iam-$(value version)
	cp run.sh build/go-iam-$(value version)/bin
	cd build && tar cvzf go-iam.tar.gz go-iam-$(value version)
	# rm -r build/go-iam-$(value version)

go-iam: release-version
	mkdir -p build/go-iam-$(value version)/bin
	go build ${BUILD_FLAGS} -o build/go-iam-$(value version)/bin/go-iam

clean:
	rm -rf build

