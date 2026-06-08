package s3client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	Bucket          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        *string
}

type Client struct {
	s3       *s3.Client
	bucket   string
	region   string
	endpoint *string
}

func New(cfg Config) (*Client, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("s3client: bucket is required")
	}
	if cfg.Region == "" {
		return nil, fmt.Errorf("s3client: region is required")
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("s3client: access key ID is required")
	}
	if cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("s3client: secret access key is required")
	}

	s3Client := s3.New(s3.Options{
		Region: cfg.Region,
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			),
		),
		BaseEndpoint: cfg.Endpoint,
		UsePathStyle: true,
	})

	return &Client{
		s3:       s3Client,
		bucket:   cfg.Bucket,
		region:   cfg.Region,
		endpoint: cfg.Endpoint,
	}, nil
}

// Upload puts an object into S3 at the given key.
// The caller is responsible for closing body if needed.
func (c *Client) Upload(ctx context.Context, key, contentType string, body io.Reader) error {
	_, err := c.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		Body:        body,
	})
	if err != nil {
		return fmt.Errorf("s3client: upload %q: %w", key, err)
	}
	return nil
}

// Delete removes an object from S3 at the given key.
func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("s3client: delete %q: %w", key, err)
	}
	return nil
}

// Get retrieves an object from S3 at the given key.
// The caller must close the returned ReadCloser when done.
func (c *Client) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	resp, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("s3client: get %q: %w", key, err)
	}
	return resp.Body, nil
}

// URL returns the public URL for the given key.
// Assumes the bucket is public and has no CloudFront or custom domain in front.
// If you use CloudFront, build the URL from your distribution domain instead.
func (c *Client) URL(key string) string {
	if c.endpoint != nil {
		return fmt.Sprintf("http://localhost:3909/api/browse/%s/%s", c.bucket, key)
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", c.bucket, c.region, key)
}

func (c *Client) GetKey(url string) string {
	if c.endpoint != nil {
		return strings.TrimPrefix(url, fmt.Sprintf("http://localhost:3909/api/browse/%s/", c.bucket))
	}
	return strings.TrimPrefix(url, fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", c.bucket, c.region))
}

// HashKey returns the SHA-256 hash of the given reader as a hex string.
func HashKey(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", fmt.Errorf("s3client: hash key: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
