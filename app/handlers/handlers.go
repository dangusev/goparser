package handlers

import (
	"github.com/dangusev/goparser/app/utils"
	"net/http"
)

func MainHandler(c *utils.GlobalContext, w http.ResponseWriter, r *http.Request) {
	c.Templates["main.html"].ExecuteTemplate(w, "base.html", utils.Context{})
}


func QueriesListHandler(c *utils.GlobalContext, w http.ResponseWriter, r *http.Request) {
}

func QueriesAddHandler(c *utils.GlobalContext, w http.ResponseWriter, r *http.Request) {
}

func QueriesDetailHandler(c *utils.GlobalContext, w http.ResponseWriter, r *http.Request) {
}

func QueriesUpdateHandler(c *utils.GlobalContext, w http.ResponseWriter, r *http.Request) {
}

func QueriesDeleteHandler(c *utils.GlobalContext, w http.ResponseWriter, r *http.Request) {
}

func ItemsListHandler(c *utils.GlobalContext, w http.ResponseWriter, r *http.Request) {
}
