.PHONY:run send

SCRIPT_RUN := cmd/main.go
run:
	go run $(CRIPT_RUN)

SCRIPT_SEND := publisher/main.go
send:
	go run $(SCRIPT_SEND)