package handlers

import (
	"HTTP_Sever/model"
	"HTTP_Sever/views"
	"fmt"
	"github.com/a-h/templ"
	"net/http"
)

func RenderHello(globalState model.GlobalState) templ.Component {
	return views.Hello(globalState)
}

func RenderIndex(globalState model.GlobalState) templ.Component {
	return views.Index(globalState)
}

func RenderDashboard(dashboardData model.DashboardData, globalState model.GlobalState) templ.Component {
	return views.Dashboard(dashboardData, globalState)
}

func RenderDashboardUpdate(dashboardData model.DashboardData, globalState *model.GlobalState) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		fmt.Printf("PROJECT: %s\n", project)

		globalState.UpdateGlobalStateProject(project)
		fmt.Printf("GLOBAL STATE: %s\n", *globalState)
		RenderDashboard(dashboardData, *globalState)
	})
}
