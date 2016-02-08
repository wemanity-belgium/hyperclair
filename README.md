

# hyperclair

[![Build Status](https://travis-ci.org/wemanity-belgium/hyperclair.svg?branch=develop)](https://travis-ci.org/wemanity-belgium/hyperclair) [![Join the chat at https://gitter.im/wemanity-belgium/hyperclair](https://badges.gitter.im/wemanity-belgium/hyperclair.svg)](https://gitter.im/wemanity-belgium/hyperclair?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

A CLI tool for using CoreOS Clair with Docker Registry.

> The Registry is a stateless, highly scalable server side application that stores and lets you distribute Docker images. The Registry is open-source, under the permissive Apache license.
>
>*From https://docs.docker.com/registry/*

> Clair is a container vulnerability analysis service. It provides a list of vulnerabilities that threaten a container, and can notify users when new vulnerabilities that affect existing containers become known.
>
>*From https://github.com/coreos/clair*

hyperclair is tool to make the link between the Docker Registry and the CoreOS Clair tool.

![hyperclair](https://cloud.githubusercontent.com/assets/3304363/12849755/9caa0fac-cc21-11e5-8b89-ddfa8535a3dc.png)

# Notification
1. Api: `hyperclair serve` run a web server to interact with the Registry and Clair. It play as Reverse Proxy for Registry Authentication
2. On-Demand: the CLI tool is used to pull image from Registry then push it to Clair

# Reporting

**hyperclair** get vulnerabilities report from Clair and generate HTML report

hyperclair can be used for Docker Hub and self-hosted Registry

# Command

```
Analyse your docker image with Clair, directly from your registry.

Usage:
  hyperclair [command]

Available Commands:
  analyse     Analyse Docker image
  health      Get Health of Hyperclair and underlying services
  pull        Pull Docker image information
  push        Push Docker image to Clair
  report      Generate Docker Image vulnerabilities report
  serve       Create hyperclair Server
  version     Get Versions of Hyperclair and underlying services

Flags:
      --config string   config file (default is ./.hyperclair.yml)
  -h, --help            help for hyperclair

Use "hyperclair [command] --help" for more information about a command.

```

# Configuration

```yaml
clair:
  port: 6060
  uri: http://clair
  priority: Low
  report:
    path: reports
    format: html
auth:
  user: jgsqware
  password: jgsqware
  insecureSkipVerify: true
hyperclair:
  uri: http://hyperclair
  port: 9999
```

# Contribution and Test

Go to /contrib folder
