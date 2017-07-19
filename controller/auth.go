package controller

import (
	"gopkg.in/authboss.v1"
	_ "gopkg.in/authboss.v1/auth"
	_ "gopkg.in/authboss.v1/confirm"
	_ "gopkg.in/authboss.v1/lock"
	_ "gopkg.in/authboss.v1/recover"
	_ "gopkg.in/authboss.v1/register"
	_ "gopkg.in/authboss.v1/remember"
	"net/http"
	"github.com/justinas/nosurf"
	"html/template"
	"time"
	"github.com/toasterson/mozaik/logger"
	"github.com/toasterson/mozaik/models"
	"github.com/toasterson/mozaik/auth"
	"github.com/gorilla/mux"
)

var (
	AuthBoss = authboss.New()
	authEnabled bool = false
)

func enableAuth(router *mux.Router){
	if authEnabled == false {
		router.PathPrefix("/auth").Handler(AuthBoss.NewRouter())
		authEnabled = true
	}
}

func SetupAuthboss(tplGlob string, debug bool, partialGlob string, funcs template.FuncMap) {
	AuthBoss.Storer = auth.Database
	AuthBoss.OAuth2Storer = auth.Database
	AuthBoss.MountPath = "/auth"
	AuthBoss.LogWriter = logger.GetLoggerWriter()
	AuthBoss.LayoutDataMaker = layoutData
	AuthBoss.NotFoundHandler = http.HandlerFunc(NotFound)


	AuthBoss.XSRFName = "csrf_token"
	AuthBoss.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
		return nosurf.Token(r)
	}

	AuthBoss.CookieStoreMaker = auth.NewCookieStorer
	AuthBoss.SessionStoreMaker = auth.NewSessionStorer

	AuthBoss.Mailer = authboss.LogMailer(logger.GetLoggerWriter())

	AuthBoss.Policies = []authboss.Validator{
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

	if err := AuthBoss.Init(); err != nil {
		logger.Critical(err)
	}
}

func layoutData(w http.ResponseWriter, r *http.Request) authboss.HTMLData {
	currentUserName := ""
	userInter, err := AuthBoss.CurrentUser(w, r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*models.User).Username
	}

	return authboss.HTMLData{
		"loggedin":               userInter != nil,
		"username":               "",
		authboss.FlashSuccessKey: AuthBoss.FlashSuccess(w, r),
		authboss.FlashErrorKey:   AuthBoss.FlashError(w, r),
		"current_user_name":      currentUserName,
	}
}

var Funcs = template.FuncMap{
	"formatDate": func(date time.Time) string {
		return date.Format("2006/01/02 03:04pm")
	},
	"yield": func() string { return "" },
}