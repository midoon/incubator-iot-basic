#include <WiFi.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>
#include <DHT.h>
#include <Wire.h>
#include <LiquidCrystal_I2C.h>

// ============================================================
// Konfigurasi — sesuaikan dengan environment kamu
// ============================================================

// WiFi
const char* WIFI_SSID     = "NAMA_WIFI_KAMU";
const char* WIFI_PASSWORD = "PASSWORD_WIFI_KAMU";

// MQTT Broker — isi dengan IP komputer yang menjalankan Mosquitto
// Contoh: "192.168.1.100"
const char* MQTT_BROKER   = "192.168.1.100";
const int   MQTT_PORT     = 1883;
const char* MQTT_CLIENT_ID = "esp32-incubator";

// Topic publish (ESP32 → Backend)
const char* TOPIC_TEMPERATURE = "incubator/temperature";
const char* TOPIC_HUMIDITY    = "incubator/humidity";
const char* TOPIC_MODE        = "incubator/mode";
const char* TOPIC_LAMP        = "incubator/lamp";

// Topic subscribe (Backend → ESP32)
const char* TOPIC_CMD_MODE    = "incubator/cmd/mode";
const char* TOPIC_CMD_LAMP    = "incubator/cmd/lamp";

// Pin
#define DHTPIN    4
#define DHTTYPE   DHT22
#define RELAY_PIN 5

// Interval publish sensor (ms)
const unsigned long PUBLISH_INTERVAL = 2000;

// Threshold suhu untuk mode auto (°C)
// Relay nyala jika suhu di bawah nilai ini
const float TEMP_THRESHOLD = 37.5;

// ============================================================
// Inisialisasi objek
// ============================================================

DHT dht(DHTPIN, DHTTYPE);
LiquidCrystal_I2C lcd(0x27, 16, 2);
WiFiClient wifiClient;
PubSubClient mqttClient(wifiClient);

// ============================================================
// State mesin — sumber kebenaran di ESP32
// ============================================================

String currentMode = "auto";   // "auto" | "manual"
bool   lampState   = false;    // status relay/lampu saat ini

unsigned long lastPublishTime = 0;

// ============================================================
// Fungsi kontrol relay
// ============================================================

void setRelay(bool on) {
  lampState = on;
  // Relay aktif HIGH — sesuaikan dengan modul relay kamu
  // Jika relay aktif LOW, ganti HIGH ↔ LOW di bawah
  digitalWrite(RELAY_PIN, on ? HIGH : LOW);
  Serial.print("[RELAY] ");
  Serial.println(on ? "ON" : "OFF");
}

// ============================================================
// Fungsi publish ke MQTT
// ============================================================

void publishSensorData(float temperature, float humidity) {
  StaticJsonDocument<64> tempDoc;
  tempDoc["temperature"] = round(temperature * 10.0) / 10.0; // 1 desimal
  char tempBuf[32];
  serializeJson(tempDoc, tempBuf);
  mqttClient.publish(TOPIC_TEMPERATURE, tempBuf);

  StaticJsonDocument<64> humDoc;
  humDoc["humidity"] = round(humidity);
  char humBuf[32];
  serializeJson(humDoc, humBuf);
  mqttClient.publish(TOPIC_HUMIDITY, humBuf);

  Serial.printf("[MQTT] Published → temp: %.1f°C, hum: %.0f%%\n", temperature, humidity);
}

void publishStatus() {
  // Publish mode saat ini
  StaticJsonDocument<32> modeDoc;
  modeDoc["mode"] = currentMode;
  char modeBuf[32];
  serializeJson(modeDoc, modeBuf);
  mqttClient.publish(TOPIC_MODE, modeBuf, true); // retain=true agar backend tahu state awal

  // Publish status lampu saat ini
  StaticJsonDocument<32> lampDoc;
  lampDoc["lamp"] = lampState;
  char lampBuf[32];
  serializeJson(lampDoc, lampBuf);
  mqttClient.publish(TOPIC_LAMP, lampBuf, true); // retain=true
}

// ============================================================
// Callback MQTT — dipanggil saat ada pesan masuk dari backend
// ============================================================

void onMQTTMessage(char* topic, byte* payload, unsigned int length) {
  // Konversi payload ke string
  char message[length + 1];
  memcpy(message, payload, length);
  message[length] = '\0';

  Serial.printf("[MQTT] Received ← topic: %s, payload: %s\n", topic, message);

  StaticJsonDocument<64> doc;
  DeserializationError err = deserializeJson(doc, message);
  if (err) {
    Serial.print("[MQTT] JSON parse error: ");
    Serial.println(err.c_str());
    return;
  }

  // ----- Command: ganti mode -----
  if (String(topic) == TOPIC_CMD_MODE) {
    if (!doc.containsKey("mode")) return;

    String newMode = doc["mode"].as<String>();

    if (newMode != "auto" && newMode != "manual") {
      Serial.println("[CMD] Invalid mode received, ignored.");
      return;
    }

    currentMode = newMode;
    Serial.printf("[CMD] Mode changed to: %s\n", currentMode.c_str());

    // Jika kembali ke auto: serahkan kontrol relay ke logika threshold
    // (akan ditangani di loop berikutnya)
    if (currentMode == "auto") {
      // Tidak langsung set relay di sini — biarkan loop() yang putuskan
      // berdasarkan suhu terkini
    }

    // Publish balik status mode yang sudah berubah
    publishStatus();
  }

  // ----- Command: kontrol lampu (hanya berlaku di mode manual) -----
  else if (String(topic) == TOPIC_CMD_LAMP) {
    if (!doc.containsKey("lamp")) return;

    if (currentMode != "manual") {
      Serial.println("[CMD] Lamp command ignored — not in manual mode.");
      return;
    }

    bool lampOn = doc["lamp"].as<bool>();
    setRelay(lampOn);

    // Publish balik status lampu yang sudah berubah
    publishStatus();
  }
}

