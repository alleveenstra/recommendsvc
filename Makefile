main: src/httpd.go
	go build src/httpd.go
	cp src/httpd recommendsvc
