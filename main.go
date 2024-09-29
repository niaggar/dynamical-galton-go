package main

import (
	"flag"
	"go-galtonboard/logic"
	"go-galtonboard/utils"
	"log"
	"runtime"
	"sync"
	"time"
)

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	var projectRoutes utils.FlagSlice
	flag.Var(&projectRoutes, "projectRoute", "Project routes (can be specified multiple times)")
	createDefaultConfig := flag.Bool("createDefaultConfig", false, "Create a default configuration file")
	flag.Parse()

	if len(projectRoutes) == 0 {
		log.Println("No project routes specified")
		return
	}

	start := time.Now()

	var wg sync.WaitGroup
	for _, projectRoute := range projectRoutes {
		wg.Add(1)
		go runConfiguration(projectRoute, *createDefaultConfig, &wg)
	}
	wg.Wait()

	elapsed := time.Since(start)
	log.Println("------------------------------------")
	log.Println("All simulations finished in:", elapsed)
}

func runConfiguration(projectRoute string, createDefaultConfig bool, group *sync.WaitGroup) {
	defer group.Done()

	if createDefaultConfig {
		err := utils.CreateBaseConfig(projectRoute)
		if err != nil {
			log.Println("Error creating the configuration file for", projectRoute)
			return
		}
	}

	config, err := utils.LoadConfig(projectRoute)
	if err != nil {
		log.Println("Error loading the configuration file for", projectRoute)
		return
	}

	engine := logic.NewEngine(*config, projectRoute)

	log.Println("Running simulation for: ", projectRoute)
	start := time.Now()
	engine.Run()
	elapsed := time.Since(start)
	log.Println("Simulation for", projectRoute, "finished in:", elapsed)
}
