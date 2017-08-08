package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/addomafi/ggbeat/beater"
)

func main() {
	err := beat.Run("ggbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
