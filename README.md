# 🚀 Go URL Shortener

Ushbu loyiha Go (Golang) tilida yozilgan zamonaviy va tezkor URL qisqartirgich servisidir. Loyiha orqali uzun havolalarni qisqa va qulay ko'rinishga keltirish, ularning statistikasini kuzatish va boshqarish mumkin.

## ✨ Xususiyatlari

* **URL Qisqartirish:** Uzun havolalarni 6 xonali noyob kalitlar orqali qisqartirish.
* **Redirect (Yo'naltirish):** Qisqa havola orqali darhol asl manzilga o'tish.
* **Click Counter:** Har bir havola necha marta bosilganini real-vaqtda hisoblash.
* **Persistence (Saqlash):** Ma'lumotlar `links.json` faylida saqlanadi, server o'chsa ham ma'lumotlar saqlanib qoladi.
* **Dashboard:** Barcha havolalarni ko'rish va keraksizlarini o'chirish imkoniyati.
* **Concurrency Safe:** `sync.Mutex` yordamida bir vaqtning o'zida ko'plab so'rovlar bilan xavfsiz ishlash.

## 🛠 Texnologiyalar

* **Backend:** Go (Golang)
* **Frontend:** HTML5, CSS3, JavaScript (Fetch API)
* **Storage:** JSON-based database

## 📁 Fayllar strukturasi

Sizning loyihangiz quyidagi tartibda tuzilgan:
```text
GO-URL-SHORTENER/
├── static/
│   └── index.html  # Asosiy UI sahifasi
├── engine.go       # Yordamchi funksiyalar (Kalit yaratish)
├── main.go         # Server logikasi va marshrutlar
├── go.mod          # Go modul fayli
├── links.json      # Ma'lumotlar bazasi
└── .gitignore      # Git uchun keraksiz fayllar ro'yxati