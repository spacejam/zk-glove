# zk-glove
a small collection of command line tools for distributed orchestration, using zookeeper
```
Glove Usage 
  -data="hostname": contents of znode (for discovery)
  -exec="echo yaaas": command to execute
  -threshold=3: max concurrent commands
  -zk="zk://127.0.0.1:2181/somedir": zookeeper URI
```
```
Hat Usage
  -delim=",": delimiter of node contents
  -exec="echo node data:{}": command to execute, {} is replaced with delimiter-separated znode data
  -hardLimit=true: only run command when threshold is reached with different members
  -pollFreq=30: zk polling frequency
  -pollJitter=30: zk polling random jitter
  -threshold=3: max concurrent commands
  -zk="zk://127.0.0.1:2181/somedir": zookeeper URI

```
