package redis

import (
	"fmt"
	"github.com/tornadoyi/viking/redis"
	"gopkg.in/yaml.v2"
	"testing"
)

const (
	network = "tcp"
	address = ""
	name = "test"
)

func skip() bool{ return len(network) == 0 || len(address) == 0 }


func TestRedisPoolCreator(t *testing.T) {
	if(skip()) {t.Skip("Set redis network and port address activate redis test")}

	type RedisConfig struct{
		Name				string					`yaml:"name"`
		Network				string					`yaml:"network"`
		Host				string					`yaml:"host"`
		PoolConfig			*redis.PoolConfig		`yaml:"pool"`
		DialConfig			*redis.DialConfig		`yaml:"dial"`
	}

	cfgContent := fmt.Sprintf(`
name: %v
network: %v
host: %v
pool:
  max_idle: 10
  MaxActive: 10
  IdleTimeout: 10s
  Wait: true
  MaxConnLifetime: 3s
dial:
  database: 0
  connect_timeout: 1s
  read_timeout: 1s
  write_timeout: 1s
`, name, network, address)

	cfg := &RedisConfig{}
	if err := yaml.Unmarshal([]byte(cfgContent), cfg); err != nil { t.Error(err) }
	_, err := redis.CreatePool(cfg.Name, cfg.Network, cfg.Host, cfg.PoolConfig.PoolOptions(), cfg.DialConfig.DialOptions())
	//_, err := redis.CreatePool(name, network, address, nil, nil)
	if err != nil { t.Error(err) }
}


func TestRedisSet(t *testing.T) {
	pool,ok := redis.GetPool(name)
	if !ok { t.Error(fmt.Sprintf("Can not get redis pool %v", name)) }

	r := pool.Do("SET", "test", 123)
	if r.Error() != nil { t.Error(r.Error())}

	r = pool.Do("GET", "test");
	if r.Error() != nil { t.Error(r.Error())}
	v, err := r.Int();
	if err != nil { t.Error(err) }
	if v != 123 {t.Errorf("error result %v, expect %v", v, 123)}

}
