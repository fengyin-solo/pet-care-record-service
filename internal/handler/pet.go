package handler

import (
	"net/http"
	"time"

	"petsmanagement/internal/model"
	"petsmanagement/pkg/httpx"
)

func (s *Server) registerPetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/pets", s.createPet)
	mux.HandleFunc("GET /api/pets", s.listPets)
	mux.HandleFunc("GET /api/pets/{id}", s.getPet)
	mux.HandleFunc("PUT /api/pets/{id}", s.updatePet)
	mux.HandleFunc("DELETE /api/pets/{id}", s.deletePet)
	mux.HandleFunc("POST /api/pets/{id}/archive", s.archivePet)
	mux.HandleFunc("POST /api/pets/batch-archive", s.batchArchivePets)
}

type createPetRequest struct {
	OwnerID      string `json:"owner_id"`
	Name         string `json:"name"`
	Species      string `json:"species"`
	Breed        string `json:"breed"`
	Gender       string `json:"gender"`
	BirthDate    string `json:"birth_date"`
	HealthStatus string `json:"health_status"`
}

func (s *Server) createPet(w http.ResponseWriter, r *http.Request) {
	var req createPetRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreatePet(model.Pet{
		OwnerID:      req.OwnerID,
		Name:         req.Name,
		Species:      req.Species,
		Breed:        req.Breed,
		Gender:       req.Gender,
		BirthDate:    parseTime(req.BirthDate),
		HealthStatus: req.HealthStatus,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listPets(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.PetFilter{
		OwnerID:      r.URL.Query().Get("owner_id"),
		Species:      r.URL.Query().Get("species"),
		HealthStatus: r.URL.Query().Get("health_status"),
		Status:       r.URL.Query().Get("status"),
		Keyword:      r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListPets(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getPet(w http.ResponseWriter, r *http.Request) {
	p, err := s.svc.GetPet(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) updatePet(w http.ResponseWriter, r *http.Request) {
	var req createPetRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdatePet(r.PathValue("id"), model.Pet{
		OwnerID:      req.OwnerID,
		Name:         req.Name,
		Species:      req.Species,
		Breed:        req.Breed,
		Gender:       req.Gender,
		BirthDate:    parseTime(req.BirthDate),
		HealthStatus: req.HealthStatus,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) deletePet(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeletePet(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) archivePet(w http.ResponseWriter, r *http.Request) {
	p, err := s.svc.ArchivePet(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

type batchArchivePetsRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchArchivePets(w http.ResponseWriter, r *http.Request) {
	var req batchArchivePetsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	count, err := s.svc.BatchArchivePets(req.IDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"archived": count})
}

// parseTime 解析 ISO8601 日期字符串，失败返回零值。
func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			return time.Time{}
		}
	}
	return t
}
