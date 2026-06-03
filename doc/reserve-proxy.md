Benar — `BASE_URL` memang dipakai frontend untuk **request ke backend** (GET/POST/SSE). Yang perlu diluruskan adalah **endpoint backend yang di-hit browser berubah karena ada reverse proxy Nginx**.

Mari kita bedakan dengan jelas.

## 1. Jika frontend hit backend langsung (tanpa Nginx proxy)

Frontend:

```js
const BASE_URL = `http://${window.location.hostname}:8080/api`;
const SSE_URL = `http://${window.location.hostname}:8080/api/stream`;
```

Flow:

```text
Browser
    ├── FE → 5050 (nginx frontend)
    └── API → 8080 (backend langsung)
```

Browser benar-benar request ke:

```text
http://IP_SERVER:8080/api
```

Jadi port `8080` memang wajib ditulis.

---

## 2. Jika pakai Nginx reverse proxy (konfigurasi yang kamu tunjukkan)

Nginx kamu:

```nginx
location /api/ {
    proxy_pass http://backend-incubator:8080;
}
```

Artinya:

```text
Semua request /api dari browser
↓
diteruskan oleh nginx
↓
ke backend-incubator:8080
```

Flow berubah menjadi:

```text
Browser
    ↓
IP_SERVER:5050/api
    ↓
Nginx frontend
    ↓
backend-incubator:8080
```

Jadi frontend **tetap hit backend**, tetapi lewat nginx.

Browser request:

```text
http://IP_SERVER:5050/api
```

Nginx lalu proxy ke:

```text
http://backend-incubator:8080
```

Karena itu di frontend:

```js
const BASE_URL = `${window.location.origin}/api`;
const SSE_URL = `${window.location.origin}/api/stream`;
```

atau lebih sederhana:

```js
const BASE_URL = "/api";
const SSE_URL = "/api/stream";
```

masih merupakan **URL backend yang dipanggil frontend**, hanya saja backend-nya berada di balik proxy.

---

Jadi pilihannya:

### Opsi A — Frontend hit backend langsung

Tidak pakai proxy nginx:

```js
const BASE_URL = `http://${window.location.hostname}:8080/api`;
```

dan backend `8080` harus diexpose ke user.

---

### Opsi B — Frontend hit backend lewat nginx proxy

Pakai config nginx yang kamu punya:

```js
const BASE_URL = "/api";
```

dan user tidak perlu tahu `8080`.

Karena nginx kamu sudah ada `location /api`, setup yang konsisten adalah **Opsi B**.
