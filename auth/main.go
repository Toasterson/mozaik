package auth

import (
	"gopkg.in/authboss.v1"
	"io/ioutil"
	"path/filepath"
	"net/http"
	"github.com/justinas/nosurf"
	"os"
	"log"
	"github.com/toasterson/mozaik/models"
	"html/template"
	"time"
)

var (
	Ab        = authboss.New()
	database  = NewMemStorer()
)

func SetupAuthboss() {
	Ab.Storer = database
	Ab.OAuth2Storer = database
	Ab.MountPath = "/auth"
	Ab.ViewsPath = "templates/ab_views"
	Ab.RootURL = `http://localhost:8080`

	Ab.LayoutDataMaker = layoutData

	b, err := ioutil.ReadFile(filepath.Join("templates", "base.html.tpl"))
	if err != nil {
		panic(err)
	}
	Ab.Layout = template.Must(template.New("base").Funcs(Funcs).Parse(string(b)))

	Ab.XSRFName = "csrf_token"
	Ab.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
		return nosurf.Token(r)
	}

	Ab.CookieStoreMaker = NewCookieStorer
	Ab.SessionStoreMaker = NewSessionStorer

	Ab.Mailer = authboss.LogMailer(os.Stdout)

	Ab.Policies = []authboss.Validator{
		authboss.Rules{
			FieldName:       "email",
			Required:        true,
			AllowWhitespace: false,
		},
		authboss.Rules{
			FieldName:       "password",
			Required:        true,
			MinLength:       4,
			MaxLength:       8,
			AllowWhitespace: false,
		},
	}

	if err := Ab.Init(); err != nil {
		log.Fatal(err)
	}
}

func layoutData(w http.ResponseWriter, r *http.Request) authboss.HTMLData {
	currentUserName := ""
	userInter, err := Ab.CurrentUser(w, r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*models.User).Username
	}

	return authboss.HTMLData{
		"loggedin":               userInter != nil,
		"username":               "",
		authboss.FlashSuccessKey: Ab.FlashSuccess(w, r),
		authboss.FlashErrorKey:   Ab.FlashError(w, r),
		"current_user_name":      currentUserName,
	}
}

var Funcs = template.FuncMap{
	"formatDate": func(date time.Time) string {
		return date.Format("2006/01/02 03:04pm")
	},
	"yield": func() string { return "" },
}