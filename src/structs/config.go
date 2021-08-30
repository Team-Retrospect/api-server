package structs

type Config struct {
	// example debug: false
	Debug bool `yaml:"debug"`

	// example cluster: "cassandra.domain.com"
	Cluster string `yaml:"cluster"`

	// example port: ":443"
	Port string `yaml:"port"`

	// example use_https: true
	UseHTTPS bool `yaml:"use_https"`
	// example fullcert: "/etc/letsencrypt/live/api.domain.com/fullchain.pem"
	FullCert string `yaml:"fullcert"`
	// example privatekey: "/etc/letsencrypt/live/api.domain.com/privkey.pem"
	PrivateKey string `yaml:"privatekey"`
}
