package election

import "time"

type Election struct {
	Key, Value        string
	TermExpire        time.Duration
	ElectionInterval  time.Duration
	KeepAliveInterval time.Duration
	Verify            func(localValue, remoteValue string) (bool, error)
	isLeader          int32
}

func NewSimpleElection(k, v string, termExpireTime time.Duration) *Election {
	return NewElection(k, v, termExpireTime, 0, 0, nil)
}

func NewElection(k, v string, termExpire, electionInterval, keepAliveInterval time.Duration, verifyFunc func(localValue, remoteValue string) (bool, error)) *Election {
	e := &Election{
		Key:               k,
		Value:             v,
		TermExpire:        termExpire,
		ElectionInterval:  electionInterval,
		KeepAliveInterval: keepAliveInterval,
		Verify:            verifyFunc,
	}
	//default value
	if e.TermExpire <= 0 {
		e.TermExpire = DefaultTermExpireTime
	}
	if e.ElectionInterval <= 0 {
		e.ElectionInterval = e.TermExpire / 3
	}
	if e.KeepAliveInterval <= 0 {
		e.KeepAliveInterval = e.TermExpire / 2
	}
	if e.Verify == nil {
		e.Verify = func(lv, rv string) (bool, error) {
			return lv == rv, nil
		}
	}
	return e
}
