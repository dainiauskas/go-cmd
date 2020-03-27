package cmd

import config "bitbucket.org/butenta/pkg-config"

type Configuration interface {
	GetApp() *config.App
	GetDB() *config.Database
}
