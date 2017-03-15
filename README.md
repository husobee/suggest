# Suggest

Suggest is an auto-suggestion micro-service which will take a stream of input
phrases, and provide a REST API for auto-suggestions of inputs.

## Running Tests

```bash
go test ./... -gcflags -l -cover
```

## Building Server

```bash
go build github.com/husobee/suggest/cmd/suggest
```

## Running Server

`./suggest`

suggest uses glog for logging, which will allow for verbosity levels on the 
command line: `./suggest -logtostderr -v 1`

suggest also accepts multiple environment variables:

```
ADDR - server address to listen on, defaults ":8080"
CORS_ALLOW_CREDENTIALS
CORS_ALLOW_HEADERS
CORS_ALLOW_ORIGIN
CORS_EXPOSE_HEADERS
CORS_MAX_AGE
TRACE - set allow trace on endpoints
TRIE_INSERTION_BUFFER - size of the trie insertion buffer, defaults 1
TRIE_RETRIEVE_BUFFER - size of the trie retrieve buffer, defaults 1
```

## Examples

Insertion:

```
curl -XPOST http://localhost:8080/ -d{"key":"projects"}
{"status":"OK","message":"successful insertion of term"}
curl -XPOST http://localhost:8080/ -d{"key":"protobuf"}
{"status":"OK","message":"successful insertion of term"}
curl -XPOST http://localhost:8080/ -d{"key":"puzzle"}
{"status":"OK","message":"successful insertion of term"}
curl -XPOST http://localhost:8080/ -d{"key":"python"}
{"status":"OK","message":"successful insertion of term"}
```

Retrieve:

```
curl http://127.0.0.1:8080/?key=pro
{"status":"OK","message":"successful in retrieving results","payload":[{"key":"projects","value":null},{"key":"protobuf","value":null}]}
curl http://127.0.0.1:8080/?key=proj
{"status":"OK","message":"successful in retrieving results","payload":[{"key":"projects","value":null}]}
```
