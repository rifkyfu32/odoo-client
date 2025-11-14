package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/rifkyfu32/odoo-client/helper"
	"github.com/rifkyfu32/odoo-client/model/domain"
)

/* type LoginRequest struct {
	DB       string `json:"db"`
	Login    string `json:"login"`
	Password string `json:"password"`
} */

func main() {
	odooURL := "http://localhost:8069/jsonrpc"
	dbName := "manufaktur"
	adminUser := "admin"
	adminPass := "admin"
	log.Println("Mencoba menghubungkan ke Odoo di", odooURL)
	log.Println("Langkah 1: Autentikasi...")
	authArgs := []any{
		dbName,
		adminUser,
		adminPass,
		map[string]any{},
	}

	authResult, err := helper.CallOdoo(odooURL, "common", "authenticate", authArgs)
	if err != nil {
		log.Fatalf("Gagal autentikasi: %v", err)
	}
	uidFloat, ok := authResult.(float64)
	if !ok {
		log.Fatalf("Error: Hasil autentikasi bukan angka (UID).")
	}
	uid := int(uidFloat)
	log.Printf("Autentikasi BERHASIL! UID Anda: %d\n", uid)

	/* log.Println("Langkah 2: Membuat produk baru...")

	// Siapkan data untuk produk baru.
	// Ini adalah 'map' yang berisi field dan value yang ingin di-create.
	// 'name' adalah satu-satunya field yang WAJIB.
	newProductData := map[string]any{
		"name":           "Produk Keren Buatan Go", // Nama Produk
		"list_price":     1250.50,                  // Harga Jual
		"standard_price": 700.00,                   // Harga Modal (Cost)
		"default_code":   "GO-PRD-001",             // Kode Internal
		"type":           "product",                // Tipe: 'product' (Stockable), 'consu' (Consumable), 'service' (Service)
	}

	// Siapkan argumen untuk 'execute_kw'
	// Perhatikan Arg 5: ini adalah list (slice) yang berisi SATU elemen,
	// yaitu 'newProductData' map yang kita buat di atas.
	createArgs := []any{
		dbName,            // Arg 0: Nama DB
		uid,               // Arg 1: User ID
		adminPass,         // Arg 2: Password
		"product.product", // Arg 3: Model Odoo
		"create",          // Arg 4: Method Model yang dipanggil

		// Arg 5: Argumen Positional untuk method 'create'.
		// Method 'create' mengharapkan sebuah list dari 'vals' (values).
		// Kita kirim list berisi satu 'vals'
		[]any{newProductData},

		// Arg 6: Keyword Arguments (kwargs) - bisa kosong
		map[string]any{},
	}

	// Panggil Odoo!
	createResult, err := helper.CallOdoo(odooURL, "object", "execute_kw", createArgs)
	if err != nil {
		log.Fatalf("Gagal membuat produk: %v", err)
	}

	// Jika sukses, Odoo akan mengembalikan ID dari record baru
	newIDFloat, ok := createResult.(float64)
	if !ok {
		log.Fatalf("Error: Hasil create bukan sebuah ID (angka). Malah: %T", createResult)
	}

	newID := int(newIDFloat)
	log.Printf("LANGKAH 2 BERHASIL! Produk baru telah dibuat dengan ID: %d\n", newID) */

	// ========================================================
	// LANGKAH 2: MEMBUAT MULTIPLE PRODUK (DENGAN LOOP)
	// ========================================================
	log.Println("Langkah 2: Membangun 'slice' produk dari loop...")

	// 1. SIMULASI DATA MASUK (dari API, file, dll.)
	// Ini adalah data mentah kita, dalam format Go struct yang bersih.
	sourceData := []domain.Product{
		{ProductName: "VGA Super Kencang", SKU: "GO-LOOP-VGA", Price: 3500.00, Cost: 2000.00},
		{ProductName: "Power Supply 800W", SKU: "GO-LOOP-PSU", Price: 950.00, Cost: 500.00},
		{ProductName: "Casing PC RGB", SKU: "GO-LOOP-CASE", Price: 600.00, Cost: 300.00},
	}
	log.Printf("Data sumber berisi %d item.", len(sourceData))

	// 2. BUAT SLICE KOSONG
	// Ini adalah 'slice' yang akan kita kirim ke Odoo.
	// Odoo mengharapkan `[]any`, di mana setiap elemen adalah `map[string]any`
	var productsToCreate []any

	// 3. JALANKAN LOOP
	// Kita 'loop' data sumber dan "menerjemahkannya" ke format 'map' Odoo
	for _, srcProduct := range sourceData {
		// Buat 'map' untuk SATU produk
		odooProductMap := map[string]any{
			"name":           srcProduct.ProductName,
			"list_price":     srcProduct.Price,
			"standard_price": srcProduct.Cost,
			"default_code":   srcProduct.SKU,
			"type":           "product",
		}

		// Tambahkan 'map' produk ini ke 'slice' utama kita
		productsToCreate = append(productsToCreate, odooProductMap)
	}

	log.Printf("Slice 'productsToCreate' berhasil dibuat dengan %d produk.", len(productsToCreate))

	// 4. SIAPKAN ARGUMEN & KIRIM
	// Perhatikan Arg 5: kita sekarang MASUKKAN 'slice' yang baru kita buat
	// 'productsToCreate' sudah dalam format: []interface{}{ map1, map2, map3 }
	createArgs := []interface{}{
		dbName,            // Arg 0: Nama DB
		uid,               // Arg 1: User ID
		adminPass,         // Arg 2: Password
		"product.product", // Arg 3: Model Odoo
		"create",          // Arg 4: Method Model yang dipanggil

		// Arg 5: Argumen Positional untuk method 'create'.
		// INI ADALAH SLICE YANG KITA BUAT DARI LOOP
		//productsToCreate,
		[]interface{}{productsToCreate},

		// Arg 6: Keyword Arguments (kwargs) - bisa kosong
		map[string]interface{}{},
	}

	// Panggil Odoo! (Hanya satu kali panggil)
	createResult, err := helper.CallOdoo(odooURL, "object", "execute_kw", createArgs)
	if err != nil {
		log.Fatalf("Gagal membuat multiple produk: %v", err)
	}

	// Penanganan hasil sama seperti sebelumnya
	newIDsResult, ok := createResult.([]interface{})
	if !ok {
		log.Fatalf("Error: Hasil create (multi) bukan sebuah list ID. Malah: %T", createResult)
	}

	var newIDs []string
	for _, idInterface := range newIDsResult {
		idFloat, ok := idInterface.(float64)
		if !ok {
			continue
		}
		newIDs = append(newIDs, fmt.Sprintf("%d", int(idFloat)))
	}

	log.Printf("LANGKAH 3 BERHASIL! %d produk baru (dari loop) telah dibuat dengan IDs: %s\n", len(newIDs), strings.Join(newIDs, ", "))

	log.Println("Langkah 3: Membaca 5 data produk...")
	// Untuk membaca data, kita panggil service 'object' dan method 'execute_kw'
	// 'execute_kw' adalah "pintu" untuk memanggil metode model Odoo.

	// Kita akan memanggil: model 'product.product', method 'search_read'
	searchReadArgs := []any{
		dbName,            // Arg 0: Nama DB
		uid,               // Arg 1: User ID (hasil dari Langkah 1)
		adminPass,         // Arg 2: Password
		"product.product", // Arg 3: Nama Model Odoo
		"search_read",     // Arg 4: Nama Method Model yang dipanggil

		// Arg 5: Domain (filter) - List kosong [] berarti "semua"
		[]any{},

		// Arg 6: Keyword Arguments (kwargs) - seperti 'fields' dan 'limit'
		map[string]any{
			"fields": []string{"name", "default_code", "list_price"},
			"limit":  5,
		},
	}

	searchResult, err := helper.CallOdoo(odooURL, "object", "execute_kw", searchReadArgs)
	if err != nil {
		log.Fatalf("Gagal membaca data produk: %v", err)
	}

	// Hasil 'search_read' adalah sebuah List (slice) dari Map
	// Kita perlu type assertion lagi
	products, ok := searchResult.([]any)
	if !ok {
		log.Fatalf("Error: Hasil search_read bukan sebuah list.")
	}

	log.Printf("Berhasil mengambil %d data produk:\n", len(products))
	log.Println("-------------------------------------------")

	if len(products) == 0 {
		log.Println("Tidak ada produk ditemukan.")
		log.Println("CATATAN: Jika Anda menginstal TANPA DATA DEMO, ini normal.")
	}

	// Iterasi (loop) hasil dan cetak ke konsol
	for i, p := range products {
		// Setiap 'p' adalah 'map[string]any'
		product, ok := p.(map[string]any)
		if !ok {
			log.Printf("Error: Item %d bukan map data produk.", i)
			continue
		}

		// Ambil data dari map. Hati-hati dengan 'false' vs 'nil'
		nama := product["name"].(string)

		// 'default_code' (Kode Internal) bisa jadi tidak diisi (false/nil)
		kode, ok := product["default_code"].(string)
		if !ok || kode == "" {
			kode = "N/A"
		}

		harga := product["list_price"].(float64)

		fmt.Printf("  Produk %d:\n", i+1)
		fmt.Printf("    Nama  : %s\n", nama)
		fmt.Printf("    Kode  : %s\n", kode)
		fmt.Printf("    Harga : %.2f\n", harga)
		fmt.Println("-------------------------------------------")
	}

}
