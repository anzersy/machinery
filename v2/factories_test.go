package machinery_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/anzersy/machinery/v2"
	"github.com/anzersy/machinery/v2/backends"
	"github.com/anzersy/machinery/v2/brokers"
	"github.com/anzersy/machinery/v2/config"
	"github.com/stretchr/testify/assert"
)

func TestBrokerFactory(t *testing.T) {
	var cnf config.Config

	// 1) AMQP broker test

	cnf = config.Config{
		Broker:       "amqp://guest:guest@localhost:5672/",
		DefaultQueue: "machinery_tasks",
		AMQP: &config.AMQPConfig{
			Exchange:      "machinery_exchange",
			ExchangeType:  "direct",
			BindingKey:    "machinery_task",
			PrefetchCount: 1,
		},
	}

	actual, err := machinery.BrokerFactory(&cnf)
	if assert.NoError(t, err) {
		_, isAMQPBroker := actual.(*brokers.AMQPBroker)
		assert.True(
			t,
			isAMQPBroker,
			"Broker should be instance of *brokers.AMQPBroker",
		)
		expected := brokers.NewAMQPBroker(&cnf)
		assert.True(
			t,
			reflect.DeepEqual(actual, expected),
			fmt.Sprintf("conn = %v, want %v", actual, expected),
		)
	}

	// 2) Redis broker test

	// with password
	cnf = config.Config{
		Broker:       "redis://password@localhost:6379",
		DefaultQueue: "machinery_tasks",
	}

	actual, err = machinery.BrokerFactory(&cnf)
	if assert.NoError(t, err) {
		_, isRedisBroker := actual.(*brokers.RedisBroker)
		assert.True(
			t,
			isRedisBroker,
			"Broker should be instance of *brokers.RedisBroker",
		)
		expected := brokers.NewRedisBroker(&cnf, "localhost:6379", "password", "", 0)
		assert.True(
			t,
			reflect.DeepEqual(actual, expected),
			fmt.Sprintf("conn = %v, want %v", actual, expected),
		)
	}

	// without password
	cnf = config.Config{
		Broker:       "redis://localhost:6379",
		DefaultQueue: "machinery_tasks",
	}

	actual, err = machinery.BrokerFactory(&cnf)
	if assert.NoError(t, err) {
		_, isRedisBroker := actual.(*brokers.RedisBroker)
		assert.True(
			t,
			isRedisBroker,
			"Broker should be instance of *brokers.RedisBroker",
		)
		expected := brokers.NewRedisBroker(&cnf, "localhost:6379", "", "", 0)
		assert.True(
			t,
			reflect.DeepEqual(actual, expected),
			fmt.Sprintf("conn = %v, want %v", actual, expected),
		)
	}

	// using a socket file
	cnf = config.Config{
		Broker:       "redis+socket:///tmp/redis.sock",
		DefaultQueue: "machinery_tasks",
	}

	actual, err = machinery.BrokerFactory(&cnf)
	if assert.NoError(t, err) {
		_, isRedisBroker := actual.(*brokers.RedisBroker)
		assert.True(
			t,
			isRedisBroker,
			"Broker should be instance of *brokers.RedisBroker",
		)
		expected := brokers.NewRedisBroker(&cnf, "", "", "/tmp/redis.sock", 0)
		assert.True(
			t,
			reflect.DeepEqual(actual, expected),
			fmt.Sprintf("conn = %v, want %v", actual, expected),
		)
	}

	// 3) AWS SQS
	cnf = config.Config{
		Broker:       "https://sqs.us-east-2.amazonaws.com/123456789012",
		DefaultQueue: "machinery_tasks",
	}

	actual, err = machinery.BrokerFactory(&cnf)
	if assert.NoError(t, err) {
		_, isAWSSQSBroker := actual.(*brokers.AWSSQSBroker)
		assert.True(
			t,
			isAWSSQSBroker,
			"Broker should be instance of *brokers.AWSSQSBroker",
		)
	}
}

func TestBrokerFactoryError(t *testing.T) {
	cnf := config.Config{
		Broker: "BOGUS",
	}

	conn, err := machinery.BrokerFactory(&cnf)
	if assert.Error(t, err) {
		assert.Nil(t, conn)
		assert.Equal(t, "Factory failed with broker URL: BOGUS", err.Error())
	}
}

