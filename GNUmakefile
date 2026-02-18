SWAGGER_URL ?= https://rootly-heroku.s3.amazonaws.com/swagger/v1/swagger.tf.json
TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=hashicorp.com
NAMESPACE=eduW
NAME=rootly
BINARY=terraform-provider-${NAME}
VERSION=0.1.1
OS_ARCH=darwin_amd64

default: testacc

# Run acceptance tests
.PHONY: testacc codegen build release install test docs
build: codegen docs
	go build -o ${BINARY}

docs:
	terraform fmt -recursive examples
	go tool tfplugindocs
	mv ./docs/data-sources/ip_ranges.md ./docs/data-sources/ip_range.md
	rm ./docs/data-sources/*s.md
	mv ./docs/data-sources/ip_range.md ./docs/data-sources/ip_ranges.md
	find ./docs/resources/*.md -type f -exec node tools/clean-docs.js {} \;
	find ./docs/resources/ -type f -name 'workflow_task_*.md' -exec perl -pi -e 's/subcategory:$$/subcategory: Workflow Tasks/g' {} +
	find ./docs/resources/ -type f -name 'workflow_*.md' -exec perl -pi -e 's/subcategory:$$/subcategory: Workflows/g' {} +

release:
	@echo "Note: Actual release building is handled by CI/GoReleaser"
	@echo "This target is for local development snapshot builds only"
	goreleaser release --rm-dist --snapshot --skip-publish --skip-sign

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

tools:
	cd providerlint && go install .
	cd tools && go install github.com/bflad/tfproviderdocs
	cd tools && go install github.com/client9/misspell/cmd/misspell
	cd tools && go install github.com/golangci/golangci-lint/cmd/golangci-lint
	cd tools && go install github.com/katbyte/terrafmt
	cd tools && go install github.com/terraform-linters/tflint
	cd tools && go install github.com/pavius/impi/cmd/impi
	cd tools && go install github.com/hashicorp/go-changelog/cmd/changelog-build

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=5m -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

sweeper:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test ./internal/provider -v -tags=sweep -sweep=all -sweep-allow-failures -timeout 120m

codegen:
	curl $(SWAGGER_URL) -o schema/swagger.json
	node tools/clean-swagger.js schema/swagger.json
	cd schema && go tool oapi-codegen --config=oapi-config.yml swagger.json
	yarn run generate schema/swagger.json
	go fmt client/*
	go fmt provider/*
	go tool goimports -w provider/*
	go tool goimports -w client/*

codegen-resource:
	@if [ -z "$(RESOURCE)" ]; then \
		echo "Error: RESOURCE parameter is required. Usage: make codegen-resource RESOURCE=service"; \
		exit 1; \
	fi
	curl $(SWAGGER_URL) -o schema/swagger.json
	node tools/clean-swagger.js schema/swagger.json
	cd schema && go tool oapi-codegen --config=oapi-config.yml swagger.json
	yarn run generate schema/swagger.json $(RESOURCE)
	@RESOURCE_PLURAL=$$(node -e "const inflect = require('./tools/inflect'); console.log(inflect.pluralize('$(RESOURCE)'))"); \
	go fmt client/$${RESOURCE_PLURAL}.go 2>/dev/null || true; \
	go fmt provider/resource_$(RESOURCE).go provider/data_source_$(RESOURCE).go 2>/dev/null || go fmt provider/resource_$(RESOURCE).go 2>/dev/null || true; \
	go tool goimports -w provider/resource_$(RESOURCE).go provider/data_source_$(RESOURCE).go 2>/dev/null || go tool goimports -w provider/resource_$(RESOURCE).go 2>/dev/null || true
	@echo ""
	@echo "✅ Code generation complete for $(RESOURCE)"
	@echo "📝 Files that may have changed:"
	@RESOURCE_PLURAL=$$(node -e "const inflect = require('./tools/inflect'); console.log(inflect.pluralize('$(RESOURCE)'))"); \
	echo "   - client/$${RESOURCE_PLURAL}.go"; \
	echo "   - provider/resource_$(RESOURCE).go"; \
	echo "   - provider/data_source_$(RESOURCE).go (if exists)"; \
	echo ""; \
	echo "💡 Tip: Use 'git add -p' to selectively stage only the changes you want"

# Version management targets
# These targets manage semantic versioning using git tags
.PHONY: version-patch version-minor version-major version-show version-next version-help

version-show:
	@echo "Current version: $$(git describe --tags --abbrev=0 2>/dev/null || echo 'No tags found')"
	@echo "Next patch: $$(scripts/bump-version.sh show patch)"
	@echo "Next minor: $$(scripts/bump-version.sh show minor)"
	@echo "Next major: $$(scripts/bump-version.sh show major)"

version-patch:
	@scripts/bump-version.sh patch

version-minor:
	@scripts/bump-version.sh minor

version-major:
	@scripts/bump-version.sh major

version-next:
	@scripts/bump-version.sh show patch

version-help:
	@scripts/bump-version.sh help

# Release targets - these create git tags which trigger CI releases
.PHONY: release-patch release-minor release-major

release-patch: version-patch
	@echo "✅ Tag $$(git describe --tags --abbrev=0) pushed"
	@echo "🚀 CI will automatically build and publish the release"

release-minor: version-minor
	@echo "✅ Tag $$(git describe --tags --abbrev=0) pushed"
	@echo "🚀 CI will automatically build and publish the release"

release-major: version-major
	@echo "✅ Tag $$(git describe --tags --abbrev=0) pushed"
	@echo "🚀 CI will automatically build and publish the release"

# Help target to show available version commands
help-version:
	@echo ""
	@echo "Version Management Commands:"
	@echo "  make version-show     - Show current and next versions"
	@echo "  make version-patch    - Bump patch version (1.2.3 → 1.2.4)"
	@echo "  make version-minor    - Bump minor version (1.2.3 → 1.3.0)"
	@echo "  make version-major    - Bump major version (1.2.3 → 2.0.0)"
	@echo "  make version-next     - Show next patch version"
	@echo "  make version-help     - Show detailed version help"
	@echo ""
	@echo "Release Commands (bump version + push tag, CI builds release):"
	@echo "  make release-patch    - Bump patch and push tag (triggers CI release)"
	@echo "  make release-minor    - Bump minor and push tag (triggers CI release)"
	@echo "  make release-major    - Bump major and push tag (triggers CI release)"
	@echo ""

# General help target
help:
	@echo "Available targets:"
	@echo ""
	@echo "Build & Development:"
	@echo "  make build           - Generate code and build provider"
	@echo "  make codegen         - Download schema and regenerate code"
	@echo "  make docs            - Generate documentation"
	@echo "  make test            - Run unit tests"
	@echo "  make testacc         - Run acceptance tests"
	@echo "  make sweeper         - Delete resources created by acceptance tests. They have names starting with tf-"
	@echo "  make install         - Install provider locally"
	@echo "  make release         - Create local snapshot build"
	@echo ""
	@make help-version
