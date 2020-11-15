package main

import (
	"net/http"
	"regexp"
	"strconv"
)

type route struct {
	Matcher *regexp.Regexp
	Method  string
	Handler func(w http.ResponseWriter, r *http.Request)
}

var rUsers = regexp.MustCompile(`^/users/?$`)
var rUser = regexp.MustCompile(`^/users/(\d+)$`)

var routes = []route{
	route{
		rUsers,
		"GET",
		usersIndex,
	},
	route{
		rUser,
		"GET",
		usersGet,
	},
	route{
		rUser,
		"POST",
		usersUpdate,
	},
	route{
		rUsers,
		"POST",
		usersCreate,
	},
}

type routeHandler struct {
	routes []route
}

func (h routeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, r := range h.routes {
		if r.Matcher.MatchString(req.URL.Path) && req.Method == r.Method {
			r.Handler(w, req)
			return
		}
	}
	notFound(w, req)
}

// TODO use named capture group so that you can make a generic function to use
func extractID(req *http.Request) uint {
	matches := rUser.FindStringSubmatch(req.URL.Path)
	id, _ := strconv.Atoi(matches[1])
	return uint(id)
}
