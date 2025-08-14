package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Product struct {
    ID          int
    Name        string
    Price       int
    Image       string
    Description string
}

var products = []Product{
    {1, "Kaos Golang", 100000, "/static/kaos.jpg", "Kaos nyaman berbahan katun dengan logo Golang."},
    {2, "Topi Programmer", 75000, "/static/topi.jpg", "Topi keren untuk para programmer sejati."},
    {3, "Sticker Coding", 25000, "/static/sticker.jpg", "Sticker lucu bertema coding untuk laptop atau meja kerja."},
}

var cart = map[int]int{} // productID -> qty

func main() {
	// Folder static untuk gambar & CSS
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", productList)
	http.HandleFunc("/add", addToCart)
	http.HandleFunc("/cart", viewCart)
	http.HandleFunc("/checkout", checkout)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func productList(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/products.html")
	tmpl.Execute(w, products)
}

func addToCart(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	cart[id]++
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func viewCart(w http.ResponseWriter, r *http.Request) {
	type CartItem struct {
		Product Product
		Qty     int
		Total   int
	}
	var items []CartItem
	var total int
	for id, qty := range cart {
		for _, p := range products {
			if p.ID == id {
				items = append(items, CartItem{p, qty, p.Price * qty})
				total += p.Price * qty
			}
		}
	}
	data := struct {
		Items []CartItem
		Total int
	}{
		Items: items,
		Total: total,
	}
	tmpl, _ := template.ParseFiles("templates/cart.html")
	tmpl.Execute(w, data)
}

func checkout(w http.ResponseWriter, r *http.Request) {
	cart = map[int]int{}
	fmt.Fprintln(w, "<h1>Checkout berhasil!</h1><a href='/'>Kembali belanja</a>")
}
