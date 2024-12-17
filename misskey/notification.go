package misskey

import (
	"github.com/buger/jsonparser"
)

type Notification struct {
  Id string
	CreatedAt string
	Type string
	User *User
	Note *Note
  Reaction string
	Achievement string
}

func stringInArray(target string, arr []string) bool {
  for _, str := range arr {
    if str == target {
      return true
    }
  }
  return false
}

func NewNotification(json []byte) (Notification, error) {
	var nf Notification
	var err error

  nf.Id, _ = jsonparser.GetString(json, "id")
	nf.CreatedAt, _ = jsonparser.GetString(json, "createdAt")
	nf.CreatedAt = convert(nf.CreatedAt)
	nf.Type, _ = jsonparser.GetString(json, "type")

	if !stringInArray(nf.Type,
	  []string{
			"roleAssigned","achievementEarned","app","reaction:grouped","renote:grouped","test",
		}) {
		user, _, _, _ := jsonparser.Get(json, "user")
	  nf.User, err = NewUser(user)
	}

	if !stringInArray(nf.Type,
		[]string{
			"follow","receiveFollowRequest","followRequestAccepted","roleAssigned","achievementEairned","app","test",
		}) {
		note, _ , _ , _ := jsonparser.Get(json,"note")
		nf.Note, _ = NewNote(note)
  }

  if nf.Type == "reaction" {
		nf.Reaction, _ = jsonparser.GetString(json, "reaction")
  }

  if nf.Type == "achievementEarned" {
		nf.Achievement, _ = jsonparser.GetString(json, "achievement")
  }
  //roleAssigned has role (string)
	//app has body, header and icon (strings)
	//reaction:grouped has reactions (array of objects)
	//renote:grouped has users (array of users)
	if err != nil {
    return nf, err
	}
  return nf, nil
}
