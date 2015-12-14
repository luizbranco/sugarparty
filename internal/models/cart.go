package models

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Cart struct {
	items    map[string]int
	Items    []CartItem
	Quantity int
	Price    float64
	Ready    bool
}

type CartItem struct {
	Product
	Quantity int
	Price    float64
}

func NewCart(r *http.Request) *Cart {
	c := &Cart{items: make(map[string]int)}
	cookie, err := r.Cookie("cart")
	if err != nil {
		return c
	}
	items := strings.Split(cookie.Value, " ")
	if len(items)%2 != 0 {
		return c
	}

	var keys []interface{}

	for i := 0; i < len(items); i += 2 {
		n, err := strconv.Atoi(items[i+1])
		if err == nil {
			k := items[i]
			keys = append(keys, k)
			c.Add(k, n)
		}
	}

	products, _ := FindProducts(keys)
	c.Items = make([]CartItem, 0, len(products))
	for _, p := range products {
		qty := c.items[p.ID]
		price := p.Price * float64(qty)
		i := CartItem{p, qty, price}
		c.Quantity += qty
		c.Price += price
		c.Items = append(c.Items, i)
	}

	return c
}

func (c *Cart) Add(id string, qty int) {
	n := c.items[id] + qty
	if n >= 1 {
		c.items[id] = n
	} else {
		c.Remove(id)
	}
}

func (c *Cart) Remove(id string) {
	delete(c.items, id)
}

func (c *Cart) Save(w http.ResponseWriter) {
	var val []string
	for k, v := range c.items {
		val = append(val, k, strconv.Itoa(v))
	}
	cookie := &http.Cookie{
		Path:     "/",
		Name:     "cart",
		Value:    strings.Join(val, " "),
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func (c *Cart) Destroy(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Path:     "/",
		Name:     "cart",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func (c *Cart) Total() error {
	return nil
}
