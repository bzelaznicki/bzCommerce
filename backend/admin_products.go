package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

type AdminProductsListPageData struct {
	Products    []AdminProductRow
	Breadcrumbs []Breadcrumb
}
type AdminProductRow struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	ImagePath    string    `json:"image_path"`
	CategoryName string    `json:"category_name"`
	CategorySlug string    `json:"category_slug"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handleAdminProductList(w http.ResponseWriter, r *http.Request) {
	rows, err := cfg.db.ListProductsWithCategory(r.Context())
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not load products")
		log.Printf("could not load products: %v", err)
		return
	}

	products := make([]AdminProductRow, 0, len(rows))
	for _, row := range rows {
		products = append(products, AdminProductRow{
			ID:           row.ID,
			Name:         row.Name,
			Slug:         row.Slug,
			Description:  row.Description.String,
			ImagePath:    row.ImageUrl.String,
			CategoryName: row.CategoryName,
			CategorySlug: row.CategorySlug,
		})
	}

	data := AdminProductsListPageData{
		Products: products,
		Breadcrumbs: NewBreadcrumbTrail(
			Breadcrumb{Label: "Products"},
		),
	}

	cfg.Render(w, r, "templates/pages/admin_products.html", data)
}

func (cfg *apiConfig) handleAdminProductNewForm(w http.ResponseWriter, r *http.Request) {

	categories, err := cfg.db.GetCategories(r.Context())
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to load categories")
		log.Printf("failed to load categories: %v", err)
		return
	}

	groups := groupCategories(categories)
	options := flattenCategoryOptions(groups)

	data := struct {
		IsEdit          bool
		FormAction      string
		CategoryOptions []CategoryOption
	}{
		IsEdit:          false,
		FormAction:      "/admin/products",
		CategoryOptions: options,
	}

	cfg.Render(w, r, "templates/pages/admin_product_new.html", data)
}

func (cfg *apiConfig) handleAdminProductCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid form")
		log.Printf("error parsing form: %v", err)
		return
	}

	name := r.FormValue("name")
	slug := r.FormValue("slug")
	description := r.FormValue("description")
	imageURL := r.FormValue("image_url")
	categoryID := r.FormValue("category_id")
	if categoryID == "" {
		cfg.RenderError(w, r, http.StatusBadRequest, "Category is required")
		log.Println("category_id is empty")
		return
	}

	product, err := cfg.db.CreateProduct(r.Context(), database.CreateProductParams{
		Name:        name,
		Slug:        slug,
		Description: sql.NullString{String: description, Valid: description != ""},
		ImageUrl:    sql.NullString{String: imageURL, Valid: imageURL != ""},
		CategoryID:  uuid.MustParse(categoryID),
	})
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Error creating product")
		log.Printf("failed to create product: %v", err)
		return
	}

	variantName := r.FormValue("variant_name")
	variantPrice, err := strconv.Atoi(r.FormValue("variant_price"))
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Error converting price")
	}
	variantSku := r.FormValue("variant_sku")
	variantStock64, err := strconv.ParseInt(r.FormValue("variant_stock"), 10, 32)
	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid stock quantity")
		log.Printf("invalid stock quantity: %v", err)
		return
	}
	variantImage := r.FormValue("variant_image")

	_, err = cfg.db.CreateProductVariant(r.Context(), database.CreateProductVariantParams{
		ProductID:     product.ID,
		Name:          sql.NullString{String: variantName, Valid: variantName != ""},
		Sku:           variantSku,
		Price:         float64(variantPrice),
		StockQuantity: int32(variantStock64),
		ImageUrl:      sql.NullString{String: variantImage, Valid: variantImage != ""},
	})
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to create product variant")
		log.Printf("failed to create product variant: %v", err)
		return
	}

	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

func (cfg *apiConfig) handleAdminProductEditForm(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Invalid product ID")
		log.Printf("invalid product id: %v", err)
		return
	}

	p, err := cfg.db.GetProductById(r.Context(), id)
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Product not found")
		log.Printf("product not found (id %s): %v", id, err)
		return
	}

	categories, err := cfg.db.GetCategories(r.Context())
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to load categories")
		log.Printf("failed to load categories for product edit: %v", err)
		return
	}

	groups := groupCategories(categories)
	options := flattenCategoryOptions(groups)

	data := struct {
		IsEdit          bool
		FormAction      string
		Product         Product
		CategoryOptions []CategoryOption
	}{
		IsEdit:     true,
		FormAction: fmt.Sprintf("/admin/products/%s", id),
		Product: Product{
			ID:          p.ID,
			Name:        p.Name,
			Slug:        p.Slug,
			Description: p.Description.String,
			ImagePath:   p.ImageUrl.String,
			CategoryID:  p.CategoryID,
		},
		CategoryOptions: options,
	}

	cfg.Render(w, r, "templates/pages/admin_product_edit.html", data)
}

func (cfg *apiConfig) handleAdminProductUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid product ID")
		log.Printf("invalid product id: %v", err)
		return
	}

	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid form")
		log.Printf("error parsing form: %v", err)
		return
	}

	name := r.FormValue("name")
	slug := r.FormValue("slug")
	description := r.FormValue("description")
	imageURL := r.FormValue("image_url")
	categoryID := r.FormValue("category_id")

	err = cfg.db.UpdateProduct(r.Context(), database.UpdateProductParams{
		ID:          id,
		Name:        name,
		Slug:        slug,
		Description: sql.NullString{String: description, Valid: description != ""},
		ImageUrl:    sql.NullString{String: imageURL, Valid: imageURL != ""},
		CategoryID:  uuid.MustParse(categoryID),
	})
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to update product")
		log.Printf("failed to update product (id %s): %v", id, err)
		return
	}

	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

func (cfg *apiConfig) handleAdminProductDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Invalid product ID")
		log.Printf("invalid product id: %v", err)
		return
	}

	err = cfg.db.DeleteProduct(r.Context(), id)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to delete product")
		log.Printf("failed to delete product (id %s): %v", id, err)
		return
	}

	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

func (cfg *apiConfig) handleAdminVariantList(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Invalid product ID")
		log.Printf("invalid product id for variant list: %v", err)
		return
	}

	product, err := cfg.db.GetProductById(r.Context(), productID)
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Product not found")
		log.Printf("product not found (id %s): %v", productID, err)
		return
	}

	variants, err := cfg.db.GetVariantsByProductID(r.Context(), productID)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to load variants")
		log.Printf("failed to load variants for product (id %s): %v", productID, err)
		return
	}

	data := struct {
		Product  database.Product
		Variants []database.ProductVariant
	}{
		Product:  product,
		Variants: variants,
	}

	pageData, err := cfg.NewPageData(r.Context(), data)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Error rendering")
		log.Printf("failed to create page data: %v", err)
		return
	}

	tmpl := parsePageTemplate("templates/pages/admin_variant_list.html")
	if err := tmpl.ExecuteTemplate(w, "base.html", pageData); err != nil {
		log.Printf("failed to execute variant list template: %v", err)
	}
}

func (cfg *apiConfig) handleAdminVariantNewForm(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Invalid product ID")
		log.Printf("invalid product id for new variant: %v", err)
		return
	}

	data := struct {
		IsEdit     bool
		FormAction string
		Variant    struct {
			Sku           string
			Price         int32
			StockQuantity int32
			ImageUrl      string
			VariantName   string
		}
	}{
		IsEdit:     false,
		FormAction: fmt.Sprintf("/admin/products/%s/variants", productID),
	}

	cfg.Render(w, r, "templates/pages/admin_variant_form.html", data)
}

func (cfg *apiConfig) handleAdminVariantCreate(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Invalid product ID")
		log.Printf("invalid product id for creating variant: %v", err)
		return
	}
	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid form")
		log.Printf("error parsing variant creation form: %v", err)
		return
	}

	sku := r.FormValue("sku")
	price, err := strconv.Atoi(r.FormValue("price"))
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Error converting price")
	}
	stockQty64, err := strconv.ParseInt(r.FormValue("stock_quantity"), 10, 32)
	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid stock quantity")
		log.Printf("invalid stock quantity in variant create: %v", err)
		return
	}
	imageURL := r.FormValue("image_url")
	variantName := r.FormValue("variant_name")

	_, err = cfg.db.CreateVariant(r.Context(), database.CreateVariantParams{
		ProductID:     productID,
		Sku:           sku,
		Price:         float64(price),
		StockQuantity: int32(stockQty64),
		ImageUrl:      sql.NullString{String: imageURL, Valid: imageURL != ""},
		VariantName:   sql.NullString{String: variantName, Valid: variantName != ""},
	})
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to create variant")
		log.Printf("failed to create variant for product (id %s): %v", productID, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/products/%s/variants", productID), http.StatusSeeOther)
}

func (cfg *apiConfig) handleAdminVariantEditForm(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Invalid variant ID")
		log.Printf("invalid variant id: %v", err)
		return
	}

	v, err := cfg.db.GetVariantByID(r.Context(), id)
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Variant not found")
		log.Printf("variant not found (id %s): %v", id, err)
		return
	}

	data := struct {
		IsEdit     bool
		FormAction string
		Variant    struct {
			Sku           string
			Price         float64
			StockQuantity int32
			ImageUrl      string
			VariantName   string
		}
	}{
		IsEdit:     true,
		FormAction: fmt.Sprintf("/admin/variants/%s", id),
		Variant: struct {
			Sku           string
			Price         float64
			StockQuantity int32
			ImageUrl      string
			VariantName   string
		}{
			Sku:           v.Sku,
			Price:         v.Price,
			StockQuantity: v.StockQuantity,
			ImageUrl:      v.ImageUrl.String,
			VariantName:   v.VariantName.String,
		},
	}

	cfg.Render(w, r, "templates/pages/admin_variant_form.html", data)
}

func (cfg *apiConfig) handleAdminVariantUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Invalid variant ID")
		log.Printf("invalid variant id: %v", err)
		return
	}
	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid form")
		log.Printf("error parsing variant update form: %v", err)
		return
	}

	sku := r.FormValue("sku")
	price, err := strconv.Atoi(r.FormValue("price"))
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Error converting price")
	}
	stockQty64, err := strconv.ParseInt(r.FormValue("stock_quantity"), 10, 32)
	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid stock quantity")
		log.Printf("invalid stock quantity in variant update: %v", err)
		return
	}
	imageURL := r.FormValue("image_url")
	variantName := r.FormValue("variant_name")

	err = cfg.db.UpdateVariant(r.Context(), database.UpdateVariantParams{
		ID:            id,
		Sku:           sku,
		Price:         float64(price),
		StockQuantity: int32(stockQty64),
		ImageUrl:      sql.NullString{String: imageURL, Valid: imageURL != ""},
		VariantName:   sql.NullString{String: variantName, Valid: variantName != ""},
	})
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to update variant")
		log.Printf("failed to update variant (id %s): %v", id, err)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func (cfg *apiConfig) handleAdminVariantDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Invalid variant ID")
		log.Printf("invalid variant id: %v", err)
		return
	}

	variant, err := cfg.db.GetVariantByID(r.Context(), id)
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Variant not found")
		log.Printf("variant not found (id %s): %v", id, err)
		return
	}

	err = cfg.db.DeleteVariant(r.Context(), id)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to delete variant")
		log.Printf("failed to delete variant (id %s): %v", id, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/products/%s/variants", variant.ProductID), http.StatusSeeOther)
}
