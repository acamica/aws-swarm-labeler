Swarm Labeler
---

This service will use AWS sdk to find all instances from a cloudformation stack, retrieve every tag from them and apply those tags to the corresponding docker swarm node using the private dns hostname.

Usage
---
```
./aws_swarm_labeler -cron <cron-exp> -filter <filter-exp> -region <regsion-name> -stack <stack-name>
```
 * **-cron** 
    	cron expression, like '@every 5m' for every five minutes or '15 * * * *' for every minute at 15th second. See [docs](https://godoc.org/github.com/robfig/cron)
 * **-filter** 
    	filter tag regex (default ".*")
 * **-region**
    	aws region (default "us-east-1")
 * **-stack** 
    	cloudformation stack name (required)


AWS permissions
---
 * DescribeStackResources
 * DescribeAutoScalingGroups
 * DescribeInstances

Docker 
---
AWS credentials and a connection to the docker socket should be provided.

Example:
```
docker run -it \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v ${HOME}/.aws/credentials:/root/.aws/credentials \
	aws_swarm_labeler \
        -stack prod -cron '@every 30m'
```

