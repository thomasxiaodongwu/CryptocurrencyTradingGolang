package clientAlpha

import (
	"fmt"
	"winter/global"
	"winter/messages"
)

func Fake(chanStmsg chan messages.AggStMsg) {
	count := 0
	for {
		msg := <-chanStmsg
		count += msg.ExID
		fmt.Println(global.ExIDList[msg.ExID])
	}
}
