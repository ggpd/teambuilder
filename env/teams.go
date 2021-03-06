package env

import (
	"strings"
	"errors"
)

const (
	createTeam = "INSERT INTO teams (selector, name, tournament_id) VALUES ($1, $2, $3) RETURNING id"
	getTeam    = "SELECT id, tournament_id, name FROM teams WHERE selector=$1"
	updateTeam = "UPDATE teams SET name=$1 WHERE id=$2"
	deleteTeam = "DELETE FROM teams WHERE team_id=$1"

	insertPlayer = "INSERT INTO players (user_id, team_id, rank) VALUES ($1, $2, $3)"
	updatePlayer = "UPDATE players SET rank=$1 WHERE team_id=$2 AND user_id=$3"
	deletePlayer = "DELETE FROM players WHERE team_id=$1 AND user_id=$2" //FIX
	selectPlayer = "SELECT rank FROM players WHERE team_id=$1 AND user_id=$2"

	selectPlayers    = "SELECT users.selector, users.id, users.first_name, users.last_name, users.gender, users.dob, users.email, players.rank FROM users JOIN players ON players.user_id=users.id WHERE players.team_id=$1"
	deleteAllPlayers = "DELETE FROM players WHERE team_id=$1"

	selectAllPlayerTeam = "SELECT teams.id, teams.selector, teams.name FROM players JOIN teams ON teams.id=players.team_id WHERE players.user_id=$1"
)

type Rank int

const (
	Member  Rank = 0
	Manager Rank = 1
)

func (r Rank) String() string {
	switch r {
	case Member:
		return "Member"
	case Manager:
		return "Manager"
	}
	return ""
}

func ToRank(st string) Rank {
	switch(strings.ToLower(st)){
	case "manager":
		return Manager
	case "member":
		return Member
	}

	return Member
}

type teamDatastore interface {
	CreateTeam(team Team) (*Team, error)
	GetTeam(selector string, full bool) (*Team, error)
	GetTeams(user User) ([]Team, error)
	UpdateTeam(team Team) error
	DeleteTeam(selector Team) error

	GetRank(team Team, user User) (Rank, error)
	AddPlayer(team Team, pl Player) (*Team, error)
	UpdatePlayer(team Team, pl Player) (*Team, error)
	DeletePlayer(team Team, pl Player) (*Team, error)
}

type Team struct {
	Selector

	ID           uint
	TournamentID uint

	Name    string
	Players []*Player
}

type Player struct {
	User
	Rank
}

func (d *db) CreateTeam(team Team) (*Team, error) {
	selector := d.GenerateSelector(selectorLen)
	team.sel = selector
	err := d.QueryRow(createTeam, selector, team.Name, team.TournamentID).Scan(&team.ID)
	if err != nil {
		return nil, errors.New("failed to create team")
	}

	return &team, nil
}

func (d *db) GetTeam(selector string, full bool) (*Team, error) {
	var team Team
	team.sel = selector

	err := d.QueryRow(getTeam, team.Selector.String()).Scan(&team.ID, &team.TournamentID, &team.Name)
	if err != nil {
		d.Logger.Println(err)
		return nil, errors.New("failed to get team")
	}

	if full {
		rows, err := d.Query(selectPlayers, team.ID)
		if err == nil {
			for rows.Next() {
				var pl Player
				rows.Scan(&pl.sel, &pl.ID, &pl.FirstName, &pl.LastName, &pl.Gender, &pl.DateOfBirth, &pl.Email, &pl.Rank)
				team.Players = append(team.Players, &pl)
			}
		}
	}

	return &team, nil
}

func (d *db) UpdateTeam(team Team) error {
	_, err := d.DB.Exec(updateTeam, team.Name, team.ID)
	if err != nil {
		d.Logger.Println(err)
		return errors.New("failed to update team")
	}

	return nil
}

func (d *db) DeleteTeam(team Team) error {
	tx, err := d.DB.Begin()
	if err != nil {
		return errors.New("Couldn't delete team")
	}
	tx.Exec(deleteTeam, team.ID)
	tx.Exec(deleteAllPlayers, team.ID)
	tx.Commit()

	return nil
}

func (d *db) GetTeams(user User) ([]Team, error) {
	rows, err := d.Query(selectAllPlayerTeam, user.ID)
	if err != nil {
		return nil, err
	}
	var rt []Team
	for rows.Next(){
		var team Team
		rows.Scan(&team.ID, &team.sel, &team.Name)
		rt = append(rt, team)
	}

	return rt, nil
}

func (d *db) AddPlayer(team Team, pl Player) (*Team, error){
	_, err := d.Exec(insertPlayer, pl.ID, team.ID, int(pl.Rank))
	if err != nil {
		return nil, err
	}

	team.Players = append(team.Players, &pl)

	return &team, nil
}

func (d *db) UpdatePlayer(team Team, pl Player) (*Team, error){
	_, err := d.Exec(updatePlayer, int(pl.Rank), team.ID, pl.ID)
	if err != nil {
		d.Logger.Println(err)
		return nil, err
	}

	return &team, nil
}

func (d *db) DeletePlayer(team Team, pl Player) (*Team, error){
	_, err := d.Exec(deletePlayer, team.ID, pl.ID)
	if err != nil {
		return nil, err
	}

	return &team, nil
}


func (d *db) GetRank(team Team, user User) (Rank, error){
	var rt Rank
	err := d.QueryRow(selectPlayer, team.ID, user.ID).Scan(&rt)
	if err != nil {
		return -1, err
	}

	return rt, nil
}