// ============================================================
// Koneksi WiFi
// ============================================================

void connectWiFi() {
  Serial.printf("[WiFi] Connecting to %s", WIFI_SSID);
  WiFi.mode(WIFI_STA);
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  int attempts = 0;
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
    attempts++;
    if (attempts > 40) {
      // Setelah 20 detik gagal, restart ESP32
      Serial.println("\n[WiFi] Failed to connect, restarting...");
      ESP.restart();
    }
  }

  Serial.println();
  Serial.print("[WiFi] Connected! IP: ");
  Serial.println(WiFi.localIP());
}

// ============================================================
// Koneksi & reconnect MQTT
// ============================================================

void connectMQTT() {
  while (!mqttClient.connected()) {
    Serial.printf("[MQTT] Connecting to broker %s:%d ...\n", MQTT_BROKER, MQTT_PORT);

    if (mqttClient.connect(MQTT_CLIENT_ID)) {
      Serial.println("[MQTT] Connected!");

      // Subscribe ke topic command dari backend
      mqttClient.subscribe(TOPIC_CMD_MODE);
      mqttClient.subscribe(TOPIC_CMD_LAMP);
      Serial.println("[MQTT] Subscribed to command topics.");

      // Publish status awal agar backend tahu state ESP32 saat baru connect
      publishStatus();

    } else {
      Serial.printf("[MQTT] Failed, rc=%d. Retry in 5s...\n", mqttClient.state());
      delay(5000);
    }
  }
}

// ============================================================
// Update LCD
// ============================================================

void updateLCD(float t, float h) {
  lcd.clear();

  // Baris 1: Suhu
  lcd.setCursor(0, 0);
  lcd.print("T:");
  lcd.print(t, 1);
  lcd.print((char)223); // karakter derajat °
  lcd.print("C ");
  // Tampilkan mode singkat di pojok kanan
  lcd.setCursor(11, 0);
  lcd.print(currentMode == "auto" ? "[AUTO]" : "[MAN] ");

  // Baris 2: Kelembapan & status lampu
  lcd.setCursor(0, 1);
  lcd.print("H:");
  lcd.print(h, 0);
  lcd.print("% ");
  lcd.setCursor(9, 1);
  lcd.print("L:");
  lcd.print(lampState ? "ON " : "OFF");
}

// ============================================================
// Setup
// ============================================================

void setup() {
  Serial.begin(115200);
  Serial.println("\n=== Incubator ESP32 Starting ===");

  // Inisialisasi pin relay
  pinMode(RELAY_PIN, OUTPUT);
  setRelay(false); // Pastikan relay mati saat startup

  // Inisialisasi LCD
  lcd.init();
  lcd.backlight();
  lcd.setCursor(0, 0);
  lcd.print("  Incubator v1  ");
  lcd.setCursor(0, 1);
  lcd.print(" Menghubungkan..");
  delay(1500);

  // Inisialisasi sensor DHT
  dht.begin();

  // Koneksi WiFi
  connectWiFi();

  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("WiFi: OK");
  lcd.setCursor(0, 1);
  lcd.print(WiFi.localIP());
  delay(1500);

  // Setup MQTT
  mqttClient.setServer(MQTT_BROKER, MQTT_PORT);
  mqttClient.setCallback(onMQTTMessage);

  // Koneksi MQTT
  connectMQTT();

  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("MQTT: OK");
  delay(1000);

  Serial.println("[MAIN] Setup complete. Starting loop...");
}

// ============================================================
// Loop utama
// ============================================================

void loop() {
  // Jaga koneksi WiFi
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("[WiFi] Disconnected, reconnecting...");
    connectWiFi();
  }

  // Jaga koneksi MQTT & proses pesan masuk
  if (!mqttClient.connected()) {
    connectMQTT();
  }
  mqttClient.loop(); // WAJIB dipanggil agar callback berjalan

  // Publish sensor setiap PUBLISH_INTERVAL ms
  unsigned long now = millis();
  if (now - lastPublishTime >= PUBLISH_INTERVAL) {
    lastPublishTime = now;

    float h = dht.readHumidity();
    float t = dht.readTemperature();

    if (isnan(h) || isnan(t)) {
      Serial.println("[SENSOR] Failed to read DHT22!");
      lcd.clear();
      lcd.setCursor(0, 0);
      lcd.print("Sensor Error!");
      return;
    }

    // Publish data sensor ke backend
    publishSensorData(t, h);

    // ---- Logika kontrol relay ----
    if (currentMode == "auto") {
      // Mode auto: relay dikontrol otomatis berdasarkan threshold suhu
      bool shouldHeat = (t < TEMP_THRESHOLD);
      if (shouldHeat != lampState) {
        setRelay(shouldHeat);
        publishStatus(); // Beritahu backend bahwa status lampu berubah
      }
    }
    // Mode manual: relay dikontrol penuh oleh command dari backend
    // (sudah ditangani di onMQTTMessage, tidak ada aksi di sini)

    // Update LCD
    updateLCD(t, h);

    Serial.printf("[SENSOR] Temp: %.1f°C | Hum: %.0f%% | Mode: %s | Lamp: %s\n",
      t, h, currentMode.c_str(), lampState ? "ON" : "OFF");
  }
}
