package election

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

type Election struct {
	Key, Value        string
	TermDuration      time.Duration
	ElectionInterval  time.Duration
	KeepAliveInterval time.Duration
	Verify            func(localValue, remoteValue string) (bool, error)
	isLeader          int32
	done              chan struct{}
}

// NewSimpleElection
// The key parameter refers to the name you choose for storing election results in Redis. It is recommended to  select a key name that is closely related to the cluster name.
// The value parameter is utilized to identify the current leader of the cluster in Redis. It is recommended to associate this value with the machine/container name where the current process is running.
func NewSimpleElection(k, v string, termDuration time.Duration) *Election {
	return NewElection(k, v, termDuration, 0, 0, nil)
}

func NewElection(k, v string, termExpire, electionInterval, keepAliveInterval time.Duration, verifyFunc func(localValue, remoteValue string) (bool, error)) *Election {
	e := &Election{
		Key:               k,
		Value:             v,
		TermDuration:      termExpire,
		ElectionInterval:  electionInterval,
		KeepAliveInterval: keepAliveInterval,
		Verify:            verifyFunc,
		done:              make(chan struct{}),
	}
	//default value
	if e.TermDuration <= 0 {
		e.TermDuration = DefaultTermExpireTime
	}
	if e.ElectionInterval <= 0 {
		e.ElectionInterval = e.TermDuration / 3
	}
	if e.KeepAliveInterval <= 0 {
		e.KeepAliveInterval = e.TermDuration / 2
	}
	if e.Verify == nil {
		e.Verify = func(lv, rv string) (bool, error) {
			return lv == rv, nil
		}
	}
	return e
}

func (e *Election) Start() {
	go e.Do()
}

func (e *Election) Quit() {
	close(e.done)
}

func (e *Election) IsLeader() bool {
	return atomic.LoadInt32(&e.isLeader) == 1
}

// todo change with lua script
func (e *Election) KeepAlive() error {
	ctx := context.Background()
	tfx := func(tx *redis.Tx) error {
		curValue, err := tx.Get(ctx, e.Key).Result()
		if err != nil {
			return err
		}
		//todo with redis.nil
		isLeader, err := e.Verify(e.Value, curValue)
		if err != nil {
			return err
		}
		if !isLeader {
			return errors.New("only leader allowed to keep alive")
		}
		_, err = tx.Pipeline().TxPipelined(ctx, func(p redis.Pipeliner) error {
			tx.Expire(ctx, e.Key, e.TermDuration)
			return nil
		})
		return err
	}
	return RedisClient.Watch(ctx, tfx, e.Key)
}

func (e *Election) Elect() (bool, error) {
	rawValue, err := setNxExScript.Run(context.Background(), RedisClient, []string{e.Key}, e.Value, e.TermDuration).Result()
	if err != nil {
		return false, err
	}
	remoteValue, ok := rawValue.(string)
	if !ok {
		return false, errors.New("invalid value to verify")
	}
	if remoteValue == "OK" {
		return true, nil
	}
	isLeader, err := e.Verify(remoteValue, e.Value)
	if err != nil {
		return false, err
	}
	if isLeader {
		return true, nil
	}
	return false, err
}

var setNxExScript = redis.NewScript(`
local current_value = redis.call('GET', KEYS[1])
if current_value then
	return current_value
else
	redis.call('SET', KEYS[1], ARGV[1], 'EX', ARGV[2], 'NX')
	return "OK"
end
`)

func (e *Election) Do() {

	for {
		for {
			select {
			case <-e.done:
				return
			default:
				isLeader, err := e.Elect()
				if err == nil {
					if isLeader {
						atomic.StoreInt32(&e.isLeader, 1)
						break
					}
					atomic.StoreInt32(&e.isLeader, 0)
				} else {
					log.Printf("error occured while electing:%v", err)
				}
				time.Sleep(e.ElectionInterval)
			}
		}

		for {
			select {
			case <-e.done:
				return
			default:
				err := e.KeepAlive()
				if err != nil {
					atomic.StoreInt32(&e.isLeader, 0)
					break
				}
				time.Sleep(e.KeepAliveInterval)
			}
		}

	}
}
