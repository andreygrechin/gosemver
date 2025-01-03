VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
MOD_PATH := $(shell go list -m)
MODULE_NAME_SHORT := gosemver

build:
	go build -ldflags "-s -w -X $(MOD_PATH)/internal/config.Version=$(VERSION) -X $(MOD_PATH)/internal/config.BuildTime=$(BUILDTIME) -X $(MOD_PATH)/internal/config.Commit=$(COMMIT)" -o bin/$(MODULE_NAME_SHORT)
	bin/$(MODULE_NAME_SHORT) version

lint:
	go vet ./...
	staticcheck ./...
	golangci-lint run --show-stats

vuln:
	gosec ./...
	govulncheck ./...

check_clean:
	@if [ -n "$(shell git status --porcelain)" ]; then \
		echo "Error: Dirty working tree. Commit or stash changes before proceeding."; \
		exit 1; \
	fi

release: check_clean lint vuln
	goreleaser release --clean

define bump_version
	$(eval NEW_VERSION := $(shell semver bump $(1) $(VERSION)))
	@echo "Old version $(VERSION)"
	@echo "Bumped to version $(NEW_VERSION)"
	@git tag -a "v$(NEW_VERSION)" -m "v$(NEW_VERSION)"
	@git push origin "v$(NEW_VERSION)"
endef

bump_patch: check_clean lint vuln
	$(call bump_version,patch)

bump_minor: check_clean lint vuln
	$(call bump_version,minor)

bump_major: check_clean lint vuln
	$(call bump_version,major)

.PHONY: build lint vuln release check_clean bump_patch bump_minor bump_major
