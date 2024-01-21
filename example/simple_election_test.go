package example

import (
	"testing"
	"time"
)

func TestElect(t *testing.T){
	done1:=make(chan struct{})
	done2:=make(chan struct{})
	go elect("goroutine1",done1)
	time.Sleep(time.Second)
	go elect("goroutine2",done2)
	time.Sleep(10*time.Second)
	close(done1)
	time.Sleep(10*time.Second)
	close(done2)
	time.Sleep(2*time.Second)
	
}