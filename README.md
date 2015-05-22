# zk-glove
A small collection of command line tools for distributed orchestration, backed by zookeeper.
* *glove* is for attempting to run a command, without crossing a maximum number that may run concurrently.
* *hat* is for monitoring running commands, and performing certain actions when the membership changes, optionally only once a threshold is (re-)crossed.

Together these simple tools can facilitate a simple zk-backed service discovery mechanism.  For example, you can run some caches using glove, and use hat to send changes in membership to clients.

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
  -threshold=3: max considered nodes
  -zk="zk://127.0.0.1:2181/somedir": zookeeper URI
```
