package handlers

import (
	"HTTP_Sever/model"
	"HTTP_Sever/views"
	"github.com/a-h/templ"
	"net/http"
)

func RenderHello(globalState model.GlobalState) templ.Component {
	return views.Hello(globalState)
}

func RenderIndex(globalState model.GlobalState) templ.Component {
	return views.Index(globalState)
}

func RenderDashboard(dashboardData model.DashboardData, globalState *model.GlobalState) templ.Component {
	return views.Dashboard(dashboardData, *globalState)
}

func RenderDashboard2(dashboardData model.DashboardData, globalState *model.GlobalState) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")
		globalState.UpdateGlobalStateProject(project)
		templ.Handler(RenderDashboard(dashboardData, globalState)).ServeHTTP(w, r)
	})
}
