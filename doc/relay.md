Kalau lampu/heater kamu disambungkan ke terminal **NO (Normally Open)**, maka sebenarnya tetap aman. Yang berubah hanya logika kontrolnya saja.

Pada relay ada 3 terminal utama:

- **COM** → jalur utama/input listrik
- **NO (Normally Open)** → terhubung ke COM saat relay aktif
- **NC (Normally Closed)** → terhubung ke COM saat relay mati

Ilustrasinya:

```text
Relay OFF:
COM --- NC
COM -X- NO

Relay ON:
COM -X- NC
COM --- NO
```

Kalau kamu pakai NO:

- relay mati → lampu mati
- relay aktif → lampu menyala

Jadi perilaku fisiknya tetap benar.

Masalah di kode kamu sekarang hanya:

- `lampState=true` ternyata membuat relay fisik mati,
- `lampState=false` ternyata membuat relay fisik hidup.

Karena relay module kamu active LOW.

Artinya:

| Kondisi software  | Pin GPIO | Relay fisik | Lampu di NO |
| ----------------- | -------- | ----------- | ----------- |
| `lampState=true`  | HIGH     | OFF         | Mati        |
| `lampState=false` | LOW      | ON          | Nyala       |

Jadi lampunya tetap bekerja, tetapi arti ON/OFF di program kebalik.

Itu yang berbahaya untuk debugging dan monitoring.

Kalau dibetulkan menjadi active LOW:

```cpp
digitalWrite(RELAY_PIN, on ? LOW : HIGH);
```

maka hasilnya jadi konsisten:

| Kondisi software  | Pin GPIO | Relay fisik | Lampu di NO |
| ----------------- | -------- | ----------- | ----------- |
| `lampState=true`  | LOW      | ON          | Nyala       |
| `lampState=false` | HIGH     | OFF         | Mati        |

Dan ini jauh lebih mudah dipahami ke depannya.
