package rabbitmq

import (
	"fmt"
	"net/url"

	"github.com/ZaiiiRan/backend_labs/order-service/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqClient struct {
	cfg  *config.RabbitMqSettings
	conn *amqp.Connection
}

func NewRabbitMqClient(cfg *config.RabbitMqSettings) (*RabbitMqClient, error) {
	c := &RabbitMqClient{cfg: cfg}
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *RabbitMqClient) Channel() (*amqp.Channel, error) {
	if c.conn == nil || c.conn.IsClosed() {
		if err := c.connect(); err != nil {
			return nil, fmt.Errorf("rabbitmq reconnect: %w", err)
		}
	}
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}
	return ch, nil
}

func (c *RabbitMqClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *RabbitMqClient) connect() error {
	vhost := url.PathEscape(c.cfg.VHost)
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		c.cfg.User,
		c.cfg.Password,
		c.cfg.Host,
		c.cfg.Port,
		vhost,
	)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return fmt.Errorf("amqp dial: %w", err)
	}
	c.conn = conn
	return nil
}
