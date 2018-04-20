package vodka

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestContext_UserValue(t *testing.T) {
	ctx := &Context{}

	if ctx.UserValue("empty") != nil {
		t.Errorf("expected %v, got %v", nil, ctx.UserValue("empty"))
	}

	ctx.userValue = make(userValue, 0)

	if ctx.UserValue("empty") != nil {
		t.Errorf("expected %v, got %v", nil, ctx.UserValue("empty"))
	}

	str := "foo"
	ctx.SetUserValue("string", str)

	num := 2016
	ctx.SetUserValue("integer", num)

	if !reflect.DeepEqual(ctx.UserValue("integer"), num) {
		t.Errorf("expect %d, got %d", num, ctx.UserValue("integer"))
	}
	if !reflect.DeepEqual(ctx.UserValue("string"), str) {
		t.Errorf("expect %q, got %q", str, ctx.UserValue("string"))
	}
}

func TestContext_IsDelete(t *testing.T) {
	req, _ := http.NewRequest(MethodDelete, "", nil)

	ctx := &Context{Request: req}
	if !ctx.IsDelete() {
		t.Errorf("expected ctx.IsDelete(): %t, got %t", true, ctx.IsDelete())
	}

	req.Method = MethodPost
	if ctx.IsDelete() {
		t.Errorf("expected ctx.IsDelete() = %t, got %t", false, ctx.IsDelete())
	}
}

func TestContext_IsGet(t *testing.T) {
	req, _ := http.NewRequest(MethodGet, "", nil)

	ctx := &Context{Request: req}
	if !ctx.IsGet() {
		t.Errorf("expected ctx.IsGet(): %t, got %t", true, ctx.IsDelete())
	}

	req.Method = MethodPost
	if ctx.IsGet() {
		t.Errorf("expected ctx.IsGet() = %t, got %t", false, ctx.IsDelete())
	}
}

func TestContext_IsPost(t *testing.T) {
	req, _ := http.NewRequest(MethodPost, "", nil)

	ctx := &Context{Request: req}
	if !ctx.IsPost() {
		t.Errorf("expected ctx.IsPost(): %t, got %t", true, ctx.IsDelete())
	}

	req.Method = MethodGet
	if ctx.IsPost() {
		t.Errorf("expected ctx.IsPost() = %t, got %t", false, ctx.IsDelete())
	}
}

func TestContext_IsPut(t *testing.T) {
	req, _ := http.NewRequest(MethodPut, "", nil)

	ctx := &Context{Request: req}
	if !ctx.IsPut() {
		t.Errorf("expected ctx.IsPut(): %t, got %t", true, ctx.IsDelete())
	}

	req.Method = MethodPost
	if ctx.IsPut() {
		t.Errorf("expected ctx.IsPut() = %t, got %t", false, ctx.IsDelete())
	}
}

func TestContext_IsHead(t *testing.T) {
	req, _ := http.NewRequest(MethodHead, "", nil)

	ctx := &Context{Request: req}
	if !ctx.IsHead() {
		t.Errorf("expected ctx.IsHead(): %t, got %t", true, ctx.IsDelete())
	}

	req.Method = MethodPost
	if ctx.IsHead() {
		t.Errorf("expected ctx.IsHead() = %t, got %t", false, ctx.IsDelete())
	}
}

func TestContext_IsAjax(t *testing.T) {
	req, _ := http.NewRequest(MethodGet, "", nil)

	ctx := &Context{Request: req}
	if ctx.IsAjax() {
		t.Errorf("expected ctx.IsAjax(): %t, got %t", false, ctx.IsAjax())
	}

	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	if !ctx.IsAjax() {
		t.Errorf("expected ctx.IsAjax() = %t, got %t", true, ctx.IsAjax())
	}
}

func TestContext_FormFile(t *testing.T) {
	req, _ := http.NewRequest(MethodGet, "", nil)

	ctx := &Context{Request: req}

	key := "key"
	f1, fh1, err1 := ctx.FormFile(key)
	f2, fh2, err2 := req.FormFile(key)

	if f1 != f2 || fh1 != fh2 || err1 != err2 {
		t.Error("failed to get form file")
	}
}

func TestContext_FormValue(t *testing.T) {
	req, _ := http.NewRequest(MethodGet, "", nil)

	ctx := &Context{Request: req}
	key := "key"
	if ctx.FormValue(key) != req.FormValue(key) {
		t.Error("failed to get form value")
	}
}

func TestContext_PostFormValue(t *testing.T) {
	req, _ := http.NewRequest(MethodGet, "", nil)

	ctx := &Context{Request: req}
	key := "key"
	if ctx.PostFormValue(key) != req.PostFormValue(key) {
		t.Error("failed to get post form value")
	}
}

func TestContext_ParseForm(t *testing.T) {
	req, _ := http.NewRequest(MethodGet, "", nil)

	ctx := &Context{Request: req}
	if ctx.ParseForm() != req.ParseForm() {
		t.Error("failed to get parse form")
	}
}

func TestContext_Host(t *testing.T) {
	req, _ := http.NewRequest(MethodGet, "", nil)

	ctx := &Context{Request: req}
	if ctx.Host() != req.Host {
		t.Error("failed to get request host value")
	}
}

func TestContext_Referer(t *testing.T) {
	req, _ := http.NewRequest(MethodGet, "", nil)

	ctx := &Context{Request: req}
	if ctx.Referer() != req.Referer() {
		t.Error("failed to get referer")
	}
}

func TestContext_URL(t *testing.T) {
	req, _ := http.NewRequest(MethodGet, "", nil)

	ctx := &Context{Request: req}
	if ctx.URL() != req.URL {
		t.Error("failed to get request url")
	}
}

type testUser struct {
	Name string `json:"name" xml:"name"`
}

var user = testUser{Name: "foo"}

func TestContext_JSON(t *testing.T) {
	respJson, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	resp := &mockResponseWriter{}
	logger := new(NullLogger)
	ctx := &Context{Response: resp, logger: logger}
	ctx.JSON(http.StatusOK, user)
	if resp.statusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.statusCode)
	}
	if !bytes.Equal(resp.body, respJson) {
		t.Errorf("expected response body %q, got %q", respJson, resp.body)
	}

	ctx.JSON(http.StatusInternalServerError, new(struct{}))
	if resp.statusCode != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, resp.statusCode)
	}
}

func TestContext_XML(t *testing.T) {
	resp := &mockResponseWriter{}
	logger := new(NullLogger)
	ctx := &Context{Response: resp, logger: logger}
	ctx.XML(http.StatusOK, user)
	if resp.statusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.statusCode)
	}

	userObj := testUser{}
	if err := xml.Unmarshal(resp.body, &userObj); err != nil {
		t.Errorf("fialed to unmarshal xml data %q", err)
	} else {
		if userObj.Name != user.Name {
			t.Errorf("expected user name %q, got %q", user.Name, userObj.Name)
		}
	}

	ctx.XML(http.StatusInternalServerError, new(struct{}))
	if resp.statusCode != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, resp.statusCode)
	}

	resp = &mockResponseWriter{}
	ctx.Response = resp
	header := `<?xml version="2.0" encoding="UTF-8"?>` + "\n"
	ctx.XML(http.StatusOK, user, header)
	if len(resp.body) < len(header) {
		t.Error("incorrect xml header")
	} else if !bytes.Equal(resp.body[:len(header)], []byte(header)) {
		t.Errorf("expected xml header %q, got %q", resp.body[:len(header)], header)
	}
}

func TestContext_Logger(t *testing.T) {
	logger := new(NullLogger)
	ctx := &Context{logger: logger}
	if ctx.Logger() != logger {
		t.Error("failed to get logger")
	}
}

func TestContext_Error(t *testing.T) {
	code := http.StatusInternalServerError
	err := http.StatusText(code)

	resp := &mockResponseWriter{}
	ctx := &Context{Response: resp}
	ctx.Error(err, code)

	if resp.statusCode != code {
		t.Errorf("expected stauts code %d, got %d", code, resp.statusCode)
	}
	if string(resp.body) != fmt.Sprintln(err) {
		t.Errorf("expected response body %q, got %q", fmt.Sprintln(err), resp.body)
	}
}

func TestContext_NotFound(t *testing.T) {
	resp := &mockResponseWriter{}
	ctx := &Context{Response: resp}
	ctx.NotFound()

	if resp.statusCode != http.StatusNotFound {
		t.Errorf("expected stauts code %d, got %d", http.StatusNotFound, resp.statusCode)
	}
	body := fmt.Sprintln("404 page not found")
	if body != string(resp.body) {
		t.Errorf("expected response body %q, got %q", body, resp.body)
	}
}

func TestContext_SetContentType(t *testing.T) {
	contentType := "text/html"

	resp := &mockResponseWriter{}
	ctx := &Context{Response: resp}
	ctx.SetContentType(contentType)

	if resp.Header().Get("Content-Type") != contentType {
		t.Error("failed to set content type")
	}
}

func TestContext_SetStatusCodet(t *testing.T) {
	codes := []int{http.StatusOK, http.StatusNotFound, http.StatusInternalServerError}
	for _, code := range codes {
		resp := &mockResponseWriter{}
		ctx := &Context{Response: resp}
		ctx.SetStatusCode(code)

		if resp.statusCode != code {
			t.Errorf("expected stauts code %d, got %d", code, resp.statusCode)
		}
	}
}

func TestContext_Write(t *testing.T) {
	resp := &mockResponseWriter{}
	ctx := &Context{Response: resp}

	msg := []byte("Hello world.")

	n1, err1 := ctx.Write(msg)
	n2, err2 := resp.Write(msg)

	if n1 != n2 || err1 != err2 {
		t.Error("failed to write response")
	}
}

func TestContext_SetLogger(t *testing.T) {
	logger := new(NullLogger)
	ctx := &Context{}
	ctx.SetLogger(logger)

	if ctx.logger != logger {
		t.Error("failed to set server")
	}
}
