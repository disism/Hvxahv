package activitypub

import (
	"fmt"
	pb "github.com/disism/hvxahv/api/hvxahv/v1alpha1"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
)

// IsActorExists Get the username in the request url such,
// as "/.well-known/webfinger?resource=acct:hvturingga@0efb43b41a8a.ngrok.io" Will get hvturingga,
// If the match fails, it will return a custom username not found error.
func IsActorExists(resource string) (string, error) {
	if strings.HasPrefix(resource, "acct:") {
		resource = resource[5:]
		if ali := strings.IndexByte(resource, '@'); ali != -1 {
			resource = resource[:ali]
		}
	} else {
		return "", errors.New("Failed to get username.")
	}

	return resource, nil
}

// EXAMPLE 9
// {
//   "@context": ["https://www.w3.org/ns/activitystreams",
//                {"@language": "ja"}],
//   "type": "Person",
//   "id": "https://kenzoishii.example.com/",
//   "following": "https://kenzoishii.example.com/following.json",
//   "followers": "https://kenzoishii.example.com/followers.json",
//   "liked": "https://kenzoishii.example.com/liked.json",
//   "inbox": "https://kenzoishii.example.com/inbox.json",
//   "outbox": "https://kenzoishii.example.com/feed.json",
//   "preferredUsername": "kenzoishii",
//   "name": "石井健蔵",
//   "summary": "この方はただの例です",
//   "icon": [
//     "https://kenzoishii.example.com/image/165987aklre4"
//   ]
// }


type icon struct {
	Type      string `json:"type"`
	MediaType string `json:"mediaType"`
	Url       string `json:"url"`
}

func NewContext() []string {
	ctx := []string{"https://www.w3.org/ns/activitystreams", "https://w3id.org/security/v1"}

	return ctx

}

func NewIcon(url string) *icon {
	t := "Image"
	mt := "image/jpg"
	return &icon{Type: t, MediaType: mt, Url: url}
}


type actor struct {
	Context           []string    `json:"@context"`
	Type              string      `json:"type"`
	ID                string      `json:"id"`
	Following         string      `json:"following"`
	Followers         string      `json:"followers"`
	Liked             string      `json:"liked"`
	Inbox             string      `json:"inbox"`
	Outbox            string      `json:"outbox"`
	PreferredUsername string      `json:"preferredUsername"`
	Name              string      `json:"name"`
	Summary           string      `json:"summary"`
	PublicKey         interface{} `json:"public_key"`
	Icon              *icon       `json:"icon"`
}

// NewActor Return standard ActivityPub protocol user data.
func NewActor(a *pb.AccountsData) *actor {
	addr := viper.GetString("localhost")

	act := fmt.Sprintf("https://%s/u/%s", addr, a.Username)
	box := fmt.Sprintf("https://%s/u/%s/", addr, a.Username)

	publicKey := map[string]string{
		"id":           a.Uuid,
		"owner":        act,
		"publicKeyPem": a.PublicKey,
	}

	actor := &actor{
		Context:           NewContext(),
		Type:              "Person",
		ID:                act,
		Inbox:             box + "inbox",
		Outbox:            box + "outbox",
		PreferredUsername: a.Username,
		Name:              a.Name,
		Summary:           a.Bio,
		PublicKey:         publicKey,
		Icon:              NewIcon(a.Avatar),
	}
	return actor
}
