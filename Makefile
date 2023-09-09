NAME = cmpdb
PREFIX = /usr/local
MANPREFIX = $(PREFIX)/share/man

$(NAME): clean
	go build -o $(NAME) .

clean:
	rm -f $(NAME)

install: $(NAME)
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	install -m 755 ./$(NAME) $(DESTDIR)$(PREFIX)/bin/

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/$(NAME)

benchmark:
	@go test -v -bench=. -run=^#

test:
	@go test -v ./...

.PHONY: clean install uninstall test benchmark
