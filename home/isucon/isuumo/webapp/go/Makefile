deploy: stop isuumo start status

.PHONY: stop start status

stop:
	-sudo systemctl stop isuumo.go

start:
	sudo systemctl start isuumo.go

status:
	sudo systemctl status isuumo.go

isuumo: *.go
	go build -o isuumo
