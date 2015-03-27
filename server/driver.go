package server

import (
	"log"

	"gopkg.in/mgo.v2"
)

type DBImpl struct {
	Session *mgo.Session
	DB      *mgo.Database
}

func (s *DBImpl) InitDB() {
	var err error
	s.Session, err = mgo.Dial(DBUrl)
	if err != nil {
		log.Println(err)
	}

	s.DB = s.Session.DB(DBName)
}

func (s *DBImpl) CloseDB() {
	defer s.Session.Close()
}
