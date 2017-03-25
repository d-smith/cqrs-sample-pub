aws ecs register-task-definition --cli-input-json file://$PWD/task-def.json

aws ecs run-task --task-definition 

aws ecs run-task --cluster DemoCluster --task-definition genevents

aws ecs describe-tasks --cluster DemoCluster --tasks arn:aws:ecs:us-west-1:444472610063:task/69b799b1-5c3c-4464-a93e-765f23189bd9
