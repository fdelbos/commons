
test:
	@echo "Running tests..."
	go test \
		-count=1 \
		-cover \
		-race \
		-timeout 240s \
		./...

define mock
@echo mocking: $(1):$(2) '->' $(4):$(3)
@mockery \
	--dir=$(1) \
	--name=$(2) \
	--structname=$(3) \
	--outpkg mocks \
	--output ./internal/mocks \
	--filename $(4)
endef

mocks:
	@rm -fr ./internal/mocks && mkdir ./internal/mocks
	$(call mock,auth,CodeStore,AuthCodeStore, auth_code_store.go)
	$(call mock,auth,Mailer,AuthMailer, auth_mailer.go)
	$(call mock,auth,SessionsStore,AuthSessionsStore, auth_sessions_store.go)