package s3storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/qor5/x/v3/oss"
)

// Client S3 storage
type Client struct {
	S3     *s3.S3
	Config Config
}

// New initialize S3 storage
func New(config Config) *Client {
	if config.ACL == "" {
		config.ACL = "public-read"
	}

	client := &Client{Config: config}

	// Создаем AWS конфигурацию
	awsConfig := &aws.Config{
		Region:           aws.String(config.Region),
		S3ForcePathStyle: aws.Bool(config.S3ForcePathStyle),
	}

	// Устанавливаем endpoint если указан
	if config.S3Endpoint != "" {
		awsConfig.Endpoint = aws.String(config.S3Endpoint)
	}

	// Используем статические credentials
	awsConfig.Credentials = credentials.NewStaticCredentials(
		config.AccessID,
		config.AccessKey,
		config.SessionToken,
	)
	sess := session.Must(session.NewSession(awsConfig))
	client.S3 = s3.New(sess)

	return client
}

// Get receive file with given path
func (client Client) Get(ctx context.Context, path string) (file *os.File, err error) {
	readCloser, err := client.GetStream(ctx, path)

	ext := filepath.Ext(path)
	pattern := fmt.Sprintf("s3*%s", ext)

	if err == nil {
		if file, err = os.CreateTemp("/tmp", pattern); err == nil {
			defer readCloser.Close()
			_, err = io.Copy(file, readCloser)
			file.Seek(0, 0)
		}
	}

	return file, err
}

// GetStream get file as stream
func (client Client) GetStream(ctx context.Context, path string) (io.ReadCloser, error) {
	getResponse, err := client.S3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(client.Config.Bucket),
		Key:    aws.String(client.ToS3Key(path)),
	})
	if err != nil {
		return nil, err
	}

	return getResponse.Body, err
}

// Put store a reader into given path
func (client Client) Put(ctx context.Context, urlPath string, reader io.Reader) (*oss.Object, error) {
	if seeker, ok := reader.(io.ReadSeeker); ok {
		seeker.Seek(0, 0)
	}

	urlPath = client.ToS3Key(urlPath)
	buffer, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	fileType := mime.TypeByExtension(path.Ext(urlPath))
	if fileType == "" {
		fileType = http.DetectContentType(buffer)
	}

	params := &s3.PutObjectInput{
		Bucket:        aws.String(client.Config.Bucket), // required
		Key:           aws.String(urlPath),              // required
		ACL:           aws.String(client.Config.ACL),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(int64(len(buffer))),
		ContentType:   aws.String(fileType),
	}
	if client.Config.CacheControl != "" {
		params.CacheControl = aws.String(client.Config.CacheControl)
	}

	_, err = client.S3.PutObjectWithContext(ctx, params)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &oss.Object{
		Path:             urlPath,
		Name:             filepath.Base(urlPath),
		LastModified:     &now,
		StorageInterface: client,
	}, nil
}

// Delete delete file
func (client Client) Delete(ctx context.Context, path string) error {
	_, err := client.S3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(client.Config.Bucket),
		Key:    aws.String(client.ToS3Key(path)),
	})
	return err
}

// DeleteObjects delete files in bulk
func (client Client) DeleteObjects(ctx context.Context, paths []string) (err error) {
	var objs []*s3.ObjectIdentifier
	for _, v := range paths {
		obj := &s3.ObjectIdentifier{
			Key: aws.String(strings.TrimPrefix(client.ToS3Key(v), "/")),
		}
		objs = append(objs, obj)
	}
	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(client.Config.Bucket),
		Delete: &s3.Delete{
			Objects: objs,
		},
	}

	_, err = client.S3.DeleteObjectsWithContext(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

// List list all objects under current path
func (client Client) List(ctx context.Context, path string) ([]*oss.Object, error) {
	var objects []*oss.Object
	var prefix string

	if path != "" {
		prefix = strings.Trim(path, "/") + "/"
	}

	listObjectsResponse, err := client.S3.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(client.Config.Bucket),
		Prefix: aws.String(prefix),
	})

	if err == nil {
		for _, content := range listObjectsResponse.Contents {
			objects = append(objects, &oss.Object{
				Path:             "/" + client.ToS3Key(*content.Key),
				Name:             filepath.Base(*content.Key),
				LastModified:     content.LastModified,
				StorageInterface: client,
			})
		}
	}

	return objects, err
}

// GetEndpoint get endpoint, FileSystem's endpoint is /
func (client Client) GetEndpoint(ctx context.Context) string {
	if client.Config.Endpoint != "" {
		return client.Config.Endpoint
	}

	endpoint := client.getS3Endpoint(ctx)
	for _, prefix := range []string{"https://", "http://"} {
		endpoint = strings.TrimPrefix(endpoint, prefix)
	}

	return client.Config.Bucket + "." + endpoint
}

var urlRegexp = regexp.MustCompile(`(https?:)?//((\w+).)+(\w+)/`)

// ToS3Key process path to s3 key
func (client Client) ToS3Key(urlPath string) string {
	if urlRegexp.MatchString(urlPath) {
		if u, err := url.Parse(urlPath); err == nil {
			if client.Config.S3ForcePathStyle { // First part of path will be bucket name
				return strings.TrimPrefix(strings.TrimPrefix(u.Path, "/"+client.Config.Bucket), "/")
			}
			return strings.TrimPrefix(u.Path, "/")
		}
	}

	if client.Config.S3ForcePathStyle { // First part of path will be bucket name
		return strings.TrimPrefix(urlPath, "/"+client.Config.Bucket+"/")
	}
	return strings.TrimPrefix(urlPath, "/")
}

// GetURL get public accessible URL
func (client Client) GetURL(ctx context.Context, path string) (url string, err error) {
	if client.getS3Endpoint(ctx) == "" {
		if client.Config.ACL == "private" || client.Config.ACL == "authenticated-read" {
			req, _ := client.S3.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(client.Config.Bucket),
				Key:    aws.String(client.ToS3Key(path)),
			})

			url, err = req.Presign(1 * time.Hour)
			if err != nil {
				return "", err
			}
			return url, nil
		}
	}

	return path, nil
}

// Copy copy s3 file from "from" to "to"
func (client Client) Copy(ctx context.Context, from, to string) (err error) {
	_, err = client.S3.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(client.Config.Bucket),
		CopySource: aws.String(from),
		Key:        aws.String(to),
	})
	return err
}

func (client Client) getS3Endpoint(ctx context.Context) string {
	if client.Config.S3Endpoint != "" {
		return client.Config.S3Endpoint
	}

	// В AWS SDK v1 endpoint обычно определяется автоматически на основе региона
	// Для кастомных endpoint нужно указывать явно в конфиге
	return ""
}

// Дополнительные методы для совместимости

// PutWithContext алиас для Put с контекстом
func (client Client) PutWithContext(ctx context.Context, urlPath string, reader io.Reader) (*oss.Object, error) {
	return client.Put(ctx, urlPath, reader)
}

// GetWithContext алиас для Get с контекстом
func (client Client) GetWithContext(ctx context.Context, path string) (*os.File, error) {
	return client.Get(ctx, path)
}

// DeleteWithContext алиас для Delete с контекстом
func (client Client) DeleteWithContext(ctx context.Context, path string) error {
	return client.Delete(ctx, path)
}

func (client Client) SelectObjectContentWithContext(
	ctx context.Context,
	in *s3.SelectObjectContentInput,
) (*s3.SelectObjectContentOutput, error) {
	return client.S3.SelectObjectContentWithContext(ctx, in)
}
