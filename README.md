Short url

run:
go run bin/gin/main.go

to save long url:
POST to localhost:8080/url
with Body:
{
  "Full": "<http://google.com>"
}

to look list of urls:
GET <http://localhost:8080/url>

to use shortlink:
http::/locahost:8080/{id}
