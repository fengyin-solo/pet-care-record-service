package handler

import (
	"net/http"

	"petsmanagement/internal/model"
	"petsmanagement/pkg/httpx"
)

func (s *Server) registerFollowUpRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/follow-ups", s.createFollowUp)
	mux.HandleFunc("GET /api/follow-ups", s.listFollowUps)
	mux.HandleFunc("GET /api/follow-ups/{id}", s.getFollowUp)
	mux.HandleFunc("PUT /api/follow-ups/{id}", s.updateFollowUp)
	mux.HandleFunc("DELETE /api/follow-ups/{id}", s.deleteFollowUp)
	mux.HandleFunc("POST /api/follow-ups/{id}/complete", s.completeFollowUp)
	mux.HandleFunc("POST /api/follow-ups/{id}/cancel", s.cancelFollowUp)
}

type followUpRequest struct {
	PetID       string `json:"pet_id"`
	OwnerID     string `json:"owner_id"`
	Method      string `json:"method"`
	ScheduledAt string `json:"scheduled_at"`
	Note        string `json:"note"`
}

func (s *Server) createFollowUp(w http.ResponseWriter, r *http.Request) {
	var req followUpRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.CreateFollowUp(model.FollowUp{
		PetID:       req.PetID,
		OwnerID:     req.OwnerID,
		Method:      req.Method,
		ScheduledAt: parseTime(req.ScheduledAt),
		Note:        req.Note,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, f)
}

func (s *Server) listFollowUps(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FollowUpFilter{
		PetID:  r.URL.Query().Get("pet_id"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListFollowUps(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getFollowUp(w http.ResponseWriter, r *http.Request) {
	f, err := s.svc.GetFollowUp(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

func (s *Server) updateFollowUp(w http.ResponseWriter, r *http.Request) {
	var req followUpRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.UpdateFollowUp(r.PathValue("id"), model.FollowUp{
		PetID:       req.PetID,
		OwnerID:     req.OwnerID,
		Method:      req.Method,
		ScheduledAt: parseTime(req.ScheduledAt),
		Note:        req.Note,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

func (s *Server) deleteFollowUp(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteFollowUp(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) completeFollowUp(w http.ResponseWriter, r *http.Request) {
	f, err := s.svc.TransitionFollowUp(r.PathValue("id"), model.FollowUpDone)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

func (s *Server) cancelFollowUp(w http.ResponseWriter, r *http.Request) {
	f, err := s.svc.TransitionFollowUp(r.PathValue("id"), model.FollowUpCancelled)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}
