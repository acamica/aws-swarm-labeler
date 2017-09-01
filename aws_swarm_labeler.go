package main

import (
    "flag"
    "golang.org/x/net/context"
    "github.com/robfig/cron"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudformation"
    "github.com/aws/aws-sdk-go/service/autoscaling"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/aws"
    "fmt"
    "regexp"
    "time"
)

func main() {
    var stackName string
    var regionName string
    var filter string
    var schedule string

    flag.StringVar(&stackName , "stack" ,""         ,"cloudformation stack name (required)")
    flag.StringVar(&regionName, "region","us-east-1","aws region")
    flag.StringVar(&filter    , "filter",".*"       ,"filter tag regex")
    flag.StringVar(&schedule  , "cron"  ,""         ,"cron expression, like '5 * * * *' for every five minutes. See [docs](https://godoc.org/github.com/robfig/cron)")
    flag.Parse()

    if "" == stackName {
        fmt.Println("stack name is required")
        return
    }

    filterRegex, err := regexp.Compile(filter)
    if err != nil { panic(err) }

    if "" == schedule {
        run(regionName, stackName, *filterRegex)
    } else {
        fmt.Println("Running croned updates:", schedule)

        c := cron.New()
        err := c.AddFunc(schedule, func() {
            run(regionName, stackName, *filterRegex)
        })
        if err != nil { panic(err) }

        c.Start()
        for true {
            time.Sleep(5000)
        }
    }
}

func run(regionName string, stackName string, filterRegex regexp.Regexp) {
    // setup aws session
    aws_session, err := session.NewSessionWithOptions(session.Options{
        Config: aws.Config{ Region: &regionName, },
    })
    if err != nil { panic(err) }

    cf := cloudformation.New(aws_session)
    as := autoscaling.New(aws_session)
    ec := ec2.New(aws_session)

    // setup docker client
    cli, err := client.NewEnvClient()
    if err != nil { panic(err) }

    fmt.Println(
        "Updating tags for", stackName,
        "with filter /", filterRegex.String(), "/",
        "at", time.Now().Format(time.ANSIC),
    )

    var groups []*string
    { // get auto scaling groups from stack name
        params := &cloudformation.DescribeStackResourcesInput{
            StackName: &stackName,
        }
        resp, err := cf.DescribeStackResources(params)
        if err != nil { panic(err) }

        for _, stackResource := range resp.StackResources {
            if "AWS::AutoScaling::AutoScalingGroup" == *stackResource.ResourceType {
                groups = append(groups, stackResource.PhysicalResourceId)
            }
        }

    }

    var instanceIds []*string
    { // get tags and instances from every auto scaling group
        params := &autoscaling.DescribeAutoScalingGroupsInput{
            AutoScalingGroupNames: groups,
        }
        resp, err := as.DescribeAutoScalingGroups(params)
        if err != nil { panic(err) }

        for _, asg := range resp.AutoScalingGroups {
            for _, instance := range asg.Instances {
                instanceIds = append(instanceIds, instance.InstanceId)
            }
        }
    }

    var instances = make(map[string]map[string]string)
    { //get instances to be tagged
        params := &ec2.DescribeInstancesInput{
            InstanceIds: instanceIds,
        }
        resp, err := ec.DescribeInstances(params)
        if err != nil { panic(err) }

        for _, re := range resp.Reservations {
            for _, instance := range re.Instances {
                tags := make(map[string]string)
                for _, tag := range instance.Tags {
                    if filterRegex.MatchString(*tag.Key) {
                        tags[*tag.Key] = *tag.Value
                    }
                }
                instances[*instance.PrivateDnsName] = tags
            }
        }
    }

    { // list all nodes in swarm and update them with the tags
        nodes, err := cli.NodeList(context.Background(), types.NodeListOptions{})
        if err != nil { panic(err) }

        for _, node := range nodes {
            resp, _, err := cli.NodeInspectWithRaw(context.Background(), node.ID)
            if err != nil { panic(err) }

            spec := resp.Spec

            spec.Annotations.Labels = make(map[string]string)
            for key, value := range instances[node.Description.Hostname] {
                fmt.Println("instance:", node.ID, "add label", key, "=", value)
                spec.Annotations.Labels[key] = value
            }

            err = cli.NodeUpdate(context.Background(), node.ID, resp.Version, spec)
            if err != nil { panic(err) }
        }
    }

}

