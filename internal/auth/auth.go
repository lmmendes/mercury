package auth

import (
	"database/sql"
	"inbox451/internal/config"
	"inbox451/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zerodha/simplesessions/stores/postgres/v3"
	"github.com/zerodha/simplesessions/v3"
)

const (
	UserKey    = "auth_user"
	SessionKey = "auth_session"
)

type Callbacks struct {
	SetCookie func(cookie *http.Cookie, w interface{}) error
	GetCookie func(name string, r interface{}) (*http.Cookie, error)
	GetUser   func(id int) (models.User, error)
}

type Auth struct {
	cfg       config.Config
	sess      *simplesessions.Manager
	sessStore *postgres.Store
	callbacks *Callbacks
	log       *log.Logger
}

func New(cfg config.Config, db *sql.DB, callbacks *Callbacks, lo *log.Logger) (*Auth, error) {
	a := &Auth{
		cfg:       cfg,
		callbacks: callbacks,
		log:       lo,
	}

	a.sess = simplesessions.New(simplesessions.Options{
		EnableAutoCreate: false,
		SessionIDLength:  64,
		Cookie: simplesessions.CookieOptions{
			IsHTTPOnly: true,
			MaxAge:     time.Hour * 24 * 7,
		},
	})

	st, err := postgres.New(postgres.Opt{}, db)
	if err != nil {
		return nil, err
	}
	a.sessStore = st
	a.sess.UseStore(st)
	a.sess.SetCookieHooks(callbacks.GetCookie, callbacks.SetCookie)

	return a, nil
}

func (o *Auth) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, user, err := o.validateSession(c)
		if err != nil {
			c.Set(UserKey, echo.NewHTTPError(http.StatusForbidden, "invalid session"))
			return next(c)
		}
		c.Set(UserKey, user)
		c.Set(SessionKey, sess)
		return next(c)
	}
}

func (o *Auth) validateSession(c echo.Context) (*simplesessions.Session, models.User, error) {
	sess, err := o.sess.Acquire(c.Request().Context(), c, c)
	if err != nil {
		return nil, models.User{}, echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	vars, err := sess.GetMulti("user_id")
	if err != nil {
		return nil, models.User{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	userID, err := o.sessStore.Int(vars["user_id"], nil)
	if err != nil || userID < 1 {
		o.log.Printf("error fetching session user ID: %v", err)
		return nil, models.User{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user, err := o.callbacks.GetUser(userID)
	if err != nil {
		o.log.Printf("error fetching session user: %v", err)
	}

	return sess, user, err
}

func (o *Auth) SaveSession(u models.User, oidcToken string, c echo.Context) error {
	sess, err := o.sess.NewSession(c, c)
	if err != nil {
		o.log.Printf("error creating login session: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error creating session")
	}

	if err := sess.SetMulti(map[string]interface{}{"user_id": u.ID, "oidc_token": oidcToken}); err != nil {
		o.log.Printf("error setting login session: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error creating session")
	}

	return nil
}
