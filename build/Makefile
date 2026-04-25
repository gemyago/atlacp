.PHONY: FORCE
.SECONDEXPANSION:
.SECONDARY:

# Prevent Make from deleting built platform directories after packaging succeeds.
.PRECIOUS: dist/%

app_name:=$(shell basename `git rev-parse --show-toplevel`)

# Define a function to read values from any config file
# Usage: $(call read_config,config_file,key_name)
define read_config
$(shell grep -m 1 "^$(2)\s*=\s*" $(1) | cut -d'=' -f2 | sed 's/^[[:space:]]*//;s/[[:space:]]*$$//')
endef

# Reading config using the new function
docker_runtime_image:=$(call read_config,build.cfg,docker_runtime_image)
platforms:=$(call read_config,build.cfg,platforms)

comma:=,

# Detect OS and ARCH for host machine
os_lower := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ifeq ($(os_lower),darwin)
	detected_os ?= Darwin
else
	detected_os ?= Linux
endif

# All Go sources participating in binary builds (across the repository tree).
GO_SOURCES := $(strip $(shell find .. -name '*.go' -type f))

detected_arch ?= $(shell uname -m)
ifeq ($(detected_arch),x86_64)
	detected_arch := x86_64
else ifeq ($(detected_arch),aarch64)
	detected_arch := arm64
else ifeq ($(detected_arch),arm64)
	detected_arch := arm64
endif

# Tool Binaries (local to build directory)
bin_dir := $(CURDIR)/bin
tmp_dir := $(CURDIR)/tmp

# Crane Installation
crane_version := $(shell cat .crane-version)
crane_archive_name := go-containerregistry_$(detected_os)_$(detected_arch).tar.gz
crane_archive_url := https://github.com/google/go-containerregistry/releases/download/$(crane_version)/$(crane_archive_name)
crane_archive := $(tmp_dir)/$(crane_archive_name)
crane := $(bin_dir)/crane

$(tmp_dir):
	mkdir -p $(tmp_dir)

$(bin_dir):
	mkdir -p $(bin_dir)

$(crane_archive): .crane-version $(tmp_dir)
	@echo "Downloading crane v$(crane_version) for $(detected_os)/$(detected_arch)..."
	curl -sL $(crane_archive_url) -o $(crane_archive)

$(crane): $(crane_archive) $(bin_dir)
	@echo "Extracting crane binary..."
	tar -zxvf $(crane_archive) -C $(bin_dir)/ crane
	@echo "crane installed to $(crane)"
	touch $(crane) # Update timestamp for Make

# Target to ensure crane is installed
.PHONY: install-crane
install-crane: $(crane)

# platforms in config are in form linux/amd64,linux/arm64...
# transforming to dist/linux/amd64 dist/linux/arm64...
dist_platforms:=$(patsubst %,dist/%,$(subst $(comma), ,$(platforms)))

# Example: dist/linux/amd64
# If you care of the size, add -ldflags="-s -w"
# but keep in mind that profiling will be challenging.
# Disabled CGO explicitly speeds up the build. Remove if you need CGO.
dist/%: $(GO_SOURCES) ../go.mod ../go.sum
	$(eval GOOS:=$(word 1,$(subst /, ,$*)))
	$(eval GOARCH:=$(word 2,$(subst /, ,$*)))
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -C .. \
			-tags=release \
			-o build/dist/$(GOOS)/$(GOARCH)/ ./cmd/...;

dist: FORCE
	@make -j4 $(dist_platforms)

# npm package generation (template-driven, generated under build/npm/packages/)
npm_templates_dir := npm/templates
npm_packages_dir := npm/packages
npm_packs_dir := npm/packages/packs
npm_install_src_dir := npm/install
npm_platform_package_json := $(patsubst %,npm/packages/%/package.json,$(subst $(comma), ,$(platforms)))
npm_install_package_json := npm/packages/install/package.json
npm_pack_targets := $(foreach p,$(subst $(comma), ,$(platforms)),npm/packages/packs/atlacp-install-$(subst /,-,$(p))-$(VERSION).tgz) \
	npm/packages/packs/atlacp-install-$(VERSION).tgz
