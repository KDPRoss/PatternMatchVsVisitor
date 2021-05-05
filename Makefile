MOTMOT=Motmot

.PHONY: run run-golang run-motmot
run: run-golang run-motmot

run-golang:
	go run ArithLang.go

run-motmot:
	$(MOTMOT) -run ArithLang

