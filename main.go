package main

import (
	"flag"
	"fmt"
	"go-galtonboard/logic"
	"go-galtonboard/utils"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	var projectRoutes utils.FlagSlice
	flag.Var(&projectRoutes, "path", "Project routes (can be specified multiple times)")
	createDefaultConfig := flag.Bool("default", false, "Create a default configuration file")
	debug := flag.Bool("debug", false, "Enable debug mode, which use default values for the configuration (boolean)")
	cpuCount := flag.Int("cpu", 1, "Number of CPUs to use")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Copyright (c) 2024 Nicolas Aguilera García \nUsage: go-galtonboard [flags]")
		flag.PrintDefaults()
	}
	flag.Parse()

	cpus := *cpuCount
	if cpus > runtime.NumCPU() {
		cpus = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(cpus)

	log.Println("Using", cpus, "CPUs")

	if *debug {
		projectRoutes = append(projectRoutes, "./")
		*createDefaultConfig = true
	}

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