func TestBackendFactory(t *testing.T) {
	var cnf config.Config

	// 1) AMQP backend test

	cnf = config.Config{ResultBackend: "amqp://guest:guest@localhost:5672/"}

	actual, err := machinery.BackendFactory(&cnf)
	if assert.NoError(t, err) {
		expected := backends.NewAMQPBackend(&cnf)
		assert.True(
			t,
			reflect.DeepEqual(actual, expected),
			fmt.Sprintf("conn = %v, want %v", actual, expected),
		)
	}

	// 2) Memcache backend test

	cnf = config.Config{
		ResultBackend: "memcache://10.0.0.1:11211,10.0.0.2:11211",
	}

	actual, err = machinery.BackendFactory(&cnf)
	if assert.NoError(t, err) {
		servers := []string{"10.0.0.1:11211", "10.0.0.2:11211"}
		expected := backends.NewMemcacheBackend(&cnf, servers)
		assert.True(
			t,
			reflect.DeepEqual(actual, expected),
			fmt.Sprintf("conn = %v, want %v", actual, expected),
		)
	}

	// 2) Redis backend test

	// with password
	cnf = config.Config{
		ResultBackend: "redis://password@localhost:6379",
	}

	actual, err = machinery.BackendFactory(&cnf)
	if assert.NoError(t, err) {
		expected := backends.NewRedisBackend(&cnf, "localhost:6379", "password", "", 0)
		assert.True(
			t,
			reflect.DeepEqual(actual, expected),
			fmt.Sprintf("conn = %v, want %v", actual, expected),
		)
	}

	// without password
	cnf = config.Config{
		ResultBackend: "redis://localhost:6379",
	}

	actual, err = machinery.BackendFactory(&cnf)
	if assert.NoError(t, err) {
		expected := backends.NewRedisBackend(&cnf, "localhost:6379", "", "", 0)
		assert.True(
			t,
			reflect.DeepEqual(actual, expected),
			fmt.Sprintf("conn = %v, want %v", actual, expected),
		)
	}

	// using a socket file
	cnf = config.Config{
		ResultBackend: "redis+socket:///tmp/redis.sock",
	}

	actual, err = machinery.BackendFactory(&cnf)
	if assert.NoError(t, err) {
		expected := backends.NewRedisBackend(&cnf, "", "", "/tmp/redis.sock", 0)
		assert.True(
			t,
			reflect.DeepEqual(actual, expected),
			fmt.Sprintf("conn = %v, want %v", actual, expected),
		)
	}

	// 4) MongoDB backend test

	cnf = config.Config{
		ResultBackend: "mongodb://localhost:27017/tasks",
	}

	actual, err = machinery.BackendFactory(&cnf)
	if assert.NoError(t, err) {
		expected := backends.NewMongodbBackend(&cnf)
		assert.True(
			t,
			reflect.DeepEqual(actual, expected),
			fmt.Sprintf("conn = %v, want %v", actual, expected),
		)
	}
}

func TestBackendFactoryError(t *testing.T) {
	cnf := config.Config{
		ResultBackend: "BOGUS",
	}

	conn, err := machinery.BackendFactory(&cnf)
	if assert.Error(t, err) {
		assert.Nil(t, conn)
		assert.Equal(t, "Factory failed with result backend: BOGUS", err.Error())
	}

	if conn != nil {
		t.Errorf("conn = %v, should be nil", conn)
	}
}

func TestParseRedisURL(t *testing.T) {
	var host, pwd, url string
	var db int
	var err error

	url = "non_redis://127.0.0.1:5672"
	_, _, _, err = machinery.ParseRedisURL(url)
	assert.Error(t, err, "invalid redis scheme")

	url = "redis:/"
	_, _, _, err = machinery.ParseRedisURL(url)
	assert.Error(t, err, "invalid redis url scheme")

	url = "redis://127.0.0.1:5672"
	host, pwd, db, err = machinery.ParseRedisURL(url)
	if assert.NoError(t, err) {
		assert.Equal(t, "127.0.0.1:5672", host)
		assert.Equal(t, "", pwd)
		assert.Equal(t, 0, db)
	}

	url = "redis://pwd@127.0.0.1:5672"
	host, pwd, db, _ = machinery.ParseRedisURL(url)
	if assert.NoError(t, err) {
		assert.Equal(t, "127.0.0.1:5672", host)
		assert.Equal(t, "pwd", pwd)
		assert.Equal(t, 0, db)
	}

	url = "redis://pwd@127.0.0.1:5672/2"
	host, pwd, db, err = machinery.ParseRedisURL(url)
	if assert.NoError(t, err) {
		assert.Equal(t, "127.0.0.1:5672", host)
		assert.Equal(t, "pwd", pwd)
		assert.Equal(t, 2, db)
	}
}

func TestParseRedisSocketURL(t *testing.T) {
	var path, pwd, url string
	var db int
	var err error

	url = "non_redissock:///tmp/redis.sock"
	_, _, _, err = machinery.ParseRedisSocketURL(url)
	assert.Error(t, err, "invalid redis scheme")

	url = "redis+socket:/"
	_, _, _, err = machinery.ParseRedisSocketURL(url)
	assert.Error(t, err, "invalid redis url scheme")

	url = "redis+socket:///tmp/redis.sock"
	path, pwd, db, err = machinery.ParseRedisSocketURL(url)
	if assert.NoError(t, err) {
		assert.Equal(t, "/tmp/redis.sock", path)
		assert.Equal(t, "", pwd)
		assert.Equal(t, 0, db)
	}

	url = "redis+socket://pwd@/tmp/redis.sock"
	path, pwd, db, _ = machinery.ParseRedisSocketURL(url)
	if assert.NoError(t, err) {
		assert.Equal(t, "/tmp/redis.sock", path)
		assert.Equal(t, "pwd", pwd)
		assert.Equal(t, 0, db)
	}

	url = "redis+socket://pwd@/tmp/redis.sock:/2"
	path, pwd, db, err = machinery.ParseRedisSocketURL(url)
	if assert.NoError(t, err) {
		assert.Equal(t, "/tmp/redis.sock", path)
		assert.Equal(t, "pwd", pwd)
		assert.Equal(t, 2, db)
	}
}
