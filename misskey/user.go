package misskey

import (
	"fmt"
	"github.com/buger/jsonparser"
)


type User struct {
  Name      string
	Username  string
	Host      string
	IsCat     bool
}

func (u User) String() (string) {
	var username string
	if u.Host != "" {
		username = u.Username + "@" + u.Host
	} else {
		username = u.Username
	}
  return fmt.Sprintf("\x1b[35m%s(@%s)\x1b[0m", u.Name, username)
}

func NewUser(json []byte) (*User, error) {
	var err error
	u := new(User)

	u.Name, _ = jsonparser.GetString(json, "name")

	//投稿者ID
	u.Username, err = jsonparser.GetString(json, "username")
	if err != nil {
		return u, err
	}
	//ホスト名
	u.Host, err = jsonparser.GetString(json, "host")
  //ねこかどうか
	u.IsCat, err = jsonparser.GetBoolean(json, "isCat")
	return u, err
}

