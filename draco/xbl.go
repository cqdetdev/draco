package draco

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/sandertv/gophertunnel/minecraft/auth"
	"golang.org/x/oauth2"
)

var (
	TokenSrc oauth2.TokenSource
)

type jsonToken struct {
	Access  string `json:"access_token"`
	Type    string `json:"token_type"`
	Refresh string `json:"refresh_token"`
}

func CacheTokenNotExists() bool {
	_, s := os.Stat("./token.json")
	return os.IsNotExist(s)
}

func InitializeToken(log *log.Logger) error {
	if CacheTokenNotExists() {
		log.Printf("XBL: New Token")
		var err error
		Token, err := auth.RequestLiveTokenWriter(log.Writer())
		if err != nil {
			panic(err)
		}
		_ = WriteToken(Token)
		TokenSrc = oauth2.StaticTokenSource(Token)
	} else {
		con, _ := ioutil.ReadFile("./token.json")
		data := &jsonToken{}
		err := json.Unmarshal(con, data)
		Token := &oauth2.Token{}

		Token.AccessToken = data.Access
		Token.RefreshToken = data.Refresh
		Token.TokenType = data.Type
		Token.Expiry = time.Now().AddDate(100, 0, 0)

		TokenSrc = oauth2.StaticTokenSource(Token)
		log.Println("Cached XBL Token")
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteToken(token *oauth2.Token) error {
	bytes, err := json.MarshalIndent(*token, "", "	")
	if err != nil {
		return err
	}
	_ = ioutil.WriteFile("./token.json", bytes, 0777)
	return nil
}
