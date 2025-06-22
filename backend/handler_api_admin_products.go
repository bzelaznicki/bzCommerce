package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// handleApiAdminGetProducts handles the admin API endpoint for retrieving a paginated list of products.
// It supports optional search, sorting, and pagination parameters via query string:
//   - "search": filters products by name (case-insensitive substring match)
//   - "sort_by": specifies the column to sort by ("name", "category_name", "updated_at", or defaults to "created_at")
//   - "sort_order": specifies the sort direction ("asc" or "desc")
//   - Pagination is controlled by page and limit parameters (see getPaginationParams).
//
// The response includes a paginated list of products with their category information.
// Returns a JSON response with the products and pagination metadata, or an error if the query fails.
func (cfg *apiConfig) handleApiAdminGetProducts(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationParams(r)
	offset := (page - 1) * limit

	q := r.URL.Query()
	search := q.Get("search")
	sortBy := q.Get("sort_by")
	sortOrder := q.Get("sort_order")

	sortColumn := "p.created_at"
	switch sortBy {
	case "name":
		sortColumn = "p.name"
	case "category_name":
		sortColumn = "c.name"
	case "updated_at":
		sortColumn = "p.updated_at"
	}

	direction := "ASC"
	if sortOrder == "desc" {
		direction = "DESC"
	}
	// #nosec G201 -- sortColumn and direction are strictly whitelisted to prevent SQL injection
	query := fmt.Sprintf(`
		SELECT
		p.id, p.name, p.slug, p.description, p.image_url AS image_path,
		c.name AS category_name, c.slug AS category_slug,
		p.created_at, p.updated_at
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE p.name ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3`, sortColumn, direction)

	rows, err := cfg.sqlDB.QueryContext(r.Context(), query, "%"+search+"%", limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Query failed")
		return
	}
	defer rows.Close()

	products := []AdminProductRow{}
	for rows.Next() {
		var p AdminProductRow
		var desc, image sql.NullString
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Slug, &desc, &image,
			&p.CategoryName, &p.CategorySlug,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {

			respondWithError(w, http.StatusInternalServerError, "Scan failed")
			return
		}
		p.Description = desc.String
		p.ImagePath = image.String
		products = append(products, p)
	}

	count, err := cfg.db.CountFilteredProducts(r.Context(), "%"+search+"%")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Count failed")
		return
	}

	if products == nil {
		products = []AdminProductRow{}
	}

	resp := NewPaginatedResponse(products, page, limit, count)
	respondWithJSON(w, http.StatusOK, resp)
}

