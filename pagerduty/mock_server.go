/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package pagerduty

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/blend/go-sdk/webutil"
)

var (
	_ http.Handler = (*MockAPI)(nil)
)

// MockAPI implements methods that can be called with the client.
type MockAPI struct {
	ListIncidents  func() ListIncidentsOutput
	UpdateIncident func(string, UpdateIncidentInput) Incident
	CreateIncident func(CreateIncidentInput) Incident
}

// Handler implements http.Handler.
func (ma MockAPI) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		if req.URL.Path == "/incidents" {
			if ma.ListIncidents != nil {
				_ = webutil.WriteJSON(rw, http.StatusOK, ma.ListIncidents())
				return
			}
		}
	}
	if req.Method == http.MethodPut {
		if strings.HasPrefix(req.URL.Path, "/incidents/") {
			if ma.UpdateIncident != nil {
				incidentID := strings.TrimPrefix(req.URL.Path, "/incidents/")
				var incidentBody updateIncidentInputWrapper
				if err := json.NewDecoder(req.Body).Decode(&incidentBody); err != nil {
					http.Error(rw, err.Error(), http.StatusBadRequest)
					return
				}
				incident := ma.UpdateIncident(incidentID, incidentBody.Incident)
				_ = webutil.WriteJSON(rw, http.StatusOK, updateIncidentOutputWrapper{Incident: incident})
				return
			}
		}
	}
	if req.Method == http.MethodPost {
		if req.URL.Path == "/incidents" {
			var incidentBody createIncidentInputWrapper
			if err := json.NewDecoder(req.Body).Decode(&incidentBody); err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				return
			}
			incident := ma.CreateIncident(incidentBody.Incident)
			_ = webutil.WriteJSON(rw, http.StatusOK, updateIncidentOutputWrapper{Incident: incident})
			return
		}
	}
	http.NotFound(rw, req)
}
