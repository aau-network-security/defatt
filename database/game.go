package database

type Game struct {
	ID       string
	Scenario Scenario
	Tag      string
	Name     string
	BlueTeam []GameUser
	RedTeam  []GameUser
}

type Scenario struct {
	ID string
}
