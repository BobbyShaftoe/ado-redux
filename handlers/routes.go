package handlers

import (
	"HTTP_Sever/model"
	"HTTP_Sever/views"
	"github.com/a-h/templ"
)

func RenderHello(name string) templ.Component {
	return views.Hello("Nick")
}

func RenderIndex() templ.Component {
	return views.Index()
}

func RenderDashboard(dashboardData model.DashboardData) templ.Component {
	return views.Dashboard(dashboardData)
}
