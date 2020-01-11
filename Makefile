verifiers: staticcheck

staticcheck:
	@echo "Running $@ check"
	@GO111MODULE=on ${GOPATH}/bin/staticcheck ./...