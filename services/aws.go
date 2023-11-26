package services

import (
	"context"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"go.uber.org/zap"
	webappConfig "webapp/config"
	"webapp/logger"
)

type AWSService struct {
	AWSConfig *webappConfig.AWSConfig
}

func NewAWSService(configs *webappConfig.AWSConfig) *AWSService {
	return &AWSService{configs}
}

func (as AWSService) PublishSubmissionToSNS(submission string) {
	arn := as.AWSConfig.SNSArn
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())

	if err != nil {
		logger.Error("configuration error, ", zap.Any("error", err.Error()))
		return
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.PublishInput{
		Message:  &submission,
		TopicArn: &arn,
	}

	result, err := client.Publish(context.TODO(), input)
	if err != nil {
		logger.Error("Got an error publishing the message", zap.Any("error", err))
		return
	}

	logger.Info("Successfully sent submission to SNS", zap.Any("MessageID", *result.MessageId))
}
