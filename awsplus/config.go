package awsplus

import (
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/kinesis"
	"github.com/crowdmob/goamz/sqs"
)

type Config struct {
	AccessKey string
	SecretKey string
	Region    string
}

type AWSClient struct {
	kinesisconn *kinesis.Kinesis
	sqsconn     *sqs.SQS
}

// Client() returns a new client for accessing kinesis.

func (c *Config) Client() (interface{}, error) {
	var client AWSClient

	auth, err := c.AWSAuth()
	if err != nil {
		return nil, err
	}

	region, err := c.AWSRegion()
	if err != nil {
		return nil, err
	}

	client.kinesisconn = kinesis.New(auth, region)
	client.sqsconn = sqs.New(auth, region)

	log.Println("[INFO] Kinesis Client configured")

	return &client, nil
}

func (c *Config) AWSAuth() (aws.Auth, error) {
	exptdate := time.Now().Add(time.Hour)
	auth, err := aws.GetAuth(c.AccessKey, c.SecretKey, "", exptdate)
	if err == nil {
		c.AccessKey = auth.AccessKey
		c.SecretKey = auth.SecretKey
	}

	return auth, err
}

func (c *Config) IsValidRegion() bool {
	var regions = [11]string{"us-east-1", "us-west-2", "us-west-1", "eu-west-1",
		"eu-central-1", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1",
		"sa-east-1", "cn-north-1", "us-gov-west-1"}

	for _, valid := range regions {
		if c.Region == valid {
			return true
		}
	}
	return false
}

func (c *Config) AWSRegion() (aws.Region, error) {
	if c.Region != "" {
		if c.IsValidRegion() {
			return aws.Regions[c.Region], nil
		} else {
			return aws.Region{}, fmt.Errorf("Not a valid region: %s", c.Region)
		}
	}

	md, err := aws.GetMetaData("placement/availability-zone")
	if err != nil {
		return aws.Region{}, err
	}

	region := strings.TrimRightFunc(string(md), unicode.IsLetter)
	return aws.Regions[region], nil
}
