# General

## Fetch API Version

It returns the versions of the Hyperclair api and connected Clair. We support onyl V2 of Docker Registry

	GET /v1/versions



### Example

```
curl -s http://localhost:9999/v1/versions
```

### Response

```
HTTP/1.1 200 OK
{
  "APIVersion": "1",
  "Clair": {
    "APIVersion": "1",
    "EngineVersion": "1"
  }
}
```

## Fetch Health status

	GET /v1/health

Returns 200 if essential services are healthy (ie. database) and 503 otherwise.

### Example

```
curl -s http://localhost:9999/v1/health
```

### Success Response

```
HTTP/1.1 200 OK
{
  "clair": {
    "database": {
      "IsHealthy": true
    },
    "updater": {
      "Details": {
        "LatestSuccessfulUpdate": "2016-02-05T10:04:26Z"
      },
      "IsHealthy": true
    }
  },
  "database": {
    "IsHealthy": true
  }
}
```

### Error Response

```
HTTP/1.1 503 Service unavailable
{
  "clair": {
    "database": {
      "IsHealthy": true
    },
    "updater": {
      "Details": {
        "LatestSuccessfulUpdate": "2016-02-05T10:04:26Z"
      },
      "IsHealthy": true
    }
  },
  "database": {
    "IsHealthy": false
  }
}
```

## Pull Image from Registry

Return a light Manifest version of docker image

  GET /v1/<name>?realm=<registry (default: Docker Hub)>&reference=<reference>

The name and reference parameter identify the image and are required. The reference may include a tag or digest.

## Push Image to Clair

  POST /v1/<name>?realm=<registry (default: Docker Hub)>&reference=<reference>

## Get Image analysis as JSON

  GET /v1/<name>/<reference>/analysis

## Get Image analysis report as HTML

  GET /v1/<name/<reference>/analysis/report
