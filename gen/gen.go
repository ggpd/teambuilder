package main

import (
	"encoding/csv"
	"fmt"
	"github.com/ggpd/brackets/env"
	"github.com/spf13/viper"
	"io"
	"math/rand"
	"os"
	"time"
	"strconv"
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

	e.GenGames()

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

func (e *Env) GenTeams() {

	file, _ := os.Open("random_team.csv")
	r := csv.NewReader(file)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		fmt.Println(record)

		t := env.Team{
			Name:         record[0],
			TournamentID: uint(rand.Intn(10) + 1),
		}

		_, err = e.Db.CreateTeam(t)
		fmt.Println(err)
	}

}

/*
func (e *Env) GenPlayers() {

	for i := 2213; i <= 3121; i++ {
		pl := env.Player{

			Rank: env.Rank(rand.Intn(100) % 2),
		}

		//e.Db.AddPlayer(pl, rand.Intn(100)+1, i)
	}

}
*/

func (e *Env) GenGames() {
	file, _ := os.Open("test_games.csv")
	r := csv.NewReader(file)
	
		x := 0
		y := 0
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			i, _ := strconv.ParseInt(record[0], 10, 64)
			t := time.Unix(i,0)
	
			r1 := (rand.Intn(10) + 1) + y
			r2 := (rand.Intn(10) + 1) + y

			for r1 == r2 {
				r1 = (rand.Intn(10) + 1) + y
				r2 = (rand.Intn(10) + 1) + y
			}


			game := env.Game{
				Location: record[1],
				Time: t,
				AwayTeam: &env.Team{ID: uint(r1)},
				HomeTeam: &env.Team{ID: uint(r2)},
			}

			_, err = e.Db.CreateGame(game)

			fmt.Println(err)


			if x == 100 {
				y += 10
				x = 0
			}

			x++

		}



}
