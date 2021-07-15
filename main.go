package main

import (
	/* debug */
	"fmt"
	"log"

	/* config */
	// "Clean and minimalistic environment configuration reader for Golang"
	// "reads and parses configuration structure from the file
	// reads and overwrites configuration structure from environment variables
	// writes a detailed variable list to help output"
	//   ReadConfig method takes a string representing the name of a config file
	//    and a pointer to a struct (ex. Config)

	// may have to use this when dockerizing
	"github.com/ilyakaznacheev/cleanenv"

	/* database */
	"time"

	// gocql implements a fast and robust Cassandra client for the Go programming language
	"github.com/gocql/gocql"

  // ----
  "github.com/Team-Textrix/cassandra-connector/webserver"
  "github.com/Team-Textrix/cassandra-connector/structs"
)

func output(contents ...string) {
  if (debug) { fmt.Println(contents) }
}
/* load configs from config.yml */
// property names match the data in the config.yml


// may have to add database credentials to this struct
// when we dockerize the application
// (read from environment variables? err := cleanenv.ReadEnv(&cfg))
//     Password string `env:"PASSWORD"`

var debug bool = false;
var cfg structs.Config
// taking the information from the .yml file and putting it into a Struct
func load_cfg() {
  cleanenv.ReadConfig("config.yml", &cfg)
  debug = cfg.Debug
}


/* orchestrate */
func main() {
  load_cfg()

  // connect to the Cassandra cluster
  output("Connecting to Cassandra...")
  cluster := gocql.NewCluster(cfg.Cluster)
  cluster.Consistency = gocql.Quorum
  cluster.ProtoVersion = 4
  cluster.ConnectTimeout = time.Second * 10
  // cluster.Authenticator = gocql.PasswordAuthenticator{
  //   Username: "Username",
  //   Password: "Password",
  // } //replace the username and password fields with their real settings.
  s, err := cluster.CreateSession()
  if err != nil {
    log.Println(err)
    return
  }
  session := s
  defer session.Close()

  webserver.DeclareRouter(cfg, session)
}
