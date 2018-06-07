watchkiller
===========

# Example

```
FROM nownabe/watchkiller:1.0 AS watchkiller

FROM gcr.io/cloudsql-docker/gce-proxy:1.11

COPY --from=watchkiller /watchkiller /watchkiller
```
