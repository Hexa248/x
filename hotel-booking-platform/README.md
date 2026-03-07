# Hotel Booking Platform (Microservices Skeleton)

Arsitektur ini memecah sistem ke API Gateway + service terpisah agar alur booking lebih mendekati platform besar.

## Fitur yang dicakup
- Search/filter/sorting hotel (harga, rating, bintang, fasilitas)
- Booking + promo/voucher + service tambahan
- Cancellation + refund
- Staff check-in/check-out
- Review dan rating hotel
- Wishlist hotel
- Notifikasi event booking/status
- Admin analytics (booking, revenue, occupancy, top hotel)

## Menjalankan cepat (contoh)
- API Gateway: `go run ./api-gateway`
- Hotel service: `go run ./services/hotel-service`
- Booking service: `go run ./services/booking-service`
- Review service: `go run ./services/review-service`
- Promo service: `go run ./services/promo-service`
- Analytics service: `go run ./services/analytics-service`
