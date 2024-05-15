package clientProcess

import (
	"fmt"
	"winter/messages"
)

func Fake(chanMsg chan messages.MsgDataToProcess) {
	for {
		msg := <-chanMsg
		fmt.Println(msg)
	}
}
