# AnonGram Full-Stack Starter

Стартовый monorepo: Go API + PostgreSQL + Redis + WebSocket + JWT/refresh + S3/MinIO + React/TypeScript web + Kotlin/Compose Android + WebRTC signaling + FFmpeg pipeline + Nginx + Docker + Prometheus/Grafana/Loki + GitHub Actions.

## Быстрый запуск

```bash
cp .env.example .env
docker compose up --build
```

Веб: http://localhost:8080  
API: http://localhost:8080/api/health  
MinIO: http://localhost:9001  
Prometheus: http://localhost:9090  
Grafana: http://localhost:3000

Это рабочий starter, а не готовая публичная соцсеть. Для production нужны реальные секреты, HTTPS/WSS, FCM credentials, TURN/STUN, backups PostgreSQL, Cloudflare DNS/CDN, rate limiting и настройка S3.

## CI

GitHub Actions проверяет три части отдельно: Go backend (`go test` + `go vet`), web (`npm run build`) и Android (`assembleDebug`). Для web используется `npm install`, поэтому первый запуск автоматически создаёт lockfile в рабочем CI-кэше.

## Что реально работает

Это стартовый проект, который можно собрать. Production-сервисы (FCM, S3 uploads, TURN, Cloudflare, TLS certificates) требуют секретов и инфраструктурных настроек; они не притворяются полностью настроенными внутри репозитория.
