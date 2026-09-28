# Build verification

Проверено перед упаковкой:

- Go: все `.go` файлы прогнаны через `gofmt`; полноценный `go test` в этой среде не смог скачать внешние Go-модули из-за отсутствия сетевого доступа. GitHub Actions выполняет `go mod download`, `go test`, `go vet` и `go build`.
- Web: `package.json`, `tsconfig.json`, YAML/JSON/XML/PHP конфигурации проверены. TypeScript/JSX дополнительно проверен через `tsc` с временными типовыми заглушками; реальные npm-пакеты устанавливаются в GitHub Actions.
- Android: Gradle-проект, manifest и Kotlin-файл проверены; GitHub Actions устанавливает Gradle 8.9, Java 17 и собирает `assembleDebug`.
- Docker Compose: YAML разобран валидатором. Docker runtime в текущей среде отсутствует, поэтому контейнеры здесь не запускались.
- GitHub Actions публикует артефакты backend, web и Android APK.

## Что исправлено

- Исправлена WebSocket upgrade-реализация.
- Исправлена структура WebSocket-клиента и отправки событий.
- Защищены сообщения и подписка WebSocket проверкой членства в чате.
- Защищены API-маршруты от использования refresh-token как access-token.
- Исправлены JSON-теги регистрации/авторизации.
- Добавлено создание чатов через API.
- Добавлены проверки методов и membership.
- Исправлен TypeScript тип меню drawer.
- Исправлен Android Compose trailing content.
- Добавлен GitHub CI для Go/Web/Android.
- Добавлена выдача build artifacts.
- Docker Compose больше не требует локального `.env` для первого запуска: API использует `.env.example` как dev-конфигурацию.
- Добавлена автоматическая подготовка MinIO bucket.
- Улучшен Nginx proxy/WebSocket forwarding.
