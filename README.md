# 🚪 Sentiric API Gateway Service

[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![Language](https://img.shields.io/badge/language-Go-blue.svg)]()
[![Protocol](https://img.shields.io/badge/protocol-gRPC--Gateway_(REST)-orange.svg)]()

**Sentiric API Gateway Service**, tüm harici istemciler (`dashboard-ui`, `cli`, üçüncü parti uygulamalar) için platformun **tek, birleşik ve güvenli giriş kapısıdır.** Temel görevi, dış dünyanın konuştuğu RESTful JSON API'lerini, platformun iç dünyasının konuştuğu yüksek performanslı gRPC protokolüne çevirmektir.

Bu servis, `grpc-gateway` kütüphanesini kullanarak Protobuf tanımlarından otomatik olarak bir REST proxy'si oluşturur.

## 🎯 Temel Sorumluluklar

*   **Protokol Çevirimi:** Gelen HTTP/REST isteklerini, ilgili gRPC mikroservisine (örn: `user-service`) yönlendirir ve gRPC yanıtını HTTP/JSON formatına çevirerek istemciye geri döner.
*   **Merkezi Giriş Noktası:** Dış istemcilerin, platformdaki onlarca mikroservisin adresini tek tek bilmesine gerek kalmadan tek bir adresten tüm platforma erişmesini sağlar.
*   **Güvenlik ve Yetkilendirme (Gelecek):** Gelen tüm istekler için merkezi kimlik doğrulama (JWT) ve yetkilendirme (RBAC) katmanı olarak görev yapacaktır.
*   **Rate Limiting ve Caching (Gelecek):** Platformu kötü niyetli kullanımdan korumak ve performansı artırmak için merkezi hız sınırlama ve önbellekleme mekanizmaları uygulayacaktır.

## 🛠️ Teknoloji Yığını

*   **Dil:** Go
*   **Protokol Çevirimi:** `grpc-gateway/v2`
*   **Servisler Arası İletişim:** gRPC (mTLS ile)
*   **Loglama:** `zerolog` ile yapılandırılmış, ortama duyarlı loglama.

## 🔌 API Etkileşimleri

*   **Gelen (Sunucu):**
    *   `sentiric-dashboard-ui` (REST/JSON)
    *   `sentiric-cli` (REST/JSON)
*   **Giden (İstemci):**
    *   `sentiric-user-service` (gRPC)
    *   `sentiric-dialplan-service` (gRPC)
    *   (Gelecekte eklenecek diğer tüm gRPC servisleri...)

## 🚀 Yerel Geliştirme

1.  **Bağımlılıkları Yükleyin:**
2.  **Ortam Değişkenlerini Ayarlayın:** `.env.example` dosyasını `.env` olarak kopyalayın ve gerekli değişkenleri doldurun.
3.  **Servisi Çalıştırın:**

## 🤝 Katkıda Bulunma

Katkılarınızı bekliyoruz! Lütfen projenin ana [Sentiric Governance](https://github.com/sentiric/sentiric-governance) reposundaki kodlama standartlarına ve katkıda bulunma rehberine göz atın.

---
## 🏛️ Anayasal Konum

Bu servis, [Sentiric Anayasası'nın (v11.0)](https://github.com/sentiric/sentiric-governance/blob/main/docs/blueprint/Architecture-Overview.md) **Zeka & Orkestrasyon Katmanı**'nda yer alan merkezi bir bileşendir.