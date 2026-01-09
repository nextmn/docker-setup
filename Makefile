prefix = /usr/local
exec_prefix = $(prefix)
bindir = $(exec_prefix)/bin
BASHCOMPLETIONSDIR = $(exec_prefix)/share/bash-completion/completions


RM = rm -f
INSTALL = install -D
MKDIRP = mkdir -p

.PHONY: install uninstall update build clean default
default: build
build:
	go build
clean:
	go clean
reinstall: uninstall install
install:
	$(INSTALL) docker-setup $(DESTDIR)$(bindir)/docker-setup
	$(MKDIRP) $(DESTDIR)$(BASHCOMPLETIONSDIR)
	$(DESTDIR)$(bindir)/docker-setup completion bash > $(DESTDIR)$(BASHCOMPLETIONSDIR)/docker-setup
	@echo "================================="
	@echo ">> Now run the following command:"
	@echo "\tsource $(DESTDIR)$(BASHCOMPLETIONSDIR)/docker-setup"
	@echo "================================="
uninstall:
	$(RM) $(DESTDIR)$(bindir)/docker-setup
	$(RM) $(DESTDIR)$(BASHCOMPLETIONSDIR)/docker-setup
