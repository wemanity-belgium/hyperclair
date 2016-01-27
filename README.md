
[![Build Status](https://travis-ci.org/jgsqware/hyperclair.svg?branch=develop)](https://travis-ci.org/jgsqware/hyperclair)
# hyperclair
A CLI tool for using CoreOS Clair with Docker Registry.

> The Registry is a stateless, highly scalable server side application that stores and lets you distribute Docker images. The Registry is open-source, under the permissive Apache license.
>
>*From https://docs.docker.com/registry/*

> Clair is a container vulnerability analysis service. It provides a list of vulnerabilities that threaten a container, and can notify users when new vulnerabilities that affect existing containers become known.
>
>*From https://github.com/coreos/clair*

hyperclair is tool to make the link between the Docker Registry and the CoreOS Clair tool.

# Notification
There is two way:

1. Automatic: Registry notify through event **hyperclair** when a new image is pulled, then **hyperclair** push it to Clair

2. On-Demand: the CLI tool is used to pull image from Registry then push it to Clair

# Reporting

**hyperclair** get vulnerabilities report from Clair and generate and HTML report

hyperclair can be used for Docker Hub and for Personal Registry

# Command

```
hyperclair --help
Analyse your docker image with Clair, directly from your registry.

Usage:
  hyperclair [command]

Available Commands:
  analyse     Analyse Docker image
  pull        Pull Docker image information
  push        Push Docker image to Clair
  report      Generate Docker Image vulnerabilities report
  serve       Create hyperclair Server

Flags:
      --config string   config file (default is ./.hyperclair.yml)

Use "hyperclair [command] --help" for more information about a command.

```

# Configuration

```yaml
clair:
  port: 6060
  uri: clair
  link: registry
  priority: Low
  report:
    path: reports
    format: html
auth:
  uri: http://172.17.0.1:3000/v2/token
```

# Contribution

Go to /contrib
