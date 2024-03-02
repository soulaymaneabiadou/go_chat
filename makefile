s:
	cd server && go run .

c:
	cd client && go run .

t:
	telnet localhost 7007

.PHONY: s c t