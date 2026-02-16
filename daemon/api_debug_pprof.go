// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2015-2019 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package daemon

import (
	"net/http"
	"net/http/pprof"

	"github.com/snapcore/snapd/overlord/auth"
)

var debugPprofCmd = &Command{
	Path:       "/v2/debug/pprof/{profile...}",
	GET:        getPprof,
	ReadAccess: rootAccess{},
}

func getPprof(c *Command, r *http.Request, user *auth.UserState) Response {
	profile := r.PathValue("profile")
	switch profile {
	case "cmdline":
		return pprofHandlerFunc(pprof.Cmdline)
	case "profile":
		return pprofHandlerFunc(pprof.Profile)
	case "symbol":
		return pprofHandlerFunc(pprof.Symbol)
	case "trace":
		return pprofHandlerFunc(pprof.Trace)
	case "heap", "allocs", "block", "threadcreate", "goroutine", "mutex":
		return &pprofHandler{handler: pprof.Handler(profile)}
	default:
		return NotFound("unknown pprof profile")
	}
}

type pprofHandlerFunc func(http.ResponseWriter, *http.Request)

func (h pprofHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}

type pprofHandler struct {
	handler http.Handler
}

func (h *pprofHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
}
