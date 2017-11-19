package main

import (
	"encoding/csv"
	"fmt"
	"github.com/ggpd/brackets/env"
	"github.com/spf13/viper"
	"io"
	"os"
	"time"
)

type Env struct{ *env.Env }

func castEnv(e *env.Env) *Env {
	return &Env{e}
}

func main() {

	e := castEnv(env.New())

	viper.SetConfigName("config")
	viper.AddConfigPath("..")
	err := viper.ReadInConfig()
	if err != nil {
		e.Log.Fatalf("Error reading config file: %s \n", err)
	}

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

	e.GenUsers()

}

func (e *Env) GenUsers() {

	file, _ := os.Open("test.csv")
	r := csv.NewReader(file)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		fmt.Println(record)

		dob, _ := time.Parse("01/02/2006", record[4])

		usr := env.User{
			Email:       record[2],
			FirstName:   record[0],
			LastName:    record[1],
			Gender:      genderMap(record[3]),
			DateOfBirth: dob,
		}

		_, err = e.Db.CreateUser(usr, record[5])
		fmt.Println(err)
	}
}

func genderMap(g string) env.Gender {
	switch g {
	case "Male":
		return 0
	case "Female":
		return 1
	}
	return 2
}

func (e *Env) GenTour() {

	file, _ := os.Open("tour.csv")
	r := csv.NewReader(file)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		fmt.Println(record)

		tour := env.Tournament{
			Name: record[0],
		}

		_, err = e.Db.CreateTournament(tour)
		fmt.Println(err)
	}
}
