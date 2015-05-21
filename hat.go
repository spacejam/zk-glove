package main

import (
	"bufio"
	"flag"
	"log"
	"math/rand"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var (
	zkConnect string
	zkChroot  = "/zk-glove"
	zkServers []string
	threshold uint
	hardLimit bool
	frequency uint
	jitter    int
	cmd       string
	delim     string
)

func splay() {
	time.Sleep((time.Duration(frequency) * time.Second) + (time.Duration(rand.Intn(jitter)) * time.Second))
}

//TODO fix execution
func run(nodes string) {
	parts := strings.Fields(cmd)
	for i := 0; i < len(parts); i++ {
		if parts[i] == "{}" {
			parts[i] = nodes
		}
	}
	head := parts[0]
	tail := parts[1:len(parts)]
	command := exec.Command(head, tail...)
	stdout, err := command.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = command.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = command.Wait()
	r := bufio.NewReader(stdout)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			break
		}
		log.Println(line)
	}
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

	children := []string{}
	for {
		newChildren, _, err := c.Children(zkChroot)
		if err != nil {
			log.Fatal(err)
		}
		if len(newChildren) > int(threshold) {
			// TODO sort this and only take lower members
			newChildren = newChildren[:threshold]
		}

		log.Println("1")
		if !reflect.DeepEqual(children, newChildren) {
			log.Println("2", hardLimit)
			if len(newChildren) == int(threshold) || !hardLimit {
				log.Println("3")
				children = newChildren
				log.Printf("all running nodes: [%+v] \n",
					strings.Join(children, ", "))
				data := []string{}
				for i := 0; i < len(children); i++ {
					rawData, _, err := c.Get(zkChroot + "/" + children[i])
					if err != nil {
						log.Println(err)
					} else {
						data = append(data, string(rawData))
					}
				}
				znodes := strings.Join(data, delim)
				run(znodes)
			}
		}
		splay()
	}
}

func init() {
	flag.StringVar(&zkConnect, "zk", "zk://127.0.0.1:2181/somedir", "zookeeper URI")
	flag.StringVar(&cmd, "exec", "echo node data:{}", "command to execute, {} is replaced with delimiter-separated znode data")
	flag.StringVar(&delim, "delim", ",", "delimiter of node contents")
	flag.UintVar(&threshold, "threshold", 3, "max concurrent commands")
	flag.UintVar(&frequency, "pollFreq", 30, "zk polling frequency")
	flag.IntVar(&jitter, "pollJitter", 30, "zk polling random jitter")
	flag.BoolVar(&hardLimit, "hardLimit", true, "only run command when threshold is reached with different members")
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
