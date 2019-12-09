/*
 * Swagger Blog
 *
 * A Simple Blog
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	createdb()
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	str, _ := ioutil.ReadFile("./go/api.json")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(str)
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/v2/",
		Index,
	},

	Route{
		"GetArticleById",
		strings.ToUpper("Get"),
		"/v2/article/{id}",
		GetArticleById,
	},

	Route{
		"GetArticles",
		strings.ToUpper("Get"),
		"/v2/articles",
		GetArticles,
	},

	Route{
		"GetCommentsOfArticle",
		strings.ToUpper("Get"),
		"/v2/article/{id}/comments",
		GetCommentsOfArticle,
	},

	Route{
		"CreateComment",
		strings.ToUpper("Post"),
		"/v2/article/{id}/comment",
		CreateComment,
	},

	Route{
		"SignIn",
		strings.ToUpper("Post"),
		"/v2/auth/signin",
		SignIn,
	},

	Route{
		"SignUp",
		strings.ToUpper("Post"),
		"/v2/auth/signup",
		SignUp,
	},
	Route{
		"OPTIONS",
		strings.ToUpper("options"),
		"/v2/auth/signin",
		Options,
	},

	Route{
		"OPTIONS",
		strings.ToUpper("options"),
		"/v2/article/{id}/comment",
		Options,
	},

	Route{
		"OPTIONS",
		strings.ToUpper("options"),
		"/v2/auth/signup",
		Options,
	},
}

func Options(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,Content-Type,Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)
}
