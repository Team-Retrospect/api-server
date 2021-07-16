package main

import (
	/* debug */
	"fmt"
	"log"

	// may have to use this when dockerizing
	"github.com/ilyakaznacheev/cleanenv"

	/* database */
	"time"

	// gocql implements a fast and robust Cassandra client for the Go programming language
	"github.com/gocql/gocql"

  // ----
  "github.com/Team-Textrix/cassandra-connector/src/webserver"
  "github.com/Team-Textrix/cassandra-connector/src/structs"
)

var cfg structs.Config
// taking the information from the .yml file and putting it into a Struct
func load_cfg() {
  cleanenv.ReadConfig("./config.yml", &cfg)
}

func main() {
  load_cfg()

  // connect to the Cassandra cluster
  fmt.Println("Connecting to Cassandra...")
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
