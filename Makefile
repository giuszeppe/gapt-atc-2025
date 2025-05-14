# run the server
run:
	$(MAKE) build && $(MAKE) server & $(MAKE) client
server:
	cd backend && ./server server
client:
	cd frontend && npm run dev
build:
	cd backend && go build -o server
run-refresh:
	$(MAKE) build && $(MAKE) server-refresh && $(MAKE) run
server-refresh:
	cd backend && ./server database refresh