package handler

import (
	"net/http"

	"petsmanagement/internal/model"
	"petsmanagement/pkg/httpx"
)

func (s *Server) registerMedicalRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/medical-records", s.createMedicalRecord)
	mux.HandleFunc("GET /api/medical-records", s.listMedicalRecords)
	mux.HandleFunc("GET /api/medical-records/{id}", s.getMedicalRecord)
	mux.HandleFunc("PUT /api/medical-records/{id}", s.updateMedicalRecord)
	mux.HandleFunc("DELETE /api/medical-records/{id}", s.deleteMedicalRecord)
}

type medicalRecordRequest struct {
	PetID          string `json:"pet_id"`
	VetName        string `json:"vet_name"`
	Clinic         string `json:"clinic"`
	Diagnosis      string `json:"diagnosis"`
	Treatment      string `json:"treatment"`
	VisitAt        string `json:"visit_at"`
	FollowUpNeeded bool   `json:"follow_up_needed"`
}

func (s *Server) createMedicalRecord(w http.ResponseWriter, r *http.Request) {
	var req medicalRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.CreateMedicalRecord(model.MedicalRecord{
		PetID:          req.PetID,
		VetName:        req.VetName,
		Clinic:         req.Clinic,
		Diagnosis:      req.Diagnosis,
		Treatment:      req.Treatment,
		VisitAt:        parseTime(req.VisitAt),
		FollowUpNeeded: req.FollowUpNeeded,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, m)
}

func (s *Server) listMedicalRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MedicalRecordFilter{
		PetID:  r.URL.Query().Get("pet_id"),
		Clinic: r.URL.Query().Get("clinic"),
	}
	items, total, err := s.svc.ListMedicalRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMedicalRecord(w http.ResponseWriter, r *http.Request) {
	m, err := s.svc.GetMedicalRecord(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) updateMedicalRecord(w http.ResponseWriter, r *http.Request) {
	var req medicalRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.UpdateMedicalRecord(r.PathValue("id"), model.MedicalRecord{
		PetID:          req.PetID,
		VetName:        req.VetName,
		Clinic:         req.Clinic,
		Diagnosis:      req.Diagnosis,
		Treatment:      req.Treatment,
		VisitAt:        parseTime(req.VisitAt),
		FollowUpNeeded: req.FollowUpNeeded,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) deleteMedicalRecord(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteMedicalRecord(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
