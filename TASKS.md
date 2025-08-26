# 🚪 Sentiric API Gateway Service - Görev Listesi (v4.0)

Bu belge, `api-gateway-service`'in geliştirme görevlerini projenin genel fazlarına uygun olarak listeler.

---

### **FAZ 1: Temel Proxy Yeteneği (Mevcut Durum)**

**Amaç:** Servisin, iç gRPC servislerini dış dünyaya temel bir REST arayüzü olarak sunabilmesini sağlamak.

-   [x] **Görev ID: GW-CORE-01 - gRPC-Gateway Kurulumu**
    -   **Durum:** ✅ **Tamamlandı**
    -   **Kabul Kriterleri:** `grpc-gateway/v2` kütüphanesi projeye entegre edilmiştir ve gelen HTTP isteklerini gRPC'ye yönlendirebilmektedir.

-   [x] **Görev ID: GW-CORE-02 - Arka Uç Servis Entegrasyonu**
    -   **Durum:** ✅ **Tamamlandı**
    -   **Kabul Kriterleri:** `user-service`'in gRPC endpoint'leri `/v1/users/...` altında REST olarak başarılı bir şekilde sunulmaktadır. Arka uç gRPC servislerine bağlantı, mTLS ile güvenli bir şekilde kurulmaktadır.

---

### **FAZ 2: Güvenlik ve Yetkilendirme Katmanı (Sıradaki Öncelik)**

**Amaç:** Servisi basit bir proxy'den, platformun tüm harici trafiğini koruyan bir güvenlik bekçisine dönüştürmek.

-   [ ] **Görev ID: GW-001 - JWT Kimlik Doğrulama Middleware'i**
    -   **Açıklama:** Gelen tüm isteklere `Authorization: Bearer <token>` başlığını kontrol eden, token'ı doğrulayan ve geçerli değilse `401 Unauthorized` hatası dönen bir middleware ekle.
    -   **Kabul Kriterleri:**
        -   [ ] `/health` gibi halka açık (public) endpoint'ler hariç tüm endpoint'ler bu middleware tarafından korunmalıdır.
        -   [ ] Geçersiz, süresi dolmuş veya imzasız bir token ile yapılan istek `401 Unauthorized` HTTP hatası ile reddedilmelidir.
        -   [ ] Geçerli bir token ile yapılan istek, middleware'den başarıyla geçmeli ve token içindeki kullanıcı bilgileri (`user_id`, `tenant_id`) isteğin `context`'ine eklenmelidir.

-   [ ] **Görev ID: GW-002 - Rol Tabanlı Yetkilendirme (RBAC)**
    -   **Açıklama:** JWT token'ın içindeki rollere (`claims`) bakarak, kullanıcının erişmeye çalıştığı endpoint için yetkisi olup olmadığını kontrol et.
    -   **Kabul Kriterleri:**
        -   [ ] "agent" rolüne sahip bir kullanıcının, sadece "admin" rolüne izin verilen bir endpoint'e (örn: `POST /v1/tenants`) erişmeye çalışması `403 Forbidden` hatası ile engellenmelidir.
        -   [ ] Yetki kuralları, kodun içinde kolayca yönetilebilir bir yapıda (örn: bir `map` veya `struct`) tanımlanmalıdır.

---

### **FAZ 3: Performans ve Dayanıklılık**

**Amaç:** Servisi yüksek trafikli üretim ortamları için optimize etmek.

-   [ ] **Görev ID: GW-003 - Hız Sınırlama (Rate Limiting)**
    -   **Açıklama:** Kötü niyetli kullanımı ve DoS saldırılarını önlemek için IP adresi veya API anahtarı bazlı hız sınırlama mekanizması ekle.
    -   **Kabul Kriterleri:**
        -   [ ] Belirlenen bir zaman aralığında (örn: 1 dakika) izin verilenden fazla istek yapan bir IP adresi, `429 Too Many Requests` HTTP hatası ile geçici olarak engellenmelidir.
        -   [ ] Hız limitleri `.env` dosyası üzerinden yapılandırılabilir olmalıdır.

-   [ ] **Görev ID: GW-004 - Merkezi Caching**
    -   **Açıklama:** Sık istenen ve nadiren değişen verileri Redis'te önbelleğe alarak arka uç servislerin yükünü azalt.
    -   **Durum:** ⬜ Planlandı.