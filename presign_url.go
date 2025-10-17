package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

const (
	expirationURL = 5 * time.Minute
)

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s3Client)
	presignURL, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", err
	}
	return presignURL.URL, nil
}

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
		return video, errors.New("missing video URL")
	}

	videoURL := strings.Split(*video.VideoURL, ",")
	if len(videoURL) != 2 {
		return video, errors.New("invalid video URL")
	}
	bucket := videoURL[0]
	key := videoURL[1]

	presignedURL, err := generatePresignedURL(cfg.s3Client, bucket, key, expirationURL)
	if err != nil {
		return video, err
	}

	video.VideoURL = &presignedURL
	return video, nil
}
