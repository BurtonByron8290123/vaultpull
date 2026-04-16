// Package multienv provides support for writing Vault secrets into multiple
// environment-specific .env files in a single pull operation.
//
// # Overview
//
// When a project needs different subsets of secrets for different environments
// (e.g. .env.dev, .env.staging, .env.prod), multienv allows declaring those
// targets in a YAML configuration file:
//
//	targets:
//	  - name: dev
//	    path: .env.dev
//	    keys: [DB_HOST, DB_PORT]
//	  - name: prod
//	    path: .env.prod
//
// Each target can optionally restrict which keys are written via the "keys"
// list. An empty "keys" list means all fetched secrets are written.
//
// # Usage
//
//	cfg, err := multienv.LoadConfig("multienv.yaml")
//	w := multienv.New(".", cfg.Targets)
//	err = w.WriteAll(secrets)
//
// All output files are written with permission 0600.
package multienv
