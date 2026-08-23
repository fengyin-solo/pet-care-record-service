package handler

import (
	"net/http"

	"petsmanagement/internal/model"
	"petsmanagement/pkg/httpx"
)

func (s *Server) registerOwnerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/owners", s.createOwner)
	mux.HandleFunc("GET /api/owners", s.listOwners)
	mux.HandleFunc("GET /api/owners/{id}", s.getOwner)
	mux.HandleFunc("PUT /api/owners/{id}", s.updateOwner)
	mux.HandleFunc("DELETE /api/owners/{id}", s.deleteOwner)
}

type createOwnerRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

func (s *Server) createOwner(w http.ResponseWriter, r *http.Request) {
	var req createOwnerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	o, err := s.svc.CreateOwner(model.Owner{Name: req.Name, Phone: req.Phone, Email: req.Email, Address: req.Address})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, o)
}

func (s *Server) listOwners(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.OwnerFilter{Keyword: r.URL.Query().Get("keyword")}
	items, total, err := s.svc.ListOwners(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getOwner(w http.ResponseWriter, r *http.Request) {
	o, err := s.svc.GetOwner(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, o)
}

func (s *Server) updateOwner(w http.ResponseWriter, r *http.Request) {
	var req createOwnerRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	o, err := s.svc.UpdateOwner(r.PathValue("id"), model.Owner{Name: req.Name, Phone: req.Phone, Email: req.Email, Address: req.Address})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, o)
}

func (s *Server) deleteOwner(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteOwner(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
