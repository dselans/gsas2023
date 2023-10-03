package deps

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/InVisionApp/go-health"
	gllogrus "github.com/InVisionApp/go-logger/shims/logrus"
	"github.com/batchcorp/rabbit"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/streamdal/natty"

	"github.com/dselans/welcome-svc/backends/cache"
	"github.com/dselans/welcome-svc/config"
	"github.com/dselans/welcome-svc/services/isb"
)

const (
	DefaultHealthCheckIntervalSecs = 1
)

type customCheck struct{}

type Dependencies struct {
	// Backends
	//ISBDedicatedBackend rabbit.IRabbit // Commented out for GSAS2023
	ISBSharedBackend rabbit.IRabbit
	CacheBackend     cache.ICache
	NATSBackend      natty.INatty

	// Services
	ISBService isb.IISB

	// Commented out for GSAS2023 workshop
	// HSBBackend kafka.IKafka
	// HSBService hsb.IHSB
	//
	//HSBChan chan *events.Manifest

	Health         health.IHealth
	DefaultContext context.Context
}

func New(cfg *config.Config) (*Dependencies, error) {
	gohealth := health.New()
	gohealth.Logger = gllogrus.New(nil)

	d := &Dependencies{
		Health:         gohealth,
		DefaultContext: context.Background(),
		//HSBChan:        make(chan *events.Manifest, 0),
	}

	if err := d.setupHealthChecks(); err != nil {
		return nil, errors.Wrap(err, "unable to setup health check(s)")
	}

	if err := d.Health.Start(); err != nil {
		return nil, errors.Wrap(err, "unable to start health runner")
	}

	if err := d.setupBackends(cfg); err != nil {
		return nil, errors.Wrap(err, "unable to setup backends")
	}

	if err := d.setupServices(cfg); err != nil {
		return nil, errors.Wrap(err, "unable to setup services")
	}

	return d, nil
}

