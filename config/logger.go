package config

import "log"

func InitLogger() {
	// bisa diganti dengan logrus atau zap bila mau
	log.SetPrefix("[golanjutan] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}