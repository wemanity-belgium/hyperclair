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

	GET /v1/<name>?realm=<registry>&reference=<reference>

	Default:
		Registry: https://registry-1.docker.io/v2
		Reference: latest

The name and reference parameter identify the image and are required. The reference may include a tag or digest.

### Example

```
curl -s http://localhost:9999/v1/jgsqware/ubuntu-git?realm=registry:5000&reference=latest
```

### Success Response

```
HTTP/1.1 200 OK
{
  "Name": "jgsqware/ubuntu-git",
  "Tag": "latest",
  "Registry": "http://registry:5000/v2",
  "FsLayers": [
    {
      "BlobSum": "sha256:13be4a52fdee2f6c44948b99b5b65ec703b1ca76c1ab5d2d90ae9bf18347082e"
    },
    {
      "BlobSum": "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4"
    },
    {
      "BlobSum": "sha256:27aa681c95e5165caf287dcfe896532df4ae8b10e099500f2f8f71acf4002a89"
    },
    {
      "BlobSum": "sha256:9e0bc8a71bde464f710bc2b593a1fc21521517671e918687892303151331fa56"
    },
    {
      "BlobSum": "sha256:d89e1bee20d9cb344674e213b581f14fbd8e70274ecf9d10c514bab78a307845"
    }
  ]
}
```

## Push Image to Clair

	POST /v1/<name>?realm=<registry>&reference=<reference>

	Default:
		Registry: https://registry-1.docker.io/v2
		Reference: latest

Push image to Clair for analysis

### Example

```
curl -X POST -s http://localhost:9999/v1/jgsqware/ubuntu-git?realm=registry:5000&reference=latest
```

### Success Response

```
HTTP/1.1 201 Created
```

## Get Image analysis as JSON

	GET /v1/<name>/analysis?realm=<registry>&reference=<reference>

	Default:
		Registry: https://registry-1.docker.io/v2
		Reference: latest

The name and reference parameter identify the image and are required. The reference may include a tag or digest.

### Example

```
curl -s http://localhost:9999/v1/jgsqware/ubuntu-git/analysis?realm=registry:5000&reference=latest
```

### Success Response

```
HTTP/1.1 200 OK
{
  "Registry": "registry:5000",
  "ImageName": "jgsqware/ubuntu-git",
  "Tag": "latest",
  "Layers": [
    {
      "ID": "sha256:d89e1bee20d9cb344674e213b581f14fbd8e70274ecf9d10c514bab78a307845",
      "Vulnerabilities": [
        {
          "ID": "CVE-2015-1865",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2015-1865",
          "Priority": "Low",
          "Description": "\"time of check to time of use\" race condition fts.c",
          "CausedByPackage": "coreutils"
        },
        {
          "ID": "CVE-2014-8121",
          "Link": "https://security-tracker.debian.org/tracker/CVE-2014-8121",
          "Priority": "Medium",
          "Description": "DB_LOOKUP in nss_files/files-XXX.c in the Name Service Switch (NSS) in GNU C Library (aka glibc or libc6) 2.21 and earlier does not properly check if a file is open, which allows remote attackers to cause a denial of service (infinite loop) by performing a look-up while the database is iterated over the database, which triggers the file pointer to be reset.",
          "CausedByPackage": "eglibc"
        }
      ]
    },
    {
      "ID": "sha256:9e0bc8a71bde464f710bc2b593a1fc21521517671e918687892303151331fa56",
      "Vulnerabilities": []
    },
    {
      "ID": "sha256:27aa681c95e5165caf287dcfe896532df4ae8b10e099500f2f8f71acf4002a89",
      "Vulnerabilities": []
    },
    {
      "ID": "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4",
      "Vulnerabilities": []
    }
  ]
}
```

## Get Image analysis report as HTML

  GET /v1/<name>/analysis/report?realm=<registry>&reference=<reference>

	Default:
		Registry: https://registry-1.docker.io/v2
		Reference: latest

The name and reference parameter identify the image and are required. The reference may include a tag or digest.

### Example

```
curl -s http://localhost:9999/v1/jgsqware/ubuntu-git/analysis/report?realm=registry:5000&reference=latest
```

### Success Response

```
HTTP/1.1 200 OK
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <title>Hyperclair Reports</title>
</head>
<body>
    <div class="container-fluid">
        <div class="row">
            <div class="col-md-12">
                <div class="page-header">
                    <h1>
					Hyperclair reports
        </h1>
                    <h2>
          Image: registry:5000/jgsqware/ubuntu-git:latest
				</h2>

      [...]


  </body>
  </html>

```
