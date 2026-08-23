package handler

import (
	"net/http"
	"strconv"

	"petsmanagement/internal/model"
	"petsmanagement/pkg/httpx"
)

func (s *Server) registerVaccineRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/vaccines", s.createVaccine)
	mux.HandleFunc("GET /api/vaccines", s.listVaccines)
	mux.HandleFunc("GET /api/vaccines/due", s.listDueVaccines)
	mux.HandleFunc("GET /api/vaccines/{id}", s.getVaccine)
	mux.HandleFunc("PUT /api/vaccines/{id}", s.updateVaccine)
	mux.HandleFunc("DELETE /api/vaccines/{id}", s.deleteVaccine)
}

type vaccineRequest struct {
	PetID          string `json:"pet_id"`
	VaccineName    string `json:"vaccine_name"`
	BatchNo        string `json:"batch_no"`
	AdministeredAt string `json:"administered_at"`
	NextDueAt      string `json:"next_due_at"`
	VetName        string `json:"vet_name"`
}

func (s *Server) createVaccine(w http.ResponseWriter, r *http.Request) {
	var req vaccineRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	v, err := s.svc.CreateVaccine(model.Vaccine{
		PetID:          req.PetID,
		VaccineName:    req.VaccineName,
		BatchNo:        req.BatchNo,
		AdministeredAt: parseTime(req.AdministeredAt),
		NextDueAt:      parseTime(req.NextDueAt),
		VetName:        req.VetName,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, v)
}

func (s *Server) listVaccines(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.VaccineFilter{PetID: r.URL.Query().Get("pet_id")}
	items, total, err := s.svc.ListVaccines(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) listDueVaccines(w http.ResponseWriter, r *http.Request) {
	days := 30
	if v := r.URL.Query().Get("days"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			days = n
		}
	}
	items, err := s.svc.ListDueVaccines(days)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}

func (s *Server) getVaccine(w http.ResponseWriter, r *http.Request) {
	v, err := s.svc.GetVaccine(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

func (s *Server) updateVaccine(w http.ResponseWriter, r *http.Request) {
	var req vaccineRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	v, err := s.svc.UpdateVaccine(r.PathValue("id"), model.Vaccine{
		PetID:          req.PetID,
		VaccineName:    req.VaccineName,
		BatchNo:        req.BatchNo,
		AdministeredAt: parseTime(req.AdministeredAt),
		NextDueAt:      parseTime(req.NextDueAt),
		VetName:        req.VetName,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

func (s *Server) deleteVaccine(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteVaccine(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
