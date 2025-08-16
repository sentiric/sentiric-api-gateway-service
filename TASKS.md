# 🚪 Sentiric API Gateway Service - Görev Listesi

Bu belge, `api-gateway-service`'in geliştirme yol haritasını ve önceliklerini tanımlar.

---

### Faz 1: Temel Proxy Yeteneği (Mevcut Durum)

Bu faz, servisin temel gRPC-to-REST proxy görevini yerine getirmesini hedefler.

-   [x] **gRPC-Gateway Kurulumu:** `grpc-gateway/v2` kütüphanesini kullanarak temel bir sunucu oluşturma.
-   [x] **User Service Entegrasyonu:** `user-service`'in gRPC endpoint'lerini `/v1/users` altında REST olarak sunma.
-   [x] **mTLS İstemci:** Arka uç gRPC servislerine güvenli (mTLS) bağlantı kurabilme.

---

### Faz 2: Güvenlik ve Yetkilendirme Katmanı (Sıradaki Öncelik)

Bu faz, servisi basit bir proxy'den, platformun güvenlik bekçisine dönüştürmeyi hedefler.

-   [ ] **Görev ID: GW-001 - JWT Kimlik Doğrulama Middleware'i**
    -   **Açıklama:** Gelen tüm isteklere `Authorization: Bearer <token>` başlığını kontrol eden, token'ı doğrulayan ve geçerli değilse `401 Unauthorized` hatası dönen bir middleware ekle.
    -   **Durum:** ⬜ Planlandı.

-   [ ] **Görev ID: GW-002 - Rol Tabanlı Yetkilendirme (RBAC)**
    -   **Açıklama:** JWT token'ın içindeki rollere (`claims`) bakarak, kullanıcının erişmeye çalıştığı endpoint için yetkisi olup olmadığını kontrol et. Yetkisi yoksa `403 Forbidden` hatası dön.
    -   **Durum:** ⬜ Planlandı.

---

### Faz 3: Performans ve Dayanıklılık

Bu faz, servisi yüksek trafikli üretim ortamları için optimize etmeyi hedefler.

-   [ ] **Görev ID: GW-003 - Rate Limiting**
    -   **Açıklama:** Kötü niyetli kullanımı ve DDoS saldırılarını önlemek için IP adresi veya API anahtarı bazlı hız sınırlama mekanizması ekle.
    -   **Durum:** ⬜ Planlandı.

-   [ ] **Görev ID: GW-004 - Merkezi Caching**
    -   **Açıklama:** Sık istenen ve nadiren değişen verileri (örn: `/dialplans` listesi) Redis'te önbelleğe alarak arka uç servislerin yükünü azalt.
    -   **Durum:** ⬜ Planlandı.