// Package classification GO2HAL API.
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Marc Arndt<marc@marcarndt.com> http://www.marcarndt.com
//     Title: go2hal API
//
//     basePath: /api
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	go2hal2 "github.com/weAutomateEverything/go2hal/go2hal"
)

//go:generate swagger generate spec

func main() {
	go2hal := go2hal2.NewGo2Hal()
	go2hal.Logger.Log("terminated", <-go2hal.Err)
}
