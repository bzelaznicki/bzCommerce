package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AdminCategoryRequest struct {
	Name        string        `json:"name"`
	Slug        string        `json:"slug"`
	Description string        `json:"description"`
	ParentID    uuid.NullUUID `json:"parent_id"`
}

func (cfg *apiConfig) handleApiAdminGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := cfg.db.ListCategoriesWithParent(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to load categories")
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

func (cfg *apiConfig) handleApiAdminCreateCategory(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	params := AdminCategoryRequest{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if params.Name == "" || params.Slug == "" {
		respondWithError(w, http.StatusBadRequest, "Name and slug cannot be empty")
		return
	}
	category, err := cfg.db.CreateCategory(r.Context(), database.CreateCategoryParams{
		Name:        params.Name,
		Slug:        params.Slug,
		Description: sql.NullString{String: params.Description, Valid: params.Description != ""},
		ParentID:    params.ParentID,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if pqErr.Constraint == "categories_slug_key" {
				respondWithError(w, http.StatusConflict, "Category with the same slug already exists")
			} else {
				respondWithError(w, http.StatusConflict, "Category already exists (duplicate field)")
			}
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	respondWithJSON(w, http.StatusOK, category)
}

func (cfg *apiConfig) handleApiAdminGetCategory(w http.ResponseWriter, r *http.Request) {
	categoryId, err := uuid.Parse(r.PathValue("categoryId"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := cfg.db.GetCategoryById(r.Context(), categoryId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Category not found")
		return
	}

	respondWithJSON(w, http.StatusOK, category)
}

func (cfg *apiConfig) handleApiAdminUpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryId, err := uuid.Parse(r.PathValue("categoryId"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	decoder := json.NewDecoder(r.Body)

	params := AdminCategoryRequest{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Cannot parse params")
	}
	if params.ParentID.Valid && params.ParentID.UUID == categoryId {
		respondWithError(w, http.StatusBadRequest, "A category cannot be its own parent")
		return
	}
	if params.ParentID.Valid {
		currentParentID := params.ParentID.UUID

		for {
			parent, err := cfg.db.GetCategoryById(r.Context(), currentParentID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					respondWithError(w, http.StatusBadRequest, "Parent category does not exist")
					return
				}
				respondWithError(w, http.StatusInternalServerError, "Failed to validate parent category")
				return
			}

			if parent.ID == categoryId {
				respondWithError(w, http.StatusBadRequest, "Cannot set parent: would create a circular relationship")
				return
			}

			if !parent.ParentID.Valid {
				break
			}

			currentParentID = parent.ParentID.UUID
		}
	}

	category, err := cfg.db.UpdateCategoryById(r.Context(), database.UpdateCategoryByIdParams{
		Name:        params.Name,
		Slug:        params.Slug,
		Description: sql.NullString{String: params.Description, Valid: params.Description != ""},
		ParentID:    params.ParentID,
		ID:          categoryId,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Category not found")
			return
		}
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if pqErr.Constraint == "categories_slug_key" {
				respondWithError(w, http.StatusConflict, "Category with the same slug already exists")
			} else {
				respondWithError(w, http.StatusConflict, "Category already exists (duplicate field)")
			}
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to update category")
		return
	}

	respondWithJSON(w, http.StatusOK, category)
}

func (cfg *apiConfig) handleApiAdminDeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryId, err := uuid.Parse(r.PathValue("categoryId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	rows, err := cfg.db.DeleteCategoryById(r.Context(), categoryId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23503" {
				respondWithError(w, http.StatusConflict, "Category not empty. Remove products first.")
				return
			}
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	if rows == 0 {
		respondWithError(w, http.StatusNotFound, "Category not found")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
