package restserver

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/colbymilton/challenge-07-26-2023/internal/controller"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var tRouter *gin.Engine

const adminAuth = "e502f4c7c766c54391f08a91d6776cc42d51279f239a97e736c29fecc8c959ed"
const guestAuth = "14907954a147647744d042f874fef7504403f7b974344cbcb5e0a1da9cac783e"

func TestMain(m *testing.M) {
	tRouter = setupRouter()
	_ = server.controller.AddUser(&controller.User{Email: "guest@email.com", Role: controller.RoleGuest})

	os.Exit(m.Run())
}

func TestAuthMiddleware(t *testing.T) {
	// auth is not necessary for GET
	resp, err := doRequest(http.MethodGet, "/users", nil)
	require.NoError(t, err)
	require.NotEqual(t, http.StatusUnauthorized, resp.Code)

	// ensure bad auth doesn't break GET
	resp, err = doRequestWithAuth(http.MethodGet, "/users", "fakeauth", nil)
	require.NoError(t, err)
	require.NotEqual(t, http.StatusUnauthorized, resp.Code)

	subtestAuthMiddleWareRequired(t, http.MethodPost, "/users")
	subtestAuthMiddleWareRequired(t, http.MethodPatch, "/users")
	subtestAuthMiddleWareRequired(t, http.MethodDelete, "/users/test@email.com")
}

func subtestAuthMiddleWareRequired(t *testing.T, method, path string) {
	// works with proper auth
	resp, err := doRequestWithAuth(method, path, adminAuth, nil)
	require.NoError(t, err)
	require.NotEqual(t, http.StatusUnauthorized, resp.Code)
	require.NotEqual(t, http.StatusForbidden, resp.Code)

	// does not work without auth
	resp, err = doRequest(method, path, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.Code)

	// does not work with unauthorized
	resp, err = doRequestWithAuth(method, path, guestAuth, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, resp.Code)
}

func TestGetUsers(t *testing.T) {
	// get users
	resp, err := doRequest(http.MethodGet, "/users", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.Code)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.True(t, strings.Contains(string(respBody), `{"email":"admin@email.com","role":"admin"}`))
	require.True(t, strings.Contains(string(respBody), `{"email":"guest@email.com","role":"guest"}`))
}

func TestAddUsers(t *testing.T) {
	// add user
	userBytes := `{"email":"first@email.com","role":"guest"}`
	resp, err := doRequestWithAuth(http.MethodPost, "/users", adminAuth, []byte(userBytes))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.Code)

	// add duplicate user
	resp, err = doRequestWithAuth(http.MethodPost, "/users", adminAuth, []byte(userBytes))
	require.NoError(t, err)
	require.Equal(t, http.StatusConflict, resp.Code)

	// add invalid user
	userBytes = `{"email":"guest@email.com","role":"invalid"}`
	resp, err = doRequestWithAuth(http.MethodPost, "/users", adminAuth, []byte(userBytes))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.Code)

	// add without auth token
	resp, err = doRequest(http.MethodPost, "/users", []byte(userBytes))
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestUpdateUser(t *testing.T) {
	// update user
	userBytes := `{"email":"guest@email.com","role":"admin"}`
	resp, err := doRequestWithAuth(http.MethodPatch, "/users", adminAuth, []byte(userBytes))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.Code)

	// update non-existant user
	userBytes = `{"email":"wrong@email.com","role":"admin"}`
	resp, err = doRequestWithAuth(http.MethodPatch, "/users", adminAuth, []byte(userBytes))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.Code)
}

func TestDeleteUser(t *testing.T) {
	// delete user
	resp, err := doRequestWithAuth(http.MethodDelete, "/users/guest@email.com", adminAuth, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.Code)

	// delete non-existant user
	resp, err = doRequestWithAuth(http.MethodDelete, "/users/wrong@email.com", adminAuth, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.Code)
}

func doRequest(method, path string, body []byte) (*httptest.ResponseRecorder, error) {
	return doRequestWithAuth(method, path, "", body)
}

func doRequestWithAuth(method, path, authToken string, body []byte) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, path, bytes.NewBuffer(body))
	if err != nil {
		return &httptest.ResponseRecorder{}, err
	}
	if authToken != "" {
		req.Header.Set("Authorization", authToken)
	}

	resp := httptest.NewRecorder()
	tRouter.ServeHTTP(resp, req)
	return resp, nil
}
