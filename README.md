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
