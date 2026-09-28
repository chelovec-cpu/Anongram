# Firebase Cloud Messaging

FCM credentials не хранятся в Git. Device tokens — PostgreSQL. Push jobs — Redis Stream `anongram:push`. Worker отправляет уведомления через Firebase Admin SDK.
