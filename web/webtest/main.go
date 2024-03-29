package webtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"codeberg.org/genofire/golang-lib/database"
	"codeberg.org/genofire/golang-lib/mailer"
	"codeberg.org/genofire/golang-lib/web"
)

var (
	// DBConnection - url to database on setting up default WebService for webtest
	DBConnection = "postgresql://root:root@localhost:26257/defaultdb?sslmode=disable"
)

// Option to configure TestServer
type Option struct {
	ModuleLoader web.ModuleRegisterFunc
	Database     bool
	DBReRun      bool
	DBSetup      func(db *database.Database)
	Mailer       bool
	Logger       *zap.Logger
}

func (option Option) log() *zap.Logger {
	if option.Logger != nil {
		return option.Logger
	}
	return zap.L()
}

// TestServer - to run it without listen an server
type TestServer struct {
	DB          *database.Database
	Mails       chan *mailer.TestingMail
	Close       func()
	gin         *gin.Engine
	WS          *web.Service
	lastCookies []*http.Cookie
}

// Login Request format (maybe just internal usage)
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// New starts WebService for testing
func New(modules web.ModuleRegisterFunc) (*TestServer, error) {
	return Option{ModuleLoader: modules}.New()
}

// NewWithDBSetup allows to reconfigure before ReRun the database - e.g. for adding Migration-Steps
func NewWithDBSetup(modules web.ModuleRegisterFunc, dbCall func(db *database.Database)) (*TestServer, error) {
	return Option{
		Database:     true,
		DBReRun:      true,
		DBSetup:      dbCall,
		ModuleLoader: modules,
	}.New()
}

// New allows to configure WebService for testing
func (option Option) New() (*TestServer, error) {
	log := option.log()

	ws := &web.Service{}
	ws.SetLog(log)
	ws.Session.Name = "mysession"
	ws.Session.Secret = "hidden"

	ts := &TestServer{
		WS:    ws,
		Close: func() {},
	}

	// db setup
	if option.Database {
		ts.DB = &database.Database{
			Testdata:   true,
			Debug:      false,
			LogLevel:   0,
		}
		ts.DB.Connection.URI = DBConnection
		if option.DBSetup != nil {
			option.DBSetup(ts.DB)
		}
		var err error
		if option.DBReRun {
			err = ts.DB.ReRun()
		} else {
			err = ts.DB.Run()
		}
		if err != nil && err != database.ErrNothingToMigrate {
			return nil, err
		}
		if ts.DB.DB == nil {
			return nil, database.ErrNotConnected
		}
		ws.DB = ts.DB.DB
	}

	if option.Mailer {
		mock, mail := mailer.NewFakeServer(log)
		if err := mail.Setup(); err != nil {
			return nil, err
		}
		ws.Mailer = mail
		ts.Mails = mock.Mails
		ts.Close = mock.Close
	}

	// api setup

	gin.EnableJsonDecoderDisallowUnknownFields()
	gin.SetMode(gin.TestMode)

	ws.ModuleRegister(option.ModuleLoader)

	r := gin.Default()
	ws.LoadSession(r)
	ws.Bind(r)
	ts.gin = r
	return ts, nil
}

// DatabaseForget to run a test without a database
func (s *TestServer) DatabaseForget() {
	s.WS.DB = nil
	s.DB = nil
}

// Request sends a request to webtest WebService
func (s *TestServer) Request(method, url string, body interface{}, expectCode int, jsonObj interface{}) error {
	var jsonBody io.Reader
	if body != nil {
		if strBody, ok := body.(string); ok {
			jsonBody = strings.NewReader(strBody)
		} else {
			jsonBodyArray, err := json.Marshal(body)
			if err != nil {
				return err
			}
			jsonBody = bytes.NewBuffer(jsonBodyArray)
		}
	}
	req, err := http.NewRequest(method, url, jsonBody)
	if err != nil {
		return err
	}
	if len(s.lastCookies) > 0 {
		for _, c := range s.lastCookies {
			req.AddCookie(c)
		}
	}
	w := httptest.NewRecorder()
	s.gin.ServeHTTP(w, req)

	// valid statusCode
	if expectCode != w.Code {
		return fmt.Errorf("wrong status code, body: %v", w.Body)
	}

	if jsonObj != nil {
		// fetch JSON
		err = json.NewDecoder(w.Body).Decode(jsonObj)
		if err != nil {
			return err
		}
	}

	result := w.Result()
	if result != nil {
		cookies := result.Cookies()
		if len(cookies) > 0 {
			s.lastCookies = cookies
		}
	}
	return nil
}

// Login to API by send request
func (s *TestServer) Login(login Login) error {
	// POST: correct login
	return s.Request(http.MethodPost, "/api/v1/auth/login", &login, http.StatusOK, nil)
}

// TestLogin to API by default login data
func (s *TestServer) TestLogin() error {
	return s.Login(Login{
		Username: "admin",
		Password: "CHANGEME",
	})
}
