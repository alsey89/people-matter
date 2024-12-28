package s3conn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alsey89/people-matter/pkg/util"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Module struct {
	config *Config
	logger *zap.Logger
	scope  string
	s3     *s3.Client
	s3PS   *s3.PresignClient
}

type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Logger    *zap.Logger
}

type Config struct {
	Region        string
	ContribBucket string
	DeployBucket  string
}

const (
	DefaultRegion        = "ap-east-1"
	DefaultContribBucket = "curate.memorial"
	DefaultDeployBucket  = "contribution.curate.memorial"
)

//! Module ---------------------------------------------------------------

// Provides Module to the fx framework and registers Lifecycle hooks.
func InjectModule(scope string) fx.Option {
	return fx.Module(
		scope,
		fx.Provide(func(p Params) *Module {

			m := &Module{scope: scope}
			m.config = m.setupConfig(scope)
			m.logger = m.setupLogger(scope, p)
			m.s3 = m.setUpS3Client()
			m.s3PS = m.setUpS3PresignClient()

			return m
		}),
		fx.Invoke(func(m *Module, p Params) {
			p.Lifecycle.Append(
				fx.Hook{
					OnStart: m.onStart,
					OnStop:  m.onStop,
				},
			)
		}),
	)
}

// Instantiates a new Module without using the fx framework.
func NewS3Conn(scope string, logger *zap.Logger) *Module {
	m := &Module{scope: scope}
	m.logger = logger.Named("[" + scope + "]")
	m.config = m.setupConfig(scope)
	m.s3 = m.setUpS3Client()
	m.s3PS = m.setUpS3PresignClient()

	m.onStart(context.Background())

	return m
}

//! INTERNAL ---------------------------------------------------------------

func (m *Module) setupConfig(scope string) *Config {
	// Set default values
	viper.SetDefault(util.GetConfigPath(scope, "region"), DefaultRegion)
	viper.SetDefault(util.GetConfigPath(scope, "contrib_bucket"), DefaultContribBucket)
	viper.SetDefault(util.GetConfigPath(scope, "deploy_bucket"), DefaultDeployBucket)

	return &Config{
		Region:        viper.GetString(util.GetConfigPath(scope, "region")),
		ContribBucket: viper.GetString(util.GetConfigPath(scope, "contrib_bucket")),
		DeployBucket:  viper.GetString(util.GetConfigPath(scope, "deploy_bucket")),
	}
}

func (m *Module) setupLogger(scope string, p Params) *zap.Logger {
	logger := p.Logger.Named("[" + scope + "]")
	return logger
}

func (m *Module) setUpS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(m.config.Region))
	if err != nil {
		m.logger.Panic("Error loading AWS config", zap.Error(err))
	}

	s3Client := s3.NewFromConfig(cfg)
	return s3Client
}

func (m *Module) setUpS3PresignClient() *s3.PresignClient {
	return s3.NewPresignClient(m.s3)
}

func (m *Module) onStart(context.Context) error {
	m.logger.Info("Starting S3 connection.")

	//check s3 connection
	_, err := m.s3.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		m.logger.Panic("Error connecting to S3", zap.Error(err))
	}
	m.logger.Info("Connected to S3")

	if viper.GetString("global.log_level") == "DEBUG" || viper.GetString("global.log_level") == "debug" {
		m.logConfigurations()
	}

	return nil
}

func (m *Module) onStop(context.Context) error {
	m.logger.Info("Stopping S3 connection.")
	// No explicit close method for S3 client; clean up resources if necessary.
	return nil
}

func (m *Module) logConfigurations() {
	m.logger.Debug("----- S3 Configuration -----")
	m.logger.Debug("Region", zap.String("region", m.config.Region))
	m.logger.Debug("Bucket", zap.String("bucket", m.config.ContribBucket))
	m.logger.Debug("Bucket", zap.String("bucket", m.config.DeployBucket))
}

//! EXTERNAL ---------------------------------------------------------------

// Returns the AWS S3 client instance.
func (m *Module) GetS3Client() *s3.Client {
	return m.s3
}

// Returns the AWS S3 presigned URL client instance.
func (m *Module) GetPresignClient() *s3.PresignClient {
	return m.s3PS
}
func (m *Module) GetPresignedUploadURL(ctx context.Context, key string, exp time.Duration, metadata map[string]string) (*string, error) {

	req, err := m.s3PS.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:   &m.config.ContribBucket,
		Key:      &key,
		Metadata: metadata,
	}, s3.WithPresignExpires(exp))
	if err != nil {
		return nil, err
	}

	m.logger.Debug("Presigned URL", zap.String("url", req.URL))
	return &req.URL, nil
}
func (m *Module) GetPresignedDeployURL(ctx context.Context, key string, exp time.Duration) (*string, error) {

	req, err := m.s3PS.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &m.config.DeployBucket,
		Key:    &key,
	}, s3.WithPresignExpires(exp))
	if err != nil {
		return nil, err
	}

	m.logger.Debug("Presigned URL", zap.String("url", req.URL))
	return &req.URL, nil
}

func (m *Module) GetPresignedReadURL(ctx context.Context, s3URL string, exp time.Duration) (*string, error) {

	key := strings.TrimPrefix(s3URL, m.GetBucketPrefix(m.config.ContribBucket))

	req, err := m.s3PS.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &m.config.ContribBucket,
		Key:    &key,
	}, s3.WithPresignExpires(exp))
	if err != nil {
		return nil, err
	}

	return &req.URL, nil
}

// Returns the AWS prefix for urls in the bucket.
func (m *Module) GetBucketPrefix(bucket string) string {
	return "https://" + "s3." + m.config.Region + ".amazonaws.com/" + bucket + "/"
}

//Helpers-------------------------------------------------------------------

// Converts a URL ending with "_thumb" to "_og".
func (m *Module) ConvertObjectURLThumbToOG(thumbURL string) (*string, error) {
	if thumbURL == "" {
		return nil, fmt.Errorf("ConvertObjectURLThumbToOG: thumbURL is empty")
	}

	if !strings.HasSuffix(thumbURL, "_thumb") {
		return nil, fmt.Errorf("ConvertObjectURLThumbToOG: URL does not end with '_thumb', thumbURL: %s", thumbURL)
	}

	ogURL := strings.TrimSuffix(thumbURL, "_thumb") + "_og"
	return &ogURL, nil
}

// Converts a URL ending with "_og" to "_thumb".
func (m *Module) ConvertObjectURLOGToThumb(ogURL string) (*string, error) {
	if ogURL == "" {
		return nil, fmt.Errorf("ConvertObjectURLOGToThumb: ogURL is empty")
	}

	if !strings.HasSuffix(ogURL, "_og") {
		return nil, fmt.Errorf("ConvertObjectURLOGToThumb: URL does not end with '_og', ogURL: %s", ogURL)
	}

	thumbURL := strings.TrimSuffix(ogURL, "_og") + "_thumb"

	m.logger.Debug("Converted URL", zap.String("thumbURL", thumbURL))

	return &thumbURL, nil
}
