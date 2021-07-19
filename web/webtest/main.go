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

	"dev.sum7.eu/genofire/golang-lib/database"
	"dev.sum7.eu/genofire/golang-lib/mailer"
	"dev.sum7.eu/genofire/golang-lib/web"
)

var (
	// DBConnection - url to database on setting up default WebService for webtest
	DBConnection = "user=root password=root dbname=defaultdb host=localhost port=26257 sslmode=disable"
)

// Option to configure TestServer
type Option struct {
	ReRun        bool
	DBSetup      func(db *database.Database)
	ModuleLoader web.ModuleRegisterFunc
}

type testServer struct {
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
func New(modules web.ModuleRegisterFunc) (*testServer, error) {
	return NewWithOption(Option{ModuleLoader: modules})
}

// NewWithDBSetup allows to reconfigure before ReRun the database - e.g. for adding Migration-Steps
func NewWithDBSetup(modules web.ModuleRegisterFunc, dbCall func(db *database.Database)) (*testServer, error) {
	return NewWithOption(Option{
		ReRun:        true,
		DBSetup:      dbCall,
		ModuleLoader: modules,
	})
}

// NewWithOption allows to configure WebService for testing
func NewWithOption(option Option) (*testServer, error) {
	// db setup
	dbConfig := database.Database{
		Connection: DBConnection,
		Testdata:   true,
		Debug:      false,
		LogLevel:   0,
	}
	if option.DBSetup != nil {
		option.DBSetup(&dbConfig)
	}
	var err error
	if option.ReRun {
		err = dbConfig.ReRun()
	} else {
		err = dbConfig.Run()
	}
	if err != nil && err != database.ErrNothingToMigrate {
		return nil, err
	}
	if dbConfig.DB == nil {
		return nil, database.ErrNotConnected
	}

	// api setup
	gin.EnableJsonDecoderDisallowUnknownFields()
	gin.SetMode(gin.TestMode)

	mock, mail := mailer.NewFakeServer()

	err = mail.Setup()
	if err != nil {
		return nil, err
	}

	ws := &web.Service{
		DB:     dbConfig.DB,
		Mailer: mail,
	}
	ws.ModuleRegister(option.ModuleLoader)
	ws.Session.Name = "mysession"
	ws.Session.Secret = "hidden"

	r := gin.Default()
	ws.LoadSession(r)
	ws.Bind(r)
	return &testServer{
		DB:    &dbConfig,
		Mails: mock.Mails,
		Close: mock.Close,
		gin:   r,
		WS:    ws,
	}, nil
}

// DatabaseForget, to run a test without a database
func (s *testServer) DatabaseForget() {
	s.WS.DB = nil
	s.DB = nil
}

// Request sends a request to webtest WebService
func (s *testServer) Request(method, url string, body interface{}, expectCode int, jsonObj interface{}) error {
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
func (s *testServer) Login(login Login) error {
	// POST: correct login
	return s.Request(http.MethodPost, "/api/v1/auth/login", &login, http.StatusOK, nil)
}

// TestLogin to API by default login data
func (s *testServer) TestLogin() error {
	return s.Login(Login{
		Username: "admin",
		Password: "CHANGEME",
	})
}
