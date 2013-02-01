main: src/httpd.go
	go install recommendsvc
	go build src/httpd.go

