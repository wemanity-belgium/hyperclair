# General

## Fetch API Version

It returns the versions of the Hyperclair api and connected Clair. We support onyl V2 of Docker Registry

	GET /v1/versions



### Example

```
curl -s 127.0.0.1:7070/v1/versions
```

### Response

```
HTTP/1.1 200 OK
{
  "APIVersion": "1",
  "ClairVersion": {
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
curl -s 127.0.0.1:7070/v1/health | python -m json.tool
```

### Success Response

```
HTTP/1.1 200 OK
{  
   "database":{  
      "IsHealthy":true
   },
   "clair":{
     "database":{  
        "IsHealthy":true
     },
     "notifier":{  
        "IsHealthy":true,
        "Details":{  
           "QueueSize":0
        }
     },
     "updater":{  
        "IsHealthy":true,
        "Details":{  
           "HealthIdentifier":"cf65a8f6-425c-4a9c-87fe-f59ddf75fc87",
           "HealthLockOwner":"1e7fce65-ee67-4ca5-b2e9-61e9f5e0d3ed",
           "LatestSuccessfulUpdate":"2015-09-30T14:47:47Z",
           "ConsecutiveLocalFailures":0
        }
     }
   }
}
```

### Error Response

```
HTTP/1.1 503 Service unavailable
{  
  "database":{  
     "IsHealthy":false
  },
  "clair":{
    "database":{  
       "IsHealthy":true
    },
    "notifier":{  
       "IsHealthy":true,
       "Details":{  
          "QueueSize":0
       }
    },
    "updater":{  
       "IsHealthy":true,
       "Details":{  
          "HealthIdentifier":"cf65a8f6-425c-4a9c-87fe-f59ddf75fc87",
          "HealthLockOwner":"1e7fce65-ee67-4ca5-b2e9-61e9f5e0d3ed",
          "LatestSuccessfulUpdate":"2015-09-30T14:47:47Z",
          "ConsecutiveLocalFailures":0
       }
    }
  }
}
```

## Pull Image from Registry

Return a light Manifest version of docker image

  GET /v1/<name/<reference>

The name and reference parameter identify the image and are required. The reference may include a tag or digest.

## Push Image to Clair

  POST /v1/<name>/<reference>

## Get Image analysis as JSON

  GET /v1/<name>/<reference>/analysis

## Get Image analysis report as HTML

  GET /v1/<name/<reference>/analysis/report
