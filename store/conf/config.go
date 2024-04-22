package conf

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

func DefaultConfig() *Config {
	return &Config{
		MongoDB: newDefaultMongoDB(),
	}
}

type Config struct {
	//[mongo]
	MongoDB *mongodb `json:"mongo" toml:"mongo"`
}

func newDefaultMongoDB() *mongodb {
	m := &mongodb{
		UserName:  "book",
		Password:  "123456",
		Database:  "book",
		AuthDB:    "",
		Endpoints: []string{"10.19.4.49:27017"},
	}
	return m
}

type mongodb struct {
	Endpoints []string `toml:"endpoints" env:"MONGO_ENDPOINTS" envSeparator:","`
	UserName  string   `toml:"username" env:"MONGO_USERNAME"`
	Password  string   `toml:"password" env:"MONGO_PASSWORD"`
	Database  string   `toml:"database" env:"MONGO_DATABASE"`
	AuthDB    string   `toml:"auth_db" env:"MONGO_AUTH_DB"`

	client *mongo.Client
	lock   sync.Mutex
}

func (m *mongodb) GetAuthDB() string {
	if m.AuthDB != "" {
		return m.AuthDB
	}

	return m.Database
}

func (m *mongodb) GetDB() (*mongo.Database, error) {
	conn, err := m.Client()
	if err != nil {
		return nil, err
	}
	return conn.Database(m.Database), nil
}

// 关闭数据库连接
func (m *mongodb) Close(ctx context.Context) error {
	if m.client == nil {
		return nil
	}

	return m.client.Disconnect(ctx)
}

// Client 获取一个全局的mongodb客户端连接
func (m *mongodb) Client() (*mongo.Client, error) {
	// 加载全局数据量单例
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.client == nil {
		conn, err := m.getClient()
		if err != nil {
			return nil, err
		}
		m.client = conn
	}

	return m.client, nil
}

func (m *mongodb) getClient() (*mongo.Client, error) {
	opts := options.Client()

	if m.UserName != "" && m.Password != "" {
		cred := options.Credential{
			AuthSource: m.GetAuthDB(),
		}

		cred.Username = m.UserName
		cred.Password = m.Password
		cred.PasswordSet = true
		opts.SetAuth(cred)
	}
	opts.SetHosts(m.Endpoints)
	opts.SetConnectTimeout(5 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("new mongodb client error, %s", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongodb server(%s) error, %s", m.Endpoints, err)
	}

	return client, nil
}
