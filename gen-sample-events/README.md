# Note on Sample Event Generation

When hosting an event store in AWS, it can be useful for testing and 
demonstration purposes to create some sample events. In a typical 
AWS set up, the RDS instance is in a private subnet with no public
network access to it. One way to access the database directly is to 
launch two hosts - a bastion host in the public subnet, a private host
in the private subnet, allow ingress from the public host to the 
private, configure security groups to allow the private host access to
the RDS instance, and so on.

When the RDS is part of an application built using the AWS ECS,
assuming the ECS cluster instances have database access, running one
off tasks against the database is as easy as... well, running an
ECS task.

This directory contains a task definition skeleton to allow events to be
seeded in an event store that is part of an ECS application.

To execute the task, first customize the particulars in 
`task-def-skeleton.json` to create `task-def.json`. Once this has
been done, register the task definition:

<pre>
aws ecs register-task-definition --cli-input-json file://$PWD/task-def.json
</pre>

Once the task has been defined in ECS, you can then run it:

<pre>
aws ecs run-task --cluster DemoCluster --task-definition genevents
</pre>

You can status the task using the task arn returned by the run-task
command.

<pre>
aws ecs describe-tasks --cluster DemoCluster --tasks arn:aws:ecs:us-west-1:nnnnnnnn:task/69b799b1-5c3c-4464-a93e-765f23189bd9
</pre>