func (d *Dependencies) setupHealthChecks() error {
	cc := &customCheck{}

	err := d.Health.AddChecks([]*health.Config{
		{
			Name:     "health-check",
			Checker:  cc,
			Interval: time.Duration(DefaultHealthCheckIntervalSecs) * time.Second,
			Fatal:    true,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (d *Dependencies) setupBackends(cfg *config.Config) error {
	// CacheBackend k/v store
	cb, err := cache.New()
	if err != nil {
		return errors.Wrap(err, "unable to create new cache instance")
	}

	d.CacheBackend = cb

	// NATS backend
	n, err := natty.New(&natty.Config{
		NatsURL:           cfg.NATSURL,
		UseTLS:            cfg.NATSUseTLS,
		TLSCACertFile:     cfg.NATSTLSCaFile,
		TLSClientCertFile: cfg.NATSTLSCertFile,
		TLSClientKeyFile:  cfg.NATSTLSKeyFile,
	})

	if err != nil {
		return errors.Wrap(err, "unable to create new nats backend")
	}

	d.NATSBackend = n

	// Events rabbitmq backend (commented out for GSAS2023)
	//isbDedicatedBackend, err := rabbit.New(&rabbit.Options{
	//	URLs:      cfg.ISBDedicatedURLs,
	//	Mode:      0,
	//	QueueName: cfg.ISBDedicatedQueueName,
	//	Bindings: []rabbit.Binding{
	//		{
	//			ExchangeName:    cfg.ISBDedicatedExchangeName,
	//			ExchangeType:    amqp.ExchangeTopic,
	//			ExchangeDeclare: cfg.ISBDedicatedExchangeDeclare,
	//			BindingKeys:     cfg.ISBDedicatedBindingKeys,
	//		},
	//	},
	//	RetryReconnectSec: rabbit.DefaultRetryReconnectSec,
	//	QueueDurable:      cfg.ISBDedicatedQueueDurable,
	//	QueueExclusive:    cfg.ISBDedicatedQueueExclusive,
	//	QueueAutoDelete:   cfg.ISBDedicatedQueueAutoDelete,
	//	QueueDeclare:      cfg.ISBDedicatedQueueDeclare,
	//	AutoAck:           cfg.ISBDedicatedAutoAck,
	//	AppID:             cfg.ServiceName,
	//	UseTLS:            cfg.ISBDedicatedUseTLS,
	//	SkipVerifyTLS:     cfg.ISBDedicatedSkipVerifyTLS,
	//})
	//if err != nil {
	//	return errors.Wrap(err, "unable to create new dedicated rabbit backend")
	//}
	//
	//d.ISBDedicatedBackend = isbDedicatedBackend

	// Shared backend
	isbSharedBackend, err := rabbit.New(&rabbit.Options{
		URLs:      cfg.ISBSharedURLs,
		Mode:      0,
		QueueName: cfg.ISBSharedQueueName,
		Bindings: []rabbit.Binding{
			{
				ExchangeName:    cfg.ISBSharedExchangeName,
				ExchangeType:    amqp.ExchangeTopic,
				ExchangeDeclare: cfg.ISBSharedExchangeDeclare,
				ExchangeDurable: cfg.ISBSharedExchangeDurable,
				BindingKeys:     cfg.ISBSharedBindingKeys,
			},
		},
		RetryReconnectSec: rabbit.DefaultRetryReconnectSec,
		QueueDurable:      cfg.ISBSharedQueueDurable,
		QueueExclusive:    cfg.ISBSharedQueueExclusive,
		QueueAutoDelete:   cfg.ISBSharedQueueAutoDelete,
		QueueDeclare:      cfg.ISBSharedQueueDeclare,
		AutoAck:           cfg.ISBSharedAutoAck,
		UseTLS:            cfg.ISBSharedUseTLS,
		SkipVerifyTLS:     cfg.ISBSharedSkipVerifyTLS,
	})
	if err != nil {
		return errors.Wrap(err, "unable to create new shared rabbit backend")
	}

	d.ISBSharedBackend = isbSharedBackend

	//if cfg.HSBUseTLS {
	//	logrus.Debug("using TLS for HSB")
	//}
	//
	//hsbBackend, err := kafka.New(
	//	&kafka.Options{
	//		Topic:     cfg.HSBTopicName,
	//		Brokers:   cfg.HSBBrokerURLs,
	//		Timeout:   cfg.HSBConnectTimeout,
	//		BatchSize: cfg.HSBBatchSize,
	//		UseTLS:    cfg.HSBUseTLS,
	//	},
	//	d.DefaultContext,
	//)
	//if err != nil {
	//	return errors.Wrap(err, "unable to create new kafka instance")
	//}
	//
	//d.HSBBackend = hsbBackend

	return nil
}

func (d *Dependencies) setupServices(cfg *config.Config) error {
	isbService, err := isb.New(&isb.Config{
		NATS:  d.NATSBackend,
		Cache: d.CacheBackend,
		RabbitMap: map[string]*isb.RabbitConfig{
			// Commented out for GSAS2023
			//
			// This is another type of Rabbit consumer - messages sent here will
			// be consumed by _every_ instance of this service. Think of it as a
			// broadcast or pubsub-like behavior.
			//
			// You would use this when you want _everyone_ to do something such
			// as maybe performing a global refresh or banning an IP.
			//
			//"dedicated": {
			//	RabbitInstance: d.ISBDedicatedBackend,
			//	NumConsumers:   cfg.ISBDedicatedNumConsumers,
			//	Func:           "DedicatedConsumeFunc",
			//},

			// The shared consumer behaves you would expect - a message sent here
			// will arrive exactly on ONE instance of this service.
			"shared": {
				RabbitInstance: d.ISBSharedBackend,
				NumConsumers:   cfg.ISBSharedNumConsumers,
				Func:           "SharedConsumeFunc",
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "unable to setup event")
	}

	d.ISBService = isbService

	//hsbService, err := hsb.New(&hsb.Config{
	//	Kafka:         d.HSBBackend,
	//	NumPublishers: cfg.HSBNumPublishers,
	//	HSBChan:       d.HSBChan,
	//})
	//if err != nil {
	//	return errors.Wrap(err, "unable to setup hsb")
	//}
	//
	//d.HSBService = hsbService

	return nil
}

func createTLSConfig(caCert, clientCert, clientKey string) (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	if err != nil {
		return nil, errors.Wrap(err, "unable to load cert + key")
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(caCert))

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}, nil
}

// Satisfy the go-health.ICheckable interface
func (c *customCheck) Status() (interface{}, error) {
	if false {
		return nil, errors.New("something major just broke")
	}

	// You can return additional information pertaining to the check as long
	// as it can be JSON marshalled
	return map[string]int{}, nil
}

func (d *Dependencies) PreCreateBuckets(ctx context.Context, cfg *config.Config) error {
	buckets := map[string]time.Duration{
		//BucketNameHere:      0,
	}

	for bucketName, ttl := range buckets {
		if err := d.NATSBackend.CreateBucket(ctx, bucketName, ttl, cfg.NATSNumBucketReplicas); err != nil {
			if err == nats.ErrStreamNameAlreadyInUse {
				continue
			}

			return fmt.Errorf("unable to create bucket '%s': %s", bucketName, err)
		}
	}

	return nil
}
