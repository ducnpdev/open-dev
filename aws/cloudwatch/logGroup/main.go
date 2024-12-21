package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func main() {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-southeast-1"),
		config.WithSharedConfigProfile("aws-uat"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create CloudWatch Logs client
	client := cloudwatchlogs.NewFromConfig(cfg)

	// Get list of log groups
	ListLogGroups(client)
}

func ListLogGroups(client *cloudwatchlogs.Client) {
	input := &cloudwatchlogs.DescribeLogGroupsInput{
		Limit: aws.Int32(10), // Optional: Limit the number of log groups
	}

	for {
		// Call DescribeLogGroups
		output, err := client.DescribeLogGroups(context.TODO(), input)
		if err != nil {
			log.Fatalf("failed to describe log groups: %v", err)
		}

		// Print the log group names
		for _, logGroup := range output.LogGroups {
			fmt.Printf("Log Group: %s\n", aws.ToString(logGroup.LogGroupName))
			err = checkOldLogStreams(client, aws.ToString(logGroup.LogGroupName))
			if err != nil {
				log.Fatalf("failed to check log streams: %v", err)
			}
		}

		// Check if there are more results
		if output.NextToken == nil {
			break
		}
		// Update input with the next token
		input.NextToken = output.NextToken
	}
}

func ListLogStreams(client *cloudwatchlogs.Client, logGroupName string) error {
	input := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroupName),
		OrderBy:      types.OrderByLastEventTime, // Sort by the most recent event time
		Descending:   aws.Bool(true),             // Show the latest log streams first
		Limit:        aws.Int32(10),              // Limit the results to 10 log streams
	}

	for {
		// Call DescribeLogStreams API
		output, err := client.DescribeLogStreams(context.TODO(), input)
		if err != nil {
			return fmt.Errorf("error describing log streams: %v", err)
		}

		// Print log stream names
		for _, logStream := range output.LogStreams {
			fmt.Printf("Log Stream: %s\n", aws.ToString(logStream.LogStreamName))
			t := time.UnixMilli(int64(aws.ToInt64(logStream.LastEventTimestamp)))
			lastEventTime := t.Local().Format("2006-01-02")
			fmt.Println("Thời gian (Local):", lastEventTime)
			time.Sleep(time.Second)

		}

		// Check if there are more results
		if output.NextToken == nil {
			break
		}

		// Update the input with the next token for pagination
		input.NextToken = output.NextToken
	}

	return nil
}

func checkOldLogStreams(client *cloudwatchlogs.Client, logGroupName string) error {
	// Define the cutoff time (1 month ago)
	// oneMonthAgo := time.Now().AddDate(0, -1, 0).UnixMilli()

	// Lấy thời gian 15 ngày trước
	oneMonthAgo := time.Now().Add(-15 * 24 * time.Hour).UnixMilli()

	input := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroupName),
		OrderBy:      types.OrderByLastEventTime,
		Descending:   aws.Bool(true),
	}

	for {
		// Call DescribeLogStreams API
		output, err := client.DescribeLogStreams(context.TODO(), input)
		if err != nil {
			return fmt.Errorf("error describing log streams: %v", err)
		}

		// Filter and print log streams with LastEventTimestamp older than 1 month
		for _, logStream := range output.LogStreams {
			if logStream.LastEventTimestamp != nil && *logStream.LastEventTimestamp < oneMonthAgo {
				fmt.Printf("Log Stream: %s, Last Event Time: %v\n",
					aws.ToString(logStream.LogStreamName),
					time.UnixMilli(*logStream.LastEventTimestamp).Format("2006-01-02 15:04:05"))
				DeleteLogStream(client, logGroupName, aws.ToString(logStream.LogStreamName))

			}
		}

		// Check if there are more results
		if output.NextToken == nil {
			break
		}

		// Update the input with the next token for pagination
		input.NextToken = output.NextToken
	}

	return nil
}

func DeleteLogStream(client *cloudwatchlogs.Client, logGroupName, logStreamName string) error {

	// Prepare the input for DeleteLogStream API
	input := &cloudwatchlogs.DeleteLogStreamInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
	}

	// Call the DeleteLogStream API
	_, err := client.DeleteLogStream(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("error deleting log stream: %v", err)
	}

	fmt.Println("delete ok", logGroupName, logStreamName)
	// time.Sleep(time.Second)

	return nil
}
