/*
Go-clean-arch is a linter for enforcing clean architecture principles in Go.

The linter reports violations of clean architecture dependency rules by checking go imports
within packages against user defined dependency rules.

Usage
  go-clean-arch [flags] [path]

The flags are:
  -c  Config file containing list of clean architecture layers with hierarchy
      index from inner layers to outer laters.

Examples

To check go source code folder containing clean architecture layers:

  go-clean-arch -c config.json path

The go-clean-arch linter can be run on the Git repo https://github.com/ManuelKiessling/go-cleanarchitecture
by cloning the repo using the command;

  git clone https://github.com/ManuelKiessling/go-cleanarchitecture

Run the linter with the command:

  go-clean-arch -c layers.json path-to-repo\go-cleanarchitecture\src

where the layers config file contains the clean architecture layers:

  {
    "domain": 1,
    "usecases": 2,
    "interfaces": 3,
    "infrastructure": 4
  }

*/
package main