npm_pack_published := $(addsuffix .published,$(npm_pack_targets))
npm_pack_unpublish := $(addsuffix .unpublished,$(npm_pack_targets))
npm_install_sources := $(wildcard $(npm_install_src_dir)/*.mjs)
npm_generated_install_sources := $(patsubst $(npm_install_src_dir)/%,$(npm_packages_dir)/install/%,$(npm_install_sources))
npm_publish_flags ?= --provenance --access public
NPM_PUBLISH_NOOP ?= true
NPM_UNPUBLISH_NOOP ?= true
.SECONDARY: $(npm_platform_package_json) npm/packages/install/package.json

# Convert platform tuple to template suffix (linux/amd64 -> linux-amd64)
define platform_template_suffix
$(subst /,-,$(1))
endef

define npm_platform_pack_rule
$(2): $(npm_packages_dir)/$1/package.json $(npm_packages_dir)/.version
	@mkdir -p "$(npm_packs_dir)"
	@cd "$(npm_packages_dir)/$1" && npm pack --json > "$(abspath $(2)).json"
	@mv "$(npm_packages_dir)/$1/`cat "$(abspath $(2)).json" | jq -r '.[0].filename'`" "$(2)"; \
	rm -f "$(abspath $(2)).json"
	@touch "$(2)"
endef

$(foreach platform,$(subst $(comma), ,$(platforms)),$(eval $(call npm_platform_pack_rule,$(platform),$(npm_packs_dir)/atlacp-install-$(subst /,-,$(platform))-$(VERSION).tgz)))


# Install package tarball.
$(npm_packs_dir)/atlacp-install-$(VERSION).tgz: $(npm_install_package_json) $(npm_packages_dir)/.version
	@mkdir -p "$(npm_packs_dir)"
	@cd "$(npm_packages_dir)/install" && npm pack --json > "$(abspath $@).json"
	@mv "$(npm_packages_dir)/install/`cat "$(abspath $@).json" | jq -r '.[0].filename'`" "$@"; \
	rm -f "$(abspath $@).json"
	@touch "$@"

npm/packages/.version: FORCE
	@if [ -z "$(VERSION)" ]; then \
		echo "VERSION is required. Usage: make npm/packages VERSION=1.2.3"; exit 1; \
	fi
	@mkdir -p $(dir $@)
	@printf "$(VERSION)" > $@.tmp
	@if cmp -s $@ $@.tmp; then \
		echo "Using cached version: $(VERSION)"; \
		rm $@.tmp; \
	else \
		echo "New version detected: $(VERSION)"; \
		mv $@.tmp $@; \
	fi
 
# Generate and pack one platform package workspace.
# Example: make npm/packages/linux/amd64/package.json VERSION=1.2.3
npm/packages/%/package.json: dist/% npm/packages/.version
	$(eval GOOS := $(word 1,$(subst /, ,$*)))
	$(eval GOARCH := $(word 2,$(subst /, ,$*)))
	$(eval TEMPLATE := $(npm_templates_dir)/package-$(call platform_template_suffix,$*).json)
	@test -f "$(TEMPLATE)" || (echo "Missing template: $(TEMPLATE)"; exit 1)
	mkdir -p $(dir $@)
	jq '.version = "$(VERSION)"' "$(TEMPLATE)" > $@
	rm -rf $(npm_packages_dir)/$*/bin
	mkdir -p $(npm_packages_dir)/$*/bin
	cp dist/$(GOOS)/$(GOARCH)/* $(npm_packages_dir)/$*/bin/
	chmod +x $(npm_packages_dir)/$*/bin/*
	cd $(dir $@) && npm pack

$(npm_packages_dir)/install/%: $(npm_install_src_dir)/%
	mkdir -p $(dir $@)
	cp $< $@

npm/packages/install/package.json: $(npm_templates_dir)/package-install.json $(npm_generated_install_sources) npm/packages/.version
	mkdir -p $(dir $@)
	jq '.version = "$(VERSION)" | .optionalDependencies |= with_entries(.value = "$(VERSION)")' \
		$< > $@

npm/packages/install/.packed: npm/packages/install/package.json $(npm_generated_install_sources)
	cd $(dir $<) && npm pack
	touch $@

.PHONY: npm/packages npm-packages
npm/packages: npm/packages/check-version $(npm_platform_packages) npm/packages/install/.packed
npm-packages: npm/packages

	cd $(dir $@) && npm pack

npm/packs.txt: npm/packages/.version $(npm_pack_targets)
	echo "$(npm_pack_targets)" | tr " " "\n" | sort > $@

$(npm_packs_dir)/%.tgz.published: $(npm_packs_dir)/%.tgz
ifeq ($(NPM_PUBLISH_NOOP),false)
	npm publish $(npm_publish_flags) $<
else
	@echo "NPM_PUBLISH_NOOP=true: skipping publish for $<"
endif
	touch $@

.PHONY: npm/publish
npm/publish: npm/packages/.version $(npm_pack_published)

.PHONY: npm/unpublish
npm/unpublish: npm/packages/.version $(npm_pack_unpublish)

$(npm_packs_dir)/%.tgz.unpublished: $(npm_packs_dir)/%.tgz
	@spec_and_version=$$(tar -xOf "$<" package/package.json | jq -r '.name + "@" + .version'); \
	if [ "$(NPM_UNPUBLISH_NOOP)" = "false" ]; then \
		npm unpublish "$${spec_and_version}"; \
	else \
		echo "NPM_UNPUBLISH_NOOP=true: skipping unpublish for $< ($${spec_and_version})"; \
	fi
	@touch "$@"

.PHONY: npm/clean-published
npm/clean-published:
	rm -f $(npm_packs_dir)/*.published

.build-artifacts: dist
	@echo ./dist > $@

# Create a tarball of build artifacts
# Example: make build-artifacts.tar.gz
build-artifacts.tar.%: .build-artifacts
	rm -f $@
	tar -a -cf $@ $(shell cat $<)

.docker-artifacts: docker/.remote-images
	@rm -f $@
	@echo "docker/.remote-images" > $@
	@echo ./.git-meta >> $@

docker-artifacts.tar.%: .docker-artifacts
	rm -f $@
	tar -a -cf $@ $(shell cat $<)

# COMMIT_SHA is used since github may checkout on the temporary PR branch
# that will have a different commit SHA than the merge commit SHA
.git-meta/HEAD: FORCE
	@mkdir -p .git-meta
	@if [ -n "$$GIT_COMMIT_SHA" ]; then \
		echo "$$GIT_COMMIT_SHA" | cut -c1-7 > $@.tmp; \
	else \
		git rev-parse --short HEAD > $@.tmp; \
	fi
	@if cmp -s $@ $@.tmp; then \
		rm $@.tmp; \
	else \
		mv $@.tmp $@; \
	fi

# COMMIT_REF is used since github may checkout on the temporary PR branch
# that will always have HEAD as a ref name
.git-meta/ref: .git-meta/HEAD
	@if [ -n "$$GIT_REF" ]; then \
		echo "$$GIT_REF" > $@.tmp; \
	else \
		git rev-parse --abbrev-ref HEAD > $@.tmp; \
	fi
	@cmp -s $@.tmp $@ || mv $@.tmp $@
	@rm -f $@.tmp

.git-meta/repo-name: .git-meta/HEAD
	@basename $(shell git rev-parse --show-toplevel) > $@

.git-meta/remote-url: .git-meta/HEAD
	@git remote get-url origin | sed -E \
		-e 's|^git@([^:]+):|https://\1/|' \
		-e 's|^http:|https:|' \
		-e 's|\.git$$||' > $@.tmp
	@cmp -s $@.tmp $@ || mv $@.tmp $@
	@rm -f $@.tmp

.git-meta/repo-namespace: .git-meta/remote-url
	@cat .git-meta/remote-url | sed -E \
		-e 's|^https://[^/]+/||' \
		-e 's|/[^/]+$$||' > $@.tmp
	@cmp -s $@.tmp $@ || mv $@.tmp $@
	@rm -f $@.tmp

docker/.labels: .git-meta/HEAD .git-meta/repo-name .git-meta/remote-url .git-meta/ref
	@echo "org.opencontainers.image.created=\"$(shell git show -s --format=%cI HEAD)\"" > $@
	@echo "org.opencontainers.image.source=\"$(shell cat .git-meta/remote-url)\"" >> $@
	@echo "org.opencontainers.image.version=\"$(shell git describe --tags --always)\"" >> $@
	@echo "org.opencontainers.image.revision=\"$(shell cat .git-meta/HEAD)\"" >> $@
	@echo "org.opencontainers.image.ref.name=\"$(shell cat .git-meta/ref)\"" >> $@
	@echo "org.opencontainers.image.title=\"$(shell cat .git-meta/repo-name)\"" >> $@
	@echo "org.opencontainers.image.url=\"$(shell cat .git-meta/remote-url)\"" >> $@

docker/.local-args: docker/.labels
	@echo "--build-arg RUNTIME_IMAGE=$(call read_config,build.cfg,docker_runtime_image)" > $@
	@echo --platform $(platforms) >> $@
	@echo -f docker/Dockerfile >> $@
	@sed 's/^/--label /' docker/.labels >> $@
	@echo "../" >> $@

docker/.local-%-tags: .git-meta/HEAD .git-meta/repo-name .git-meta/ref
	$(eval registry := $(call read_config,build.cfg,docker_local_registry))
	$(eval image_base_name := $(registry)/$(shell cat .git-meta/repo-name)-$*)
	$(eval tag_names=$(shell ./scripts/resolve-docker-tags.sh \
		--latest \
		--commit-sha $(shell cat .git-meta/HEAD) \
		--git-ref $(shell cat .git-meta/ref) \
		--stable-branches $(call read_config,build.cfg,stable_branches) \
	))
	@echo $(tag_names) | tr ' ' '\n' | sed 's|^|--tag $(image_base_name):|' > $@

# Build local image
# Target must be in form: docker/.local-<binary>-image
# Example: make docker/.local-server-image
docker/.local-%-image: docker/.local-args docker/.local-%-tags
	@echo "Building local $* image"
	docker buildx build --load --build-arg TARGET_BINARY="$*" \
			--metadata-file $@ \
			--annotation org.opencontainers.image.description="Includes $* binary" \
			$(shell cat docker/.local-$*-tags) \
			$(shell cat docker/.local-args)

# Build all local images for each binary in dist/
# Fill use top folders from cmd folder to understand which binaries to build (wildcard)
.PHONY: docker/.local-images
docker/.local-images: $(patsubst %, docker/.local-%-image, $(notdir $(wildcard ../cmd/*)))
	@echo "Local images are ready: $?"

# Build image tags for remote registry
# Target must be in form: docker/.remote-<registry>-<binary>-tags
# Example: make docker/.remote-ghcr-server-tags
# The registry host is read from build.cfg under docker_<registry>_registry key
docker/.remote-%-tags: .git-meta/HEAD .git-meta/repo-name .git-meta/repo-namespace .git-meta/ref
	$(eval registry_name = $(word 1,$(subst -, ,$*)))
	$(eval target_binary = $(word 2,$(subst -, ,$*)))
	$(eval registry := $(call read_config,build.cfg,docker_$(registry_name)_registry))
	$(eval image_base_name := $(registry)/$(shell cat .git-meta/repo-namespace)/$(shell cat .git-meta/repo-name)-$(target_binary))
	$(eval tag_names=$(shell ./scripts/resolve-docker-tags.sh \
		--commit-sha $(shell cat .git-meta/HEAD) \
		--git-ref $(shell cat .git-meta/ref) \
		--stable-branches $(call read_config,build.cfg,stable_branches) \
	))
	@echo $(tag_names) | tr ' ' '\n' | sed 's|^|--tag $(image_base_name):|' > $@

# Push image to remote registry
# Target must be in form: docker/.remote-<registry>-<binary>-image
# Example: make docker/.remote-ghcr-server-image
# The registry host is read from build.cfg under docker_<registry>_registry_host key
docker/.remote-%-image: docker/.local-args docker/.remote-%-tags
	$(eval registry_name := $(word 1,$(subst -, ,$*)))
	$(eval target_binary := $(word 2,$(subst -, ,$*)))
	@echo "Building and pushing remote $* image for $(target_binary)"
	docker buildx build --push --build-arg TARGET_BINARY="$(target_binary)" \
			--debug \
			--metadata-file $@ \
			--annotation org.opencontainers.image.description="Includes $(target_binary) binary" \
			$(shell cat docker/.remote-$*-tags) \
			$(shell cat docker/.local-args)

# Push all images to remote registry
.PHONY: docker/remote-images
docker/.remote-images: .git-meta/repo-name .git-meta/repo-namespace $(eval docker_push_registries := $(call read_config,build.cfg,docker_push_registries))
docker/.remote-images: $(eval binaries := $(notdir $(wildcard ../cmd/*)))
docker/.remote-images: $(eval remote_image_targets := $(patsubst %, docker/.remote-%-image, $(foreach registry_name,$(subst $(comma), ,$(docker_push_registries)),$(foreach binary,$(binaries),$(registry_name)-$(binary)))))
docker/.remote-images: $(remote_image_targets)
	@echo "Writing unique base image names to docker/.remote-images"
	@rm -f docker/.remote-images docker/remote-images.tmp
	@touch docker/remote-images.tmp
	@$(foreach registry_name,$(subst $(comma), ,$(docker_push_registries)), \
		$(eval registry := $(call read_config,build.cfg,docker_$(registry_name)_registry)) \
		$(foreach binary,$(binaries), \
			echo "$(registry)/$(shell cat .git-meta/repo-namespace)/$(shell cat .git-meta/repo-name)-$(binary)" >> docker/remote-images.tmp; \
		) \
	)
	@sort -u docker/remote-images.tmp > docker/.remote-images
	@rm -f docker/remote-images.tmp
	@echo "Remote images pushed: $?"

# Target to generate a file containing the list of remote image base names
# without building/pushing the images.
docker/.remote-image-names: .git-meta/repo-name .git-meta/repo-namespace
	@echo "Generating list of remote image names to $@..."
	@rm -f $@ $@.tmp
	@touch $@.tmp
	$(eval docker_push_registries := $(call read_config,build.cfg,docker_push_registries))
	$(eval binaries := $(notdir $(wildcard ../cmd/*)))
	@$(foreach registry_name,$(subst $(comma), ,$(docker_push_registries)), \
		$(eval registry := $(call read_config,build.cfg,docker_$(registry_name)_registry)) \
		$(foreach binary,$(binaries), \
			echo "$(registry)/$(shell cat .git-meta/repo-namespace)/$(shell cat .git-meta/repo-name)-$(binary)" >> $@.tmp; \
		) \
	)
	@sort -u $@.tmp > $@
	@rm -f $@.tmp
	@echo "List of remote image names written to $@"

.PHONY: clean docker/clean dist/clean

docker/clean:
	rm -f docker/.local-* docker/.remote-* docker/.labels

dist/clean:
	rm -rf dist

clean: docker/clean dist/clean
	rm -rf .git-meta
	rm -f .build-artifacts
	rm -rf npm/packages
	rm -f build-artifacts.tar.*
	# Clean crane artifacts
	rm -rf $(bin_dir)
	rm -rf $(tmp_dir)

.PHONY: test
test:
	@echo "Running resolve-docker-tags self-test"
	scripts/resolve-docker-tags.sh --self-test

	@echo "Running python script tests"
	python -m unittest discover -s scripts/tests -p "test_*.py"

	@echo "Running npm/scripts resolve-semver-tag.sh self-test"
	./npm/scripts/resolve-npm-platform.sh --self-test

	@echo "Running npm/scripts parse-semver-tag.sh self-test"
	./npm/scripts/parse-semver-tag.sh --self-test
