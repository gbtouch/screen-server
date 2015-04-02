package server

import "gopkg.in/mgo.v2"

type DBImpl struct {
	Session *mgo.Session
	DB      *mgo.Database
}

func (s *DBImpl) InitDB() {
	var err error

	s.Session, err = mgo.Dial(Settings.DB["url"][0])

	check(err)

	s.DB = s.Session.DB(Settings.DB["name"][0])
}

func (s *DBImpl) CloseDB() {
	defer s.Session.Close()
}
