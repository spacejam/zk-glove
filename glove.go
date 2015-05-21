package main

import (
	"flag"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var (
	zkConnect string
	zkChroot  string = "/zk-glove"
	zkServers []string
	threshold uint
	cmd       string
)

// Necessary to prevent zk from timing out our connection
// so we can keep our ephemeral znodes
func pinger(c *zk.Conn) {
	for {
		time.Sleep(30 + time.Duration(rand.Intn(30))*time.Second)
		c.Get("/")
	}
}

func run(c *zk.Conn) {
	result, err := c.Create(zkChroot+"/",
		[]byte("yo"),
		zk.FlagEphemeral|zk.FlagSequence,
		zk.WorldACL(777))
	if err != nil {
		panic(err)
	}
	log.Printf("using znode %+v \n", result)
	go pinger(c)

	parts := strings.Fields(cmd)
	head := parts[0]
	tail := parts[1:len(parts)]
	command := exec.Command(head, tail...)
	err = command.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = command.Wait()
	log.Printf("command finished")
}

func main() {
	c, _, err := zk.Connect(zkServers, time.Second*5)
	if err != nil {
		panic(err)
	}
	_, err = c.Create(zkChroot,
		[]byte(""),
		0,
		zk.WorldACL(zk.PermAll))
	if err != nil && err != zk.ErrNodeExists {
		panic(err)
	}
	children, _, err := c.Children(zkChroot)
	if err != nil {
		panic(err)
	}
	log.Printf("all running nodes: [%+v] \n",
		strings.Join(children, ", "))

	if len(children) >= int(threshold) {
		log.Printf("already at our threshold of %+v, exiting\n", threshold)
	} else {
		run(c)
	}
}

func init() {
	flag.StringVar(&zkConnect, "zk", "zk://127.0.0.1:2181/somedir", "zookeeper URI")
	flag.StringVar(&cmd, "exec", "echo yaaas", "command to execute")
	flag.UintVar(&threshold, "threshold", 3, "max concurrent commands")
	flag.Parse()

	// this is to use the canonical zk://host1:ip,host2/zkChroot format
	strippedZKConnect := strings.TrimLeft(zkConnect, "zk://")
	parts := strings.Split(strippedZKConnect, "/")
	if len(parts) == 2 {
		zkChroot = "/" + parts[1]
		zkServers = strings.Split(parts[0], ",")
	} else if len(parts) == 1 {
		zkServers = strings.Split(parts[0], ",")
	}
	log.Printf("using chroot %+v\n", zkChroot)
}
