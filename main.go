package main

import (
	"flag"
	"go-galtonboard/logic"
	"go-galtonboard/utils"
	"log"
	"runtime"
	"time"
)

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	projectRoute := flag.String("projectRoute", "", "Project route")
	createDefaultConfig := flag.Bool("createDefaultConfig", false, "Create a default configuration file")
	flag.Parse()

	if *createDefaultConfig {
		err := utils.CreateBaseConfig(*projectRoute)
		if err != nil {
			log.Fatalf("Error creating the configuration file: %v", err)
		}
	}

	config, err := utils.LoadConfig(*projectRoute)
	if err != nil {
		log.Fatalf("Error loading the configuration file: %v", err)
	}

	log.Printf("Running the logic...")

	start := time.Now()
	engine := logic.NewEngine(*config, *projectRoute)
	engine.Run()
	elapsed := time.Since(start)

	log.Printf("Engine finished")
	log.Printf("Took %s", elapsed)
}
