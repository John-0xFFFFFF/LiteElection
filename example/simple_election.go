package example

import (
	"fmt"
	"time"

	"github.com/John-0xFFFFFF/LiteElection/election"
)

func elect(processName string,done chan struct{}){
	election.InitSimpleRedis("127.0.0.1:6379","",0)
	e:=election.NewSimpleElection("cluster_local",processName,6*time.Second)
	e.Start()
	
	for{
		select{
		case <-done:
			e.Quit()
			return
		default:
			if e.IsLeader(){
				fmt.Printf("%s is leader \n",processName)
			}else{
				fmt.Printf("%s is follower \n",processName)
			}
			time.Sleep(time.Second)
		}
	}
} 