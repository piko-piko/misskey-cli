package misskey

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/buger/jsonparser"
	"github.com/mattn/go-colorable"
)

var (
	output = colorable.NewColorableStdout()
)

// タイムライン取得
func (c *Client) GetTimeline(plainPrint bool, limit int, mode string) error {
	body := struct {
		I     string `json:"i"`
		Limit int    `json:"limit"`
	}{
		I:     c.InstanceInfo.Token,
		Limit: limit,
	}

	var endpoint string
	if mode == "local" {
		endpoint = "notes/local-timeline"
	} else if mode == "global" {
		endpoint = "notes/global-timeline"
	} else if mode == "home" {
		endpoint = "notes/timeline"
	} else {
		return errors.New("Please select mode in local/home/global")
	}

	jsonByte, err := json.Marshal(body)
	if err != nil {
		return err
	}

	if err := c.apiPost(jsonByte, endpoint); err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout,"Timeline: " + mode + "  @" + c.InstanceInfo.UserName + " (" + c.InstanceInfo.Host + ")")
	if !plainPrint {
		printLine()
	}

	jsonparser.ArrayEach(c.resBuf.Bytes(), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

		// とりあえずTextを持ってきてみる
		_, err = jsonparser.GetString(value, "renoteId")

		var note *Note

		if err != nil {
			note, err = NewNote(value)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			_, err = jsonparser.GetString(value, "replyId")
			if err == nil {
				replyParentValue, _, _, _ := jsonparser.Get(value, "reply")
				replyParent, err := NewNote(replyParentValue)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return
				}
				repStr := fmt.Sprint(replyParent.String(!plainPrint))
				fmt.Fprintln(output, repStr)
				note.Offset = "    "
			}

		} else { // renoteだったら

			renoteValue, _, _, _ := jsonparser.Get(value, "renote")

			note, err = NewNote(renoteValue)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			note.Renote = true
		}
		fmt.Fprintln(os.Stdout, note.String(!plainPrint))
	})

	if err != nil {
		fmt.Fprintln(output, err)
		return err
	}

	return nil
}


