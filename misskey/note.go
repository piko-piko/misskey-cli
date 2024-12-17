package misskey

import (
	"fmt"
  "regexp"
	"github.com/buger/jsonparser"
)

type Note struct {
	Offset    string
	User      *User
	Renote    bool
	Timestamp string
	Text      string
	Attach    string
	Id        string
}

func (note Note) String(newlines bool) (string) {
	var rnsfx string
	if note.Renote == true {
	  rnsfx = "\x1b[35m[RN]\x1b[0m"
  }
	var text string
	if newlines {
		text = note.Text
	} else {
		re := regexp.MustCompile(`\r?\n`)
		text = re.ReplaceAllString(note.Text, "\\n")
  }
  return fmt.Sprintf("%s %s%s %s \x1b[32m%s\x1b[0m\x1b[34m(%s)\x1b[0m", note.Timestamp, rnsfx, note.User, text, note.Attach, note.Id)
}

func NewNote(value []byte) (*Note, error) {
	note := new(Note)
	var err error
	user, _, _ , _ := jsonparser.Get(value,"user")
	note.User, err = NewUser(user)
	// 投稿時刻
	note.Timestamp, err = jsonparser.GetString(value, "createdAt")
	if err != nil {
		return note, err
	}
	note.Timestamp = convert(note.Timestamp)

	// 本文
	note.Text, _ = jsonparser.GetString(value, "text")
//投稿ID(元投稿)
	note.Id, err = jsonparser.GetString(value, "id")
	if err != nil {
		return note, err
	}

	// ファイルが有れば
	filesId, _, _, _ := jsonparser.Get(value, "files")
	if len(filesId) != 2 {
		note.Attach = "   (添付有り)"
	}

	return note, nil
}

