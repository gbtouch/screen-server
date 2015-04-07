package server

import "gopkg.in/mgo.v2"

type DBImpl struct {
	Session *mgo.Session
	DB      *mgo.Database
}

func (s *DBImpl) InitDB() {
	s.Session, _ = mgo.Dial(Settings.DB["url"][0])

	if s.Session != nil {
		s.DB = s.Session.DB(Settings.DB["name"][0])
	}
}
