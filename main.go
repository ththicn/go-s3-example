package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	bucket = "YOUR_AWS_BUCKET"
	region = "us-east-2"
)

func main() {
	profile := flag.String("profile", "", "AWS profile name")
	ctx := context.Background()
	config, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithSharedConfigProfile(*profile),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := s3.NewFromConfig(config)

	lists, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String("source"),
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range lists.Contents {
		log.Println(*v.Key)
		destKey := strings.Split(*v.Key, "/")
		if len(destKey) == 0 {
			log.Fatal("no contents")
		}
		output, err := client.CopyObject(ctx, &s3.CopyObjectInput{
			CopySource: aws.String(bucket + "/" + *v.Key),
			Bucket:     aws.String(bucket),
			Key:        aws.String("dest/" + destKey[len(destKey)-1]),
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", *output.CopyObjectResult)
	}
}
