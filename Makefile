
test:
	ginkgo -r -v -p --cover
	
lint:
	golangci-lint -v run
