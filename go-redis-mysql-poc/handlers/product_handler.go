package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go-redis-mysql-poc/database"
	"go-redis-mysql-poc/models"
	"go-redis-mysql-poc/redis"

	"github.com/gin-gonic/gin"
)

// CreateProduct creates a new product.
// func CreateProduct(w http.ResponseWriter, r *http.Request) {
func CreateProduct(c *gin.Context) {
	var product models.Product
	//err := json.NewDecoder(r.Body).Decode(&product)
	/*if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}*/

	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Exec("INSERT INTO products (name, price) VALUES (?, ?)", product.Name, product.Price)
	/*if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}*/
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	/*if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}*/

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	product.ID = int(id) // Set the ID after insertion

	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(product)
	c.JSON(http.StatusCreated, product)
}

// GetProduct retrieves a product by ID.
// func GetProduct(w http.ResponseWriter, r *http.Request) {
func GetProduct(c *gin.Context) {
	//vars := mux.Vars(r)
	//idStr := vars["id"]
	//id, err := strconv.Atoi(idStr)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	/*if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}*/

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Try to get from Redis first
	cacheKey := fmt.Sprintf("product:%d", id)
	cachedProduct, err := redis.Get(cacheKey)
	if err == nil && cachedProduct != "" {
		// Found in cache
		log.Println("Fetching product from Redis")
		var product models.Product
		err := json.Unmarshal([]byte(cachedProduct), &product)
		if err != nil {
			log.Printf("Error unmarshaling from cache: %v", err) // Log the error, but continue to DB
		} else {
			//w.Header().Set("Content-Type", "application/json")
			//json.NewEncoder(w).Encode(product)
			c.JSON(http.StatusOK, product)
			return
		}
	}

	var product models.Product
	err = database.DB.Get(&product, "SELECT id, name, price FROM products WHERE id = ?", id)
	/*if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}*/
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	log.Println("Fetching product from DB")

	// Store in Redis
	productJSON, err := json.Marshal(product)
	if err != nil {
		log.Printf("Error marshaling to JSON: %v", err) // Log the error, but continue
	} else {
		err = redis.Set(cacheKey, productJSON, 3600) // Cache for 1 hour
		if err != nil {
			log.Printf("Error setting cache: %v", err)
		}
	}

	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(product)
	c.JSON(http.StatusOK, product)
}

// UpdateProduct updates an existing product.
// func UpdateProduct(w http.ResponseWriter, r *http.Request) {
func UpdateProduct(c *gin.Context) {
	//vars := mux.Vars(r)
	//idStr := vars["id"]
	//id, err := strconv.Atoi(idStr)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	/*if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}*/

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	/*err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}*/

	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = database.DB.Exec("UPDATE products SET name = ?, price = ? WHERE id = ?", product.Name, product.Price, id)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("product:%d", id)
	err = redis.Delete(cacheKey)
	if err != nil {
		log.Printf("Error deleting cache: %v", err)
	}

	product.ID = id // Set the ID for the response

	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(product)
	c.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product by ID.
// func DeleteProduct(w http.ResponseWriter, r *http.Request) {
func DeleteProduct(c *gin.Context) {
	//vars := mux.Vars(r)
	//idStr := vars["id"]
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("product:%d", id)
	err = redis.Delete(cacheKey)
	if err != nil {
		log.Printf("Error deleting cache: %v", err)
	}

	//w.WriteHeader(http.StatusNoContent) // 204 No Content
	c.Status(http.StatusNoContent) // 204 No Content

}

// ListProducts retrieves all products.
// func ListProducts(w http.ResponseWriter, r *http.Request) {
func ListProducts(c *gin.Context) {
	var products []models.Product
	err := database.DB.Select(&products, "SELECT id, name, price FROM products")
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(products)
	c.JSON(http.StatusOK, products)

}
