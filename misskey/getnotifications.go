package misskey

import (
	"encoding/json"
//	"errors"
	"fmt"
	"os"
	"sort"
	"github.com/buger/jsonparser"
	"github.com/mattn/go-colorable"
)

// タイムライン取得
func (c *Client) GetNotifications(plainPrint bool, limit int, sinceId string) error {
	output = colorable.NewColorableStdout()
	var body = map[string]interface{}{
		"i": c.InstanceInfo.Token,
		"limit":  limit,
	}
	if sinceId != "" {
		body["sinceId"] = &sinceId
	}
	jsonByte, err := json.Marshal(body)
	if err != nil {
		return err
	}

	if err := c.apiPost(jsonByte,"i/notifications"); err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout,"Notifications: @" + c.InstanceInfo.UserName + " (" + c.InstanceInfo.Host + ")")
	if !plainPrint {
		printLine()
	}
	notifications := []Notification{}
	// Something peculiuar happends to the order of the notifications - it seems like the order depends on if the sinceId parameter is present, so i packed it up into Slice and sort before displaying
	jsonparser.ArrayEach(c.resBuf.Bytes(), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var nf Notification

		nf, err = NewNotification(value)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		if !((sinceId != "") && (nf.Id == sinceId)) {
			notifications = append(notifications,nf)
		}
	})
	sort.Slice(notifications,
		func(p, q int) bool {  
			return notifications[p].CreatedAt < notifications[q].CreatedAt;
		},
	)

	for _,nf := range notifications {
	var str string
	str = fmt.Sprintf("%s\t%s\t", nf.Id, nf.CreatedAt)
		switch nf.Type {
			case "follow":
				str += fmt.Sprintf("%s followed", nf.User )
			case "renote":
				str += fmt.Sprintf("%s renoted\t%s", nf.User, nf.Note.String(!plainPrint))
//				str += fmt.Sprintf("    %s", nf.Note 
			case "followRequestAccepted":
				str += fmt.Sprintf("%s accepted the follow request", nf.User )
			case "reaction":
				str += fmt.Sprintf("%s reacted with %s to\t%s", nf.User, nf.Reaction, nf.Note.String(!plainPrint) )
//				str += fmt.Sprintf("    %s", nf.Note )
			case "achievementEarned":
				str += fmt.Sprintf("Earned achievement %s", nf.Achievement )
			default:
				str += fmt.Sprintf("%s", nf.Type )
			}
		fmt.Fprintln(output, str)
	}
	return nil
}

