deploy: stop isuumo start status

.PHONY: stop start status

stop:
	-systemctl start isuumo.go

start:
	systemctl start isuumo.go

status:
	systemctl status isuumo.go

isuumo: *.go
	go build -o isuumo
