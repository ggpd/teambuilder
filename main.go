package main

import (
	"time"
	"fmt"

	"github.com/ggpd/brackets/env"
	"github.com/ggpd/brackets/routes"
	//"github.com/gin-gonic/autotls"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	router := gin.Default()
	e := routes.CastEnv(env.New())

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		e.Log.Fatalf("Error reading config file: %s \n", err)
	}

	url := viper.GetString("app.url")
	port := viper.GetInt("app.port")

	userSQL := viper.GetString("sql.username")
	passSQL := viper.GetString("sql.password")
	databaseSQL := viper.GetString("sql.database")
	hostSQL := viper.GetString("sql.host")
	portSQL := viper.GetInt("sql.port")

	passRedis := viper.GetString("redis.password")
	hostRedis := viper.GetString("redis.host")
	portRedis := viper.GetInt("redis.port")

	sqlOptions := env.SQLOptions{
		User:     userSQL,
		Password: passSQL,
		Host:     hostSQL,
		Port:     portSQL,
		Database: databaseSQL,
	}

	redisOptions := env.RedisOptions{
		Host:     hostRedis,
		Port:     portRedis,
		Password: passRedis,
	}

	e.ConnectDb(sqlOptions, redisOptions)

	fm := template.FuncMap{
		"age":     env.Age,
		"httpdob": env.HttpString,
		"time": timeSt,
		"date": dateSt,
	}

	router.SetFuncMap(fm)

	router.LoadHTMLFiles("public/home.html",
		"public/notfound.html",
		"public/user/user_index.html",
		"public/user/user_login.html",
		"public/user/user_edit.html",
		"public/user/user_register.html",
		"public/user/users.html",
		"public/team/team_index.html",
		"public/team/team_edit.html",
		"public/tournament/tournament_index.html",
		"public/tournament/tournaments.html")

	/*
	 * Register Routes
	 */

	router.NoRoute(routes.NotFoundRoute)

	router.Static("/assets", "public/assets")

	router.GET("/", e.GetHomeRoute)

	router.GET("/login", e.GetLoginRoute)
	router.POST("/login", e.PostLoginRoute)

	router.GET("/logout", e.GetLogoutRoute)

	router.GET("/register", e.GetRegisterRoute)
	router.POST("/register", e.PostRegisterRoute)

	router.GET("/settings", e.GetSettingsRoute)
	router.POST("/settings", e.PostSettingsRoute)

	router.GET("/tournament/:selector", e.GetTournamentRoute)
	router.GET("/tournament", e.GetTournamentsRoute)
	router.POST("/tournament/:selector/new", e.PostCreateTeamRoute)
	router.POST("/tournament", e.PostCreateTournamentRoute)

	router.GET("/team/:selector", e.GetTeamRoute)

	router.GET("/team/:selector/edit", e.GetEditTeamRoute)
	router.POST("/team/:selector/edit", e.PostEditTeamRoute)
	router.POST("/team/:selector/add", e.PostAddPlayerRoute)
	router.GET("/team/:selector/edit/:user/delete", e.GetDeletePlayerRoute)
	router.POST("/team/:selector/game", e.PostCreateGameRoute)

	router.GET("/user/:selector", e.GetUserRoute)
	router.GET("/user", e.GetUsersRoute)

	e.Log.Printf("Server starting...")
	//e.Log.Fatal(autotls.Run(router, fmt.Sprintf("%v:%v", url, port)))
	router.Run(fmt.Sprintf("%v:%v", url, port))
}

func dateSt(t time.Time) string{
	return t.Format("Monday, January 2, 2006")
}

func timeSt(t time.Time) string {
	return t.Format("3:04 PM")
}

