package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	app := App{}

	log.Info().Msgf("Created app %s.", app)
}
