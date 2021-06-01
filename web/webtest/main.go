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
	"github.com/stretchr/testify/assert"

	"dev.sum7.eu/genofire/golang-lib/database"
	"dev.sum7.eu/genofire/golang-lib/web"
)

var (
	// DBConnection - url to database on setting up default WebService for webtest
	DBConnection = "user=root password=root dbname=defaultdb host=localhost port=26257 sslmode=disable"
)

type testServer struct {
	db          *database.Database
	gin         *gin.Engine
	ws          *web.Service
	assert      *assert.Assertions
	lastCookies []*http.Cookie
}

// Login Request format (maybe just internal usage)
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// New starts WebService for testing
func New(assert *assert.Assertions) *testServer {
	// db setup
	dbConfig := database.Database{
		Connection: DBConnection,
		Testdata:   true,
		Debug:      false,
		LogLevel:   0,
	}
	err := dbConfig.Run()
	if err != nil && err != database.ErrNothingToMigrate {
		fmt.Println(err.Error())
		assert.Nil(err)
	}
	assert.NotNil(dbConfig.DB)

	// api setup
	gin.EnableJsonDecoderDisallowUnknownFields()
	gin.SetMode(gin.TestMode)

	ws := &web.Service{
		DB: dbConfig.DB,
	}
	ws.Session.Name = "mysession"
	ws.Session.Secret = "hidden"

	r := gin.Default()
	ws.LoadSession(r)
	ws.Bind(r)
	return &testServer{
		db:     &dbConfig,
		gin:    r,
		ws:     ws,
		assert: assert,
	}
}

// DatabaseMigration set up a migration on webtest WebService
func (s *testServer) DatabaseMigration(f func(db *database.Database)) {
	f(s.db)
	s.db.MigrateTestdata()
}

// Request sends a request to webtest WebService
func (s *testServer) Request(method, url string, body interface{}, expectCode int, jsonObj interface{}) {
	var jsonBody io.Reader
	if body != nil {
		if strBody, ok := body.(string); ok {
			jsonBody = strings.NewReader(strBody)
		} else {
			jsonBodyArray, err := json.Marshal(body)
			s.assert.Nil(err, "no request created")
			jsonBody = bytes.NewBuffer(jsonBodyArray)
		}
	}
	req, err := http.NewRequest(method, url, jsonBody)
	s.assert.Nil(err, "no request created")
	if len(s.lastCookies) > 0 {
		for _, c := range s.lastCookies {
			req.AddCookie(c)
		}
	}
	w := httptest.NewRecorder()
	s.gin.ServeHTTP(w, req)

	// valid statusCode
	s.assert.Equal(expectCode, w.Code, "expected http status code")
	if expectCode != w.Code {
		fmt.Printf("wrong status code, body:%v\n", w.Body)
		return
	}

	if jsonObj != nil {
		// fetch JSON
		err = json.NewDecoder(w.Body).Decode(jsonObj)
		s.assert.Nil(err, "decode json")
	}

	result := w.Result()
	if result != nil {
		cookies := result.Cookies()
		if len(cookies) > 0 {
			s.lastCookies = cookies
		}
	}
}

// Login to API by send request
func (s *testServer) Login(login Login) {
	// POST: correct login
	s.Request(http.MethodPost, "/api/v1/auth/login", &login, http.StatusOK, nil)
}

// TestLogin to API by default login data
func (s *testServer) TestLogin() {
	s.Login(Login{
		Username: "admin",
		Password: "CHANGEME",
	})
}
