#include "DHT.h"
#include <WiFi.h>
#include <HTTPClient.h>

#define DHTPIN 4
#define DHTTYPE DHT22
DHT dht(DHTPIN, DHTTYPE);

#define SOUND_PIN 34
#define LED_PIN 25
#define BUZZER_PIN 26

const char* ssid = "YOUR_WIFI";
const char* password = "YOUR_PASSWORD";
const char* serverURL = "http://YOUR_SERVER_IP:8080/data";

void setup() {
  Serial.begin(115200);
  pinMode(LED_PIN, OUTPUT);
  pinMode(BUZZER_PIN, OUTPUT);
  dht.begin();

  WiFi.begin(ssid, password);
  Serial.print("Connecting to WiFi");
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("\nConnected to WiFi");
}

void loop() {
  float temp = dht.readTemperature();
  float hum = dht.readHumidity();
  int soundLevel = analogRead(SOUND_PIN);

  if (!isnan(temp) && !isnan(hum)) {
    HTTPClient http;
    http.begin(serverURL);
    http.addHeader("Content-Type", "application/json");

    String json = "{\"temperature\":" + String(temp) +
                  ",\"humidity\":" + String(hum) +
                  ",\"sound\":" + String(soundLevel) + "}";

    int httpCode = http.POST(json);

    if (httpCode == 200) {
      String payload = http.getString();
      if (payload.indexOf("ALERT") >= 0) {
        digitalWrite(LED_PIN, HIGH);
        digitalWrite(BUZZER_PIN, HIGH);
      } else {
        digitalWrite(LED_PIN, LOW);
        digitalWrite(BUZZER_PIN, LOW);
      }
    }

    http.end();
  }

  delay(1000); 
}
