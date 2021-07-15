package structs

type Config struct {
  // debug: false
  Debug         bool    `yaml:"debug"`

  // cluster: "cassandra.domain.com"
  Cluster       string  `yaml:"cluster"`

  // port: ":443"
  Port          string  `yaml:"port"`

  // use_https: true
  UseHTTPS      bool    `yaml:"use_https"`
  // fullcert: "/etc/letsencrypt/live/api.domain.com/fullchain.pem"
  FullCert      string  `yaml:"fullcert"`
  // privatekey: "/etc/letsencrypt/live/api.domain.com/privkey.pem"
  PrivateKey    string  `yaml:"privatekey"`
}
