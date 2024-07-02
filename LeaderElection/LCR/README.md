# How to run?
1. Create 5 terminals in present working directory.
2. In each terminal execute exactly one of the following

```
go run . 0
```
```
go run . 1
```
```
go run . 2
```
```
go run . 3
```
```
go run . 4
```

#### Note
Execute all commands within 15 seconds of each other. The algorithm expects the ring to be formed after 15 seconds, and triggers election immediately.

#### Observe
1 of the terminals would be leader and would receive heartbeat from the other terminals every ~3 seconds.