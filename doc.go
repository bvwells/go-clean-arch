
  /*
  Go-clean-arch is a linter for enforcing clean architecture principles in Go.
  
  Usage:
  	go-clean-arch [flags] [path ...]
  
  The flags are:
  	-c
  		Config file containing list of clean architecture layers from
        inner layers to outer laters.  
    
  Examples
  
  To check go source code folder containing clean architecture layers:
  
  	go-clean-arch -c config.cfg path_to_src
  */
  package main