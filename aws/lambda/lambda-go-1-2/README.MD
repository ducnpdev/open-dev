# lambda simple:
- create simple lambda function
- read blog: https://viblo.asia/p/golang-aws-lambda-thong-qua-serverless-framework-phan-1-EoW4ob9xVml


aws lambda create-event-source-mapping --function-name lambda-go-dev-receive  --batch-size 10 --region us-east-1 \
--event-source-arn arn:aws:sqs:us-east-1:064038607558:sqs-demo


aws lambda list-event-source-mappings --function-name lambda-go-dev-receive --region us-east-1 \
--event-source-arn arn:aws:sqs:us-east-1:064038607558:sqs-demo  