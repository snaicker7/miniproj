package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var app Application

// func TestMain(m *testing.M) {
// 	app.DB = &dbrepo.TestDBRepo{}
// 	os.Ex
// }

func TestRegisterService(t *testing.T) {
	app := &Application{}
	testServer := httptest.NewServer(app.RegisterService())
	defer func() {
		testServer.Close()
	}()

	t.Run("Check if route exists", func(t *testing.T) {
		router := app.RegisterService()
		var registered = []struct {
			path     string
			expected bool
		}{
			{"/addProduct", true},
			{"/getProduct", false},
		}
		for _, route := range registered {
			if route.expected != routeExists(router, route.path) {
				t.Errorf("route %s is not registered", route.path)
			}
		}
	})
}

func routeExists(router *mux.Router, testRoute string) bool {
	found := false
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil && pathTemplate == testRoute {
			found = true
		}
		return nil
	})
	if !found {
		found = false
	}
	return found
}
func TestAddProduct(t *testing.T) {
	app := &Application{}
	testServer := httptest.NewServer(app.RegisterService())
	requestBody := []byte(`{
		"user_id": 19900,
		"product_name": "Women's Cloth",
		"product_description": "great outerwear jackets for Spring/Autumn/Winter, suitable for many occasions, such as working, hiking, camping, mountain/rock climbing, cycling, traveling or other outdoors. Good gift choice for you or your family member. A warm hearted love to Father, husband or son in this thanksgiving or Christmas Day.",
		"product_images": [
			"https://fakestoreapi.com/img/71li-ujtlUL._AC_UX679_.jpg",
			"https://fakestoreapi.com/img/71-3HjGNDUL._AC_SY879._SX._UX._SY._UY_.jpg"
		],
		"product_price": 2000
	}`)
	defer func() {
		testServer.Close()
	}()
	t.Run("Check success", func(t *testing.T) {
		req, err := http.NewRequest("POST", testServer.URL+"/addProduct", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
		expectedBody := "Product added successfully"
		actualBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(actualBody) != expectedBody {
			t.Errorf("Expected response body '%s', got '%s'", expectedBody, string(actualBody))
		}
	})

}
