## Go Cache Server

`GoCache` is a simple Cache Server based on [LRU](https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU)) algorithm which is also used on **MemCache**.

### Setup and Run
- Clone this repo: `https://github.com/kadnan/GoCache.git`
- `go run gocached.go lru.go` - By Default it runs on port `9000` with the capacity `5`.
- You can also specify port while running: `go run gocached.go lru.go -port=9002 -capacity=20`

### Commands
- get <keyname> to retrieve a key. If not available it gives a message. For example `get name`.
- set <keyname> <value> to set a key. It returns `1`. For example `set name adnan`.

## Etcetera

- This is my first ever Go program and also first ever implementation of LRU Cache Algorithm so do mention mistake where found or come up witha  PR request.

- There is also a sample go client available in the file `client.go`.