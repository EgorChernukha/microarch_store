package amqp

import (
	stderrors "errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Config contains parameters for AMQP connection
type Config struct {
	User           string
	Password       string
	Host           string
	ConnectTimeout time.Duration // 0 means default timeout (60 seconds)
}

type Connection interface {
	Start() error
	Stop() error
	AddChannel(channel Channel)
}

type Channel interface {
	Connect(conn *amqp.Connection) error
}

func NewAMQPConnection(cfg *Config) Connection {
	return &connection{cfg: cfg}
}

var (
	errNilAMQPConnection    = stderrors.New("amqp connection is empty")
	errClosedAMQPConnection = stderrors.New("amqp connection is closed")
)

type connection struct {
	cfg      *Config
	amqpConn *amqp.Connection
	channels []Channel
}

func (c *connection) Start() error {
	url := fmt.Sprintf("amqp://%s:%s@%s/", c.cfg.User, c.cfg.Password, c.cfg.Host)

	err := backoff.Retry(func() error {
		connection, cErr := amqp.Dial(url)
		c.amqpConn = connection
		return errors.Wrap(cErr, "failed to connect to amqp")
	}, newBackOff(c.cfg.ConnectTimeout))

	if err == nil {
		if err = c.validateConnection(c.amqpConn); err != nil {
			return err
		}

		for _, channel := range c.channels {
			if err = channel.Connect(c.amqpConn); err != nil {
				return err
			}
		}

		connErrorChan := c.amqpConn.NotifyClose(make(chan *amqp.Error))
		go c.processConnectErrors(connErrorChan)
	}
	return err
}

func (c *connection) Stop() error {
	return c.amqpConn.Close()
}

func (c *connection) AddChannel(channel Channel) {
	c.channels = append(c.channels, channel)
}

func (c *connection) validateConnection(conn *amqp.Connection) error {
	if conn == nil {
		return errors.WithStack(errNilAMQPConnection)
	}
	if conn.IsClosed() {
		return errors.WithStack(errClosedAMQPConnection)
	}
	return nil
}

// channel will be closed then the connection is closed so this function will exit, no need for custom graceful shutdown
func (c *connection) processConnectErrors(ch chan *amqp.Error) {
	err := <-ch
	if err == nil {
		return
	}

	logrus.Error(err, "AMQP connection error, trying to reconnect")
	for {
		err := c.Start()
		if err == nil {
			logrus.Info("AMQP connection restored")
			break
		} else {
			logrus.Error(err, "failed to reconnect to AMQP")
		}
	}
}

func newBackOff(timeout time.Duration) backoff.BackOff {
	exponentialBackOff := backoff.NewExponentialBackOff()
	const defaultTimeout = 60 * time.Second
	if timeout != 0 {
		exponentialBackOff.MaxElapsedTime = timeout
	} else {
		exponentialBackOff.MaxElapsedTime = defaultTimeout
	}
	exponentialBackOff.MaxInterval = 5 * time.Second
	return exponentialBackOff
}