// handleApiCreateProduct handles the creation of a new product and its initial variant via the admin API.
// It expects a JSON payload containing product details and a product variant. The function validates the input,
// starts a database transaction, creates the product and its variant, and commits the transaction. If any step fails,
// it responds with an appropriate error message. On success, it returns the created product and variant as JSON.
//
// Expected JSON request body:
//
//	{
//	  "name": "Product Name",
//	  "slug": "product-slug",
//	  "description": "Product description",
//	  "image_url": "https://example.com/image.jpg",
//	  "category_id": "uuid-string",
//	  "product_variant": {
//	    "sku": "SKU123",
//	    "price": 19.99,
//	    "stock_quantity": 100,
//	    "image_url": "https://example.com/variant.jpg",
//	    "name": "Size M"
//	  }
//	}
//
// Responses:
//
//	200 OK: Returns the created product and variant as JSON.
//	400 Bad Request: If required fields are missing or variant creation fails.
//	500 Internal Server Error: On decoding, transaction, or product creation errors.
//
// Parameters:
//   - w: http.ResponseWriter to write the response.
//   - r: *http.Request containing the request data.
//
// The ProductVariant structure includes fields: id, product_id, sku, price, stock_quantity, image_url, variant_name, created_at, updated_at.
func (cfg *apiConfig) handleApiCreateProduct(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name        string    `json:"name"`
		Slug        string    `json:"slug"`
		Description string    `json:"description"`
		ImageURL    string    `json:"image_url"`
		CategoryID  uuid.UUID `json:"category_id"`
		Variant     struct {
			Sku           string  `json:"sku"`
			Price         float64 `json:"price"`
			StockQuantity int32   `json:"stock_quantity"`
			ImageUrl      string  `json:"image_url"`
			Name          string  `json:"name"`
		} `json:"product_variant"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if params.Name == "" || params.Slug == "" || params.CategoryID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Name, slug and category_id cannot be empty")
		return
	}
	tx, err := cfg.sqlDB.BeginTx(r.Context(), nil)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer tx.Rollback()

	qtx := cfg.db.WithTx(tx)
	product, err := qtx.CreateProduct(r.Context(), database.CreateProductParams{
		Name:        params.Name,
		Slug:        params.Slug,
		Description: sql.NullString{Valid: params.Description != "", String: params.Description},
		ImageUrl:    sql.NullString{Valid: params.ImageURL != "", String: params.ImageURL},
		CategoryID:  params.CategoryID,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if pqErr.Constraint == "products_slug_key" {
				respondWithError(w, http.StatusConflict, "Product with the same slug already exists")
			} else {
				respondWithError(w, http.StatusConflict, "Product already exists (duplicate field)")
			}
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	newVariant := database.CreateProductVariantParams{
		ProductID:     product.ID,
		Name:          sql.NullString{Valid: params.Variant.Name != "", String: params.Variant.Name},
		Sku:           params.Variant.Sku,
		Price:         params.Variant.Price,
		StockQuantity: params.Variant.StockQuantity,
		ImageUrl:      sql.NullString{Valid: params.Variant.ImageUrl != "", String: params.Variant.ImageUrl},
	}

	variant, err := qtx.CreateProductVariant(r.Context(), newVariant)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if pqErr.Constraint == "product_variants_sku_key" {
				respondWithError(w, http.StatusConflict, "Variant with the same SKU already exists")
			} else {
				respondWithError(w, http.StatusConflict, "Product variant already exists (duplicate field)")
			}
			return
		}
		respondWithError(w, http.StatusBadRequest, "Failed to create variant")
		return
	}

	err = tx.Commit()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	type response struct {
		ID          uuid.UUID               `json:"id"`
		CategoryID  uuid.UUID               `json:"category_id"`
		Name        string                  `json:"name"`
		Slug        string                  `json:"slug"`
		ImageUrl    sql.NullString          `json:"image_url"`
		Description sql.NullString          `json:"description"`
		CreatedAt   time.Time               `json:"created_at"`
		UpdatedAt   time.Time               `json:"updated_at"`
		Variant     database.ProductVariant `json:"product_variant"`
	}

	resp := response{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Slug:        product.Slug,
		ImageUrl:    product.ImageUrl,
		Description: product.Description,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		Variant:     variant,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handleApiUpdateProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	type parameters struct {
		Name        string    `json:"name"`
		Slug        string    `json:"slug"`
		Description string    `json:"description"`
		ImageURL    string    `json:"image_url"`
		CategoryID  uuid.UUID `json:"category_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if params.Name == "" || params.Slug == "" || params.CategoryID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Name, slug, and category_id cannot be empty")
		return
	}

	updatedProduct, err := cfg.db.UpdateProduct(r.Context(), database.UpdateProductParams{
		ID:          productID,
		Name:        params.Name,
		Slug:        params.Slug,
		Description: sql.NullString{Valid: params.Description != "", String: params.Description},
		ImageUrl:    sql.NullString{Valid: params.ImageURL != "", String: params.ImageURL},
		CategoryID:  params.CategoryID,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if pqErr.Constraint == "products_slug_key" {
				respondWithError(w, http.StatusConflict, "Product with the same slug already exists")
			} else {
				respondWithError(w, http.StatusConflict, "Product already exists (duplicate field)")
			}
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	respondWithJSON(w, http.StatusOK, updatedProduct)
}

func (cfg *apiConfig) handleApiAdminGetProductDetails(w http.ResponseWriter, r *http.Request) {
	productId, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := cfg.db.GetProductById(r.Context(), productId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	respondWithJSON(w, http.StatusOK, product)
}

func (cfg *apiConfig) handleApiAdminDeleteProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	rows, err := cfg.db.DeleteProduct(r.Context(), productId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	if rows == 0 {
		respondWithError(w, http.StatusNotFound, "Product not found")
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
