package db

//Database interface is to allow implementing Plug method for any type of Database
//Return a database handle
type Database interface {
	Plug(dbURL string) (interface{}, error)
}

