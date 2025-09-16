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

-   **Görev ID: GW-005 - "Sessiz" Sağlık Kontrolü Endpoint'i Ekle**
    -   **Durum:** ✅ **Tamamlandı**
    -   **Açıklama:** Log kirliliğini önlemek ve `docker-compose`'da `service_healthy` koşulunu desteklemek için, loglama middleware'inden geçmeyen bir `/healthz` endpoint'i eklendi.
    -   **Kabul Kriterleri:**
        -   [x] `/healthz` endpoint'i `200 OK` yanıtı dönmeli.
        -   [x] Bu endpoint'e yapılan istekler, `loggingMiddleware` tarafından loglanmamalıdır.  

---


### **Mevcut Durum ve Rolü: "Zırhlı Ön Kapı (SBC)"**

`sip-gateway-service`, basit bir proxy'den çok daha fazlası olarak tasarlanmış. `README.md` dosyasında da belirtildiği gibi, bir **Session Border Controller (SBC)**'nin temel görevlerini yerine getiriyor. Bu, platformunuzu kurumsal düzeyde bir telekomünikasyon çözümüne dönüştüren en önemli adımlardan biridir.

Şu anda **başarıyla tamamladığı** görevler şunlardır:

1.  **Güvenlik Kalkanı:** Dış dünyadan gelen ham SIP trafiğini ilk karşılayan servis budur. `sip-signaling-service` gibi daha karmaşık mantık içeren servisleri, internetin potansiyel tehlikelerinden (hatalı formatlanmış paketler, tarama girişimleri vb.) tamamen izole eder.

2.  **Ağ Tercümanı (NAT Traversal - EN ÖNEMLİ GÖREVİ):**
    *   **Sorun:** Ev veya ofis ağlarının arkasındaki SIP istemcileri (MicroSIP gibi), `INVITE` paketlerinde kendi yerel IP adreslerini (`192.168.1.3` gibi) gönderirler. Platform, bu adrese doğrudan ses paketi (RTP) gönderemez.
    *   **Çözüm:** `sip-gateway`, `[18:57:08]` logunda görüldüğü gibi, paketin **gerçekte geldiği** genel IP adresini (`100.104.184.19:61097`) görür. `src/main.rs` dosyasındaki `Transactions` `HashMap`'i sayesinde, bu gerçek adresi çağrının `Call-ID`'si ile eşleştirerek hafızasında tutar.
    *   `sip-signaling`'den bir yanıt geldiğinde (`[18:57:08] handle_response_from_signaling`), `gateway` bu hafızayı kullanarak yanıtı doğru istemcinin gerçek genel IP adresine geri gönderir. Bu olmadan, sesli iletişim **asla kurulamazdı.**

3.  **Akıllı Yönlendirici:** Gelen tüm geçerli SIP trafiğini tek bir hedefe, yani `sip-signaling:5060` adresine yönlendirerek iç ağdaki mimariyi basitleştirir.

**Teknik Değerlendirme:** Rust ve Tokio ile yazılmış bu servis, `stateful` (durum bilgisi tutan) bir proxy olarak son derece verimli ve doğru çalışıyor. `Call-ID` ve `CSeq` metodu bazlı işlem takibi, standartlara uygun ve sağlam bir yaklaşımdır.

---

### **Gelecek Planları ve Yol Haritası**

`TASKS.md` dosyanız, bu sağlam temel üzerine inşa edilecek geleceği net bir şekilde ortaya koyuyor:

#### **Faz 2: Güvenlik ve Dayanıklılık (Sıradaki Adımlar)**
*   **Görev ID: `GW-SIP-001` - Hız Sınırlama (Rate Limiting):** Bu, platformunuzu basit DoS (Denial of Service) saldırılarına karşı koruyacak olan bir sonraki mantıksal adımdır. Belirli bir IP'den saniyede gelen istek sayısını sınırlayarak, `sip-signaling` servisinin aşırı yüklenmesini engelleyecektir.
*   **Görev ID: `GW-SIP-002` - IP Beyaz/Kara Liste:** Güvenliği bir adım öteye taşıyarak, sadece bilinen ve güvenilen telekom operatörlerinden veya müşterilerden gelen trafiğe izin vermenizi sağlar.

#### **Faz 3: Gelişmiş Protokol Desteği (Stratejik Hedef)**
*   **Görev ID: `GW-SIP-003` - WebRTC Entegrasyonu:** Bu, platformunuzun en büyük stratejik hedeflerinden biridir. `web-agent-ui` gibi tarayıcı tabanlı uygulamaların, mikrofonlarını kullanarak doğrudan platformla sesli iletişim kurabilmesini sağlar. Bu görev, `sip-gateway`'in SIP over WebSocket (WSS) trafiğini alıp, iç ağdaki standart SIP/UDP'ye çevirmesini gerektirecektir.

---

### **Özet ve Stratejik Önem**

Evet, `sip-gateway` için sadece bir planımız yok, aynı zamanda **halihazırda çalışan, kritik bir bileşenimiz var.** Bu servis olmasaydı, platformunuz:
*   **Kırılgan** olurdu (her türlü hatalı SIP paketi iç mantığı çökertebilirdi).
*   **Güvensiz** olurdu (iç ağ yapısı dış dünyaya açık olurdu).
*   **Fonksiyonel olmazdı** (NAT arkasındaki kullanıcılarla sesli iletişim kuramazdı).

Loglarınız, bu servisin görevini kusursuz bir şekilde yerine getirdiğini kanıtlıyor. `TASKS.md`'deki gelecek planları ise onu basit bir proxy'den, tam teşekküllü, kurumsal düzeyde bir **iletişim güvenlik duvarına** dönüştürme vizyonunu gösteriyor.