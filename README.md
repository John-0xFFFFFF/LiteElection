# Lite Election Algorithm

## Overview

This README provides a concise introduction to a lightweight and stable election algorithm implemented in the Go programming language for Redis. The algorithm ensures that once a role becomes a leader, it will remain the leader as long as it does not experience a crash. In the event of a crash, a new leader will be elected within a specified time interval.

## Features

- Lightweight and stable election algorithm.
- Persistent leader status as long as the leader role is not interrupted by a crash.
- Automatic leader re-election within a defined time window in case of a leader crash.

## Key Parameters

The algorithm involves three crucial parameters: **term duration, election interval, and keep-alive interval**.

- Term Duration: The total duration a leader remains in power.
- Election Interval: The time interval between leader elections.
- Keep-alive Interval: The time interval for renewing the leader's status.

By default, the election interval is set to **1/3** of the term duration, and the keep-alive interval is **1/2** of the expiration time. Consequently, it can be straightforwardly deduced that the time interval for electing a new leader after a leader crash is **[term duration - renewal interval, term duration + election interval]**. In other words, if the term duration is set to 6 seconds, a new leader can be elected within the interval [3, 8] seconds following a leader crash.

## Usage

```go
import github.com/John-0xFFFFFF/LiteElection

func main(){
    //First of all you need to init your redis
    InitSimpleRedis(RedisAddr,Pwd,YourDBNum)
    //or you could use this way to init
    InitRedis(YourExistedRedisClient)

    //Next to init Election Algorithm
    //The key parameter refers to the name you choose for storing election results in Redis. 
    //It is recommended to  select a key name that is closely related to the cluster name.
    key:="election_cluster001"
    //The value parameter is utilized to identify the current leader of the cluster in Redis. 
    //It is recommended to associate this value with the machine/container name where the current process is running.
    value:="machine_001"
    //Then set you leader's term duration.
    //If you also want set election interval and keep-alive interval,
    //you should use function NewElection()
    termDuration:=10*time.Second
    e:=NewSimpleElection(key,value,termDuration)
    e.Start()

    //If you want to check whether you are the leader
    if e.IsLeader(){
        //some leader logic
    }else{
        //some follower logic
    }

    //If you want to stop/quit election
    e.Quit()

}
```

## Contributions

Contributions to enhance the algorithm's features, stability, or documentation are welcome. Feel free to submit pull requests or open issues to discuss potential improvements.

## License

This Redis election algorithm is distributed under the BSD-3 License. See the LICENSE file for more information.
