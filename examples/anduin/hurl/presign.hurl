GET http://localhost:5000/v2/hello/manifests/v1
[BasicAuth]
admin: password
HTTP 200
[Captures]
layer_digest: jsonpath "$['layers'][0]['digest']"

GET http://localhost:5000/v2/_zot/ext/blobpresign/hello/blobs/{{layer_digest}}
[BasicAuth]
admin: password
