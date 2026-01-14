package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

type URLData struct {
	OriginalURL string `json:"original_url"`
	Clicks      int    `json:"clicks"`
}

type URLManager struct {
	mu    sync.Mutex
	Links map[string]URLData `json:"links"`
	File  string             `json:"-"`
}

func NewURLManager(filename string) *URLManager {
	mgr := &URLManager{
		Links: make(map[string]URLData),
		File:  filename,
	}
	file, err := os.ReadFile(filename)
	if err == nil {
		json.Unmarshal(file, &mgr.Links)
	}
	return mgr
}

func (m *URLManager) save() {
	data, _ := json.MarshalIndent(m.Links, "", "  ")
	os.WriteFile(m.File, data, 0644)
}

func generateKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seed.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	db := NewURLManager("links.json")

	// 1. Asosiy sahifa (Frontend)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, `
			<html>
			<head>
				<title>Go Shortener</title>
				<style>
					body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: #f0f2f5; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
					.container { background: white; padding: 40px; border-radius: 15px; box-shadow: 0 10px 25px rgba(0,0,0,0.1); width: 450px; text-align: center; }
					h2 { color: #1c1e21; margin-bottom: 25px; }
					input { width: 100%%; padding: 12px; margin-bottom: 15px; border: 1px solid #ddd; border-radius: 8px; box-sizing: border-box; font-size: 16px; }
					button { width: 100%%; padding: 12px; background: #0084ff; color: white; border: none; border-radius: 8px; font-size: 16px; cursor: pointer; transition: background 0.3s; }
					button:hover { background: #0073e6; }
					#result { margin-top: 20px; padding: 15px; border-radius: 8px; background: #e7f3ff; color: #0052cc; font-weight: bold; display: none; word-break: break-all; }
					.nav { margin-top: 20px; font-size: 14px; }
					.nav a { color: #606770; text-decoration: none; }
				</style>
			</head>
			<body>
				<div class="container">
					<h2>🚀 URL Qisqartirgich</h2>
					<input id="url" type="text" placeholder="https://uzun-linkni-shu-yerga-qo'ying.com">
					<button onclick="shorten()">Qisqartirish</button>
					<div id="result"></div>
					<div class="nav"><a href="/list">📊 Statistikani ko'rish</a></div>
				</div>
				<script>
					async function shorten() {
						let u = document.getElementById('url').value;
						if(!u) return;
						let res = await fetch('/shorten?url=' + encodeURIComponent(u));
						let text = await res.text();
						let resDiv = document.getElementById('result');
						resDiv.innerText = text;
						resDiv.style.display = 'block';
					}
				</script>
			</body>
			</html>
		`)
	})

	// 2. URL qisqartirish API
	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		targetURL := r.URL.Query().Get("url")
		if targetURL == "" {
			http.Error(w, "URL kiritilmadi", 400)
			return
		}
		db.mu.Lock()
		defer db.mu.Unlock()
		key := ""
		for k, v := range db.Links {
			if v.OriginalURL == targetURL {
				key = k
				break
			}
		}
		if key == "" {
			key = generateKey()
			db.Links[key] = URLData{OriginalURL: targetURL, Clicks: 0}
			db.save()
		}
		fmt.Fprintf(w, "http://localhost:8080/r/%s", key)
	})

	// 3. Yo'naltirish (Redirect)
	http.HandleFunc("/r/", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[len("/r/"):]
		db.mu.Lock()
		data, ok := db.Links[key]
		if ok {
			data.Clicks++
			db.Links[key] = data
			db.save()
		}
		db.mu.Unlock()
		if !ok {
			http.Error(w, "Havola topilmadi", 404)
			return
		}
		http.Redirect(w, r, data.OriginalURL, http.StatusFound)
	})

	// 4. Chiroyli Statistika va O'chirish funksiyasi
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		db.mu.Lock()
		defer db.mu.Unlock()

		// O'chirish so'rovi kelsa
		if delKey := r.URL.Query().Get("delete"); delKey != "" {
			delete(db.Links, delKey)
			db.save()
			http.Redirect(w, r, "/list", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<html>
			<head>
				<title>Statistika</title>
				<style>
					body { font-family: sans-serif; background: #f4f7f6; padding: 50px; }
					table { width: 100%%; border-collapse: collapse; background: white; box-shadow: 0 5px 15px rgba(0,0,0,0.05); }
					th, td { padding: 15px; text-align: left; border-bottom: 1px solid #eee; }
					th { background: #0084ff; color: white; }
					.btn-del { color: #ff4d4d; text-decoration: none; font-weight: bold; border: 1px solid #ff4d4d; padding: 5px 10px; border-radius: 5px; }
					.btn-del:hover { background: #ff4d4d; color: white; }
					.back { margin-bottom: 20px; display: inline-block; text-decoration: none; color: #0084ff; }
				</style>
			</head>
			<body>
				<a href="/" class="back">← Bosh sahifa</a>
				<h2>📊 Havolalar statistikasi</h2>
				<table>
					<tr><th>Qisqa kod</th><th>Asl havola</th><th>Clicks</th><th>Amal</th></tr>`)
		for k, v := range db.Links {
			fmt.Fprintf(w, `
				<tr>
					<td><b>%s</b></td>
					<td style="color: #666;">%s</td>
					<td><b style="color: #28a745;">%d</b></td>
					<td><a href="/list?delete=%s" class="btn-del" onclick="return confirm('Ochirilsinmi?')">O'chirish</a></td>
				</tr>`, k, v.OriginalURL, v.Clicks, k)
		}
		fmt.Fprintf(w, `</table></body></html>`)
	})

	fmt.Println("Dastur ishga tushdi: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
