# gowatcher for stand alone docker hosts in aws autoscaling groups

This was created for a specific use case of stand alone docker hosts that are running in aws in an autoscaling group.
It has a filter for the containers to watch once all the containers finishes it will detach and reduce autoscale desired count, and will terminate it self.

The following container needs to have the following environment variables set.

AWS Permissions needs to have autoscaling access.
>If you are not using IAM_Role use environment variables when running the containers
>AWS_ACCESS_KEY_ID, AWS_ACCESS_KEY

```bash
FILTER='container image to filter in for watching'
TIME_INTERVAL='30s' # Optional default to 60
DEBUG=true # Optional, when stated terminate action will be skipped
```

To build:

```bash
docker build --no-cache -t omerha/gowatcher:latest .
```

To Run it:

```bash
docker run --memory 30M --cpus 0.1 -v /var/run/docker.sock:/var/run/docker.sock:ro -e FILTER=container_to_watch --restart always omerha/gowatcher:latest

## Optional -e DEBUG=true to skip termination
## Optional -e TIME_INTERVAL=10s to reduce ticks
```