package test

import (
	"encoding/json"
	"github.com/margostino/climateline-processor/api"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockServer struct {
	Method  string
	Status  int
	Queries map[string]string
	Body    interface{}
	Url     string
}

func NewMockServer() *MockServer {
	return &MockServer{}
}

func (m *MockServer) withMethod(method string) *MockServer {
	m.Method = method
	return m
}

func (m *MockServer) withStatus(status int) *MockServer {
	m.Status = status
	return m
}

func (m *MockServer) withBody(body interface{}) *MockServer {
	m.Body = body
	return m
}

func (m *MockServer) withQueries(queries map[string]string) *MockServer {
	m.Queries = queries
	return m
}

func (m *MockServer) start(t *testing.T) *MockServer {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != m.Method {
			t.Errorf("unexpected method for mock server: got %v want %v", r.Method, m.Method)
		}

		if len(m.Queries) > 0 {
			for key, value := range m.Queries {
				if r.URL.Query().Get(key) != value {
					t.Errorf("unexpected query for mock server: got %v want %v", r.URL.Query().Get(key), value)
				}
			}
		}

		w.WriteHeader(m.Status)

		if m.Body != nil {
			res, err := json.Marshal(m.Body)
			w.Write(res)
			require.NoError(t, err)
		}

	}))
	m.Url = server.URL
	return m
}

func call(request *http.Request) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(api.Bot)
	handler.ServeHTTP(response, request)
	return response
}
