# Minired [![version](https://img.shields.io/badge/version-0.0.1-green.svg)](https://semver.org)
Minimalist cache database with publish-subscribe capabilities

## Features🍕
* Handling major Redis commands
  - SET
  - GET
  - PING
  - HSET
  - HGET
  - HGETALL
* Persistence storage using [AOF](https://redis.io/docs/latest/operate/oss_and_stack/management/persistence/)
* Subscribing to channels
* Multi-client connections
* Handling transactions.

## Build✨
Clone the repo

```git
  git clone github.com/tdadadavid/minired
  cd minired
```
Build minired🚀
```git

  go build -o minired

```
## Potential improvements🤔
* [] Add commands pipelining feature
* [] Build and Deploy a webapp for interacting with it.
* [] Add other redis commands
