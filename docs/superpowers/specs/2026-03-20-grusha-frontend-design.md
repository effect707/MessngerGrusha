# Grusha Messenger — Frontend Design Specification

## Overview

Telegram-like web-клиент для мессенджера "Груша". Демонстрирует полный flow бэкенда: регистрация, авторизация, чаты, каналы, файлы, реакции, real-time сообщения. React + Vite, подключается к Go-бэкенду через grpc-gateway (REST) и WebSocket.

## Architecture

### System Diagram

```
┌─────────────────────┐         ┌──────────────────────────────┐
│   React App (Vite)  │         │       Go Binary               │
│   localhost:5173     │         │   localhost:8080               │
│                     │  HTTP   │                                │
│  fetch('/api/...')  ├────────►│  grpc-gateway  ──► gRPC services│
│                     │  JSON   │  /api/auth/*                   │
│  WebSocket          │         │  /api/chats/*                  │
│  ws://host/ws       ├────────►│  /api/messages/*               │
│                     │         │  /api/channels/*               │
│                     │         │  /api/users/*                  │
│                     │         │  /api/files/*  (HTTP handler)  │
│                     │         │  /ws           (WebSocket)     │
└─────────────────────┘         └──────────────────────────────┘
        │                                    │
        │ Vite proxy                         │
        │ /api/* → :8080                     ▼
        │ /ws   → :8080              Postgres / Redis / MinIO
```

React app проксирует все запросы через Vite dev server. В production — nginx или Go может раздавать static.

### Frontend Tech Stack

| Concern | Choice | Rationale |
|---------|--------|-----------|
| Framework | React 19 + TypeScript | Компонентный подход, широко используется |
| Build | Vite | Быстрый dev server, proxy для API |
| State | Zustand | Минимальный бойлерплейт, отдельные stores |
| Routing | React Router v7 | Стандарт для SPA |
| Styling | CSS Modules | Scoped стили, без лишних зависимостей |
| HTTP | fetch (native) | Простой wrapper с JWT header |
| WebSocket | native WebSocket | Один коннект, auto-reconnect |

Без тяжёлых UI-библиотек (MUI, Tailwind) — чистый CSS в стиле Telegram dark theme.

## Backend Changes: grpc-gateway

### Proto-файлы — HTTP аннотации

Добавить `google.api.http` option к каждому RPC-методу. Примеры:

```protobuf
import "google/api/annotations.proto";

service AuthService {
  rpc Register(RegisterRequest) returns (RegisterResponse) {
    option (google.api.http) = {
      post: "/api/auth/register"
      body: "*"
    };
  }
  rpc Login(LoginRequest) returns (LoginResponse) {
    option (google.api.http) = {
      post: "/api/auth/login"
      body: "*"
    };
  }
  // ...
}
```

### buf.gen.yaml — добавить плагины

```yaml
version: v2
plugins:
  - local: protoc-gen-go
    out: api/gen
    opt:
      - paths=source_relative
  - local: protoc-gen-go-grpc
    out: api/gen
    opt:
      - paths=source_relative
  - local: protoc-gen-grpc-gateway
    out: api/gen
    opt:
      - paths=source_relative
```

### buf.yaml — добавить googleapis dependency

```yaml
deps:
  - buf.build/googleapis/googleapis
```

### app.go — регистрация gateway mux

```go
grpcAddr := fmt.Sprintf("localhost:%d", cfg.GRPC.Port) // e.g. localhost:50051
gwMux := runtime.NewServeMux()
opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, opts)
chatpb.RegisterChatServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, opts)
messagepb.RegisterMessageServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, opts)
channelpb.RegisterChannelServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, opts)
userpb.RegisterUserServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, opts)

// File handlers registered BEFORE gateway wildcard (more specific routes first)
httpMux.HandleFunc("/api/files/upload", fileHandler.Upload)
httpMux.HandleFunc("/api/files/download", fileHandler.Download)
httpMux.Handle("/ws", wsHandler)       // existing WebSocket handler
httpMux.Handle("/api/", gwMux)         // gateway catches remaining /api/* routes
```

**Примечания:**
- grpc-gateway подключается к gRPC серверу на `cfg.GRPC.Port` (50051). HTTP сервер на `cfg.HTTP.Port` (8080).
- grpc-gateway автоматически прокидывает `Authorization` header как gRPC metadata — существующий auth interceptor работает без изменений.
- Auth interceptor должен пропускать без токена: `Register`, `Login`, `RefreshTokens` (public endpoints).
- Go 1.22+ `ServeMux` выбирает наиболее специфичный маршрут, но для ясности файловые хендлеры регистрируются первыми.

### Proto HTTP аннотации — особые случаи

**AddMember** — `chat_id` из URL path, `user_id` из JSON body:
```protobuf
rpc AddMember(AddMemberRequest) returns (AddMemberResponse) {
  option (google.api.http) = {
    post: "/api/chats/{chat_id}/members"
    body: "*"
  };
}
```

**Notification RPCs** живут в `UserService`, но имеют отдельный URL prefix `/api/notifications/`:
```protobuf
// Inside UserService proto:
rpc GetNotifications(GetNotificationsRequest) returns (GetNotificationsResponse) {
  option (google.api.http) = { get: "/api/notifications" };
}
rpc MarkNotificationRead(MarkNotificationReadRequest) returns (MarkNotificationReadResponse) {
  option (google.api.http) = { post: "/api/notifications/{notification_id}/read" };
}
rpc MarkAllNotificationsRead(MarkAllNotificationsReadRequest) returns (MarkAllNotificationsReadResponse) {
  option (google.api.http) = { post: "/api/notifications/read-all" };
}
rpc GetUnreadCount(GetUnreadCountRequest) returns (GetUnreadCountResponse) {
  option (google.api.http) = { get: "/api/notifications/unread-count" };
}
```

### REST API Endpoints

```
AUTH
  POST /api/auth/register             → Register
  POST /api/auth/login                → Login
  POST /api/auth/logout               → Logout
  POST /api/auth/logout-all           → LogoutAll
  POST /api/auth/refresh              → RefreshTokens

CHATS
  POST /api/chats/direct              → CreateDirectChat
  POST /api/chats/group               → CreateGroupChat
  GET  /api/chats/{chat_id}           → GetChat
  GET  /api/chats/mine                → GetUserChats
  POST /api/chats/{chat_id}/members          → AddMember
  DELETE /api/chats/{chat_id}/members/{user_id} → RemoveMember

MESSAGES
  POST /api/messages/send             → SendMessage
  GET  /api/messages/history          → GetHistory (?chat_id, cursor_id, cursor_created_at, limit; omit cursors for latest)
  GET  /api/messages/search           → SearchMessages (?chat_id, query, limit)

REACTIONS
  POST /api/messages/{message_id}/reactions    → AddReaction
  DELETE /api/messages/{message_id}/reactions   → RemoveReaction
  GET  /api/messages/{message_id}/reactions    → GetReactions

ATTACHMENTS
  GET  /api/messages/{message_id}/attachments  → GetAttachments

CHANNELS
  POST /api/channels                  → CreateChannel
  GET  /api/channels/{channel_id}     → GetChannel
  PUT  /api/channels/{channel_id}     → UpdateChannel
  DELETE /api/channels/{channel_id}   → DeleteChannel
  POST /api/channels/{channel_id}/subscribe    → Subscribe
  DELETE /api/channels/{channel_id}/subscribe  → Unsubscribe
  GET  /api/channels/public           → GetPublicChannels
  GET  /api/channels/mine             → GetMyChannels

USERS
  GET  /api/users/{user_id}           → GetProfile
  POST /api/users/online-status       → GetOnlineStatus

NOTIFICATIONS
  GET  /api/notifications             → GetNotifications
  POST /api/notifications/{notification_id}/read     → MarkNotificationRead
  POST /api/notifications/read-all                   → MarkAllNotificationsRead
  GET  /api/notifications/unread-count               → GetUnreadCount

FILES (existing HTTP handler)
  POST /api/files/upload              → Upload
  GET  /api/files/download            → Download

WEBSOCKET
  GET  /ws?token=JWT                  → WebSocket connection
```

## Frontend Structure

### Directory Layout

```
web/
├── index.html
├── package.json
├── vite.config.ts
├── tsconfig.json
└── src/
    ├── main.tsx
    ├── App.tsx
    ├── api/
    │   ├── client.ts            # fetch wrapper: base URL, JWT, error handling
    │   ├── auth.ts              # register, login, logout, refresh
    │   ├── chats.ts             # CRUD, members
    │   ├── messages.ts          # send, history, search, reactions, attachments
    │   ├── channels.ts          # CRUD, subscribe/unsubscribe
    │   ├── users.ts             # profile, online status
    │   ├── files.ts             # upload, download
    │   └── notifications.ts     # list, mark read, count
    ├── ws/
    │   └── socket.ts            # WebSocket connect, reconnect, message dispatch
    ├── store/
    │   ├── authStore.ts         # user, tokens, isAuthenticated
    │   ├── chatStore.ts         # chats list, active chat, members
    │   ├── messageStore.ts      # messages by chatId, pagination cursors
    │   ├── channelStore.ts      # channels, active channel
    │   ├── notificationStore.ts # notifications, unread count
    │   └── uiStore.ts           # sidebar state, modals, search query
    ├── pages/
    │   ├── LoginPage.tsx
    │   ├── RegisterPage.tsx
    │   └── MainPage.tsx         # sidebar + content area
    ├── components/
    │   ├── sidebar/
    │   │   ├── Sidebar.tsx      # left panel container
    │   │   ├── ChatList.tsx     # list of chats/channels
    │   │   ├── ChatItem.tsx     # single chat row
    │   │   └── SearchBar.tsx    # message search
    │   ├── chat/
    │   │   ├── ChatView.tsx     # header + messages + input
    │   │   ├── ChatHeader.tsx   # name, online, actions
    │   │   ├── MessageList.tsx  # infinite scroll with cursor pagination
    │   │   ├── MessageBubble.tsx # text, file, voice variants
    │   │   ├── MessageInput.tsx  # text input + file attach + voice
    │   │   ├── ReactionBar.tsx  # emoji reactions display + add
    │   │   └── TypingIndicator.tsx
    │   ├── channel/
    │   │   ├── ChannelView.tsx  # channel messages (read-only for subscribers)
    │   │   ├── ChannelList.tsx  # public channels browser
    │   │   └── ChannelHeader.tsx
    │   ├── modals/
    │   │   ├── CreateChatModal.tsx
    │   │   ├── CreateGroupModal.tsx
    │   │   ├── CreateChannelModal.tsx
    │   │   └── UserProfileModal.tsx
    │   └── common/
    │       ├── Avatar.tsx       # colored circle with initials
    │       ├── OnlineIndicator.tsx # green dot
    │       ├── NotificationBell.tsx # bell with unread count
    │       └── FilePreview.tsx  # file/image/voice attachment display
    ├── types/
    │   └── index.ts             # User, Chat, Message, Channel, Notification, Reaction, Attachment
    └── hooks/
        ├── useWebSocket.ts      # connect, dispatch to stores, reconnect
        ├── useInfiniteScroll.ts # cursor-based pagination on scroll up
        └── useOnlineStatus.ts   # poll online statuses
```

### Routing

| Path | Component | Auth Required |
|------|-----------|---------------|
| `/login` | LoginPage | No |
| `/register` | RegisterPage | No |
| `/` | MainPage | Yes (redirect to /login) |

### State Management (Zustand)

**authStore:**
- `user: User | null` — текущий пользователь
- `accessToken: string | null` — JWT access token
- `refreshToken: string | null` — JWT refresh token
- `login(email, password)` — вызывает API, сохраняет токены в localStorage
- `register(...)` — вызывает API, редиректит на login
- `logout()` — revoke token, очистка стейта
- `refreshTokens()` — обновление пары токенов

**chatStore:**
- `chats: Chat[]` — список чатов пользователя
- `activeChatId: string | null` — выбранный чат
- `members: Map<chatId, ChatMember[]>` — участники по чату
- `fetchChats()` — загрузка списка чатов
- `createDirect(recipientId)` / `createGroup(name, memberIds)`
- `setActiveChat(chatId)`

**messageStore:**
- `messages: Map<chatId, Message[]>` — сообщения по чатам
- `cursors: Map<chatId, {id, createdAt}>` — курсоры пагинации
- `hasMore: Map<chatId, boolean>`
- `typing: Map<chatId, string[]>` — кто печатает
- `fetchHistory(chatId)` — подгрузка с курсором
- `addMessage(message)` — из WebSocket
- `setTyping(chatId, userId)`

**channelStore:**
- `channels: Channel[]` — подписки пользователя
- `publicChannels: Channel[]` — публичные каналы
- `activeChannelId: string | null`
- `subscribe(channelId)` / `unsubscribe(channelId)`

**notificationStore:**
- `notifications: Notification[]` — список уведомлений
- `unreadCount: number` — счётчик непрочитанных
- `fetchNotifications(limit, unreadOnly)` — загрузка
- `markRead(notificationId)` / `markAllRead()`
- `fetchUnreadCount()`

**uiStore:**
- `sidebarTab: 'chats' | 'channels'`
- `searchQuery: string`
- `activeModal: string | null`

### WebSocket Protocol

Подключение: `ws://localhost:8080/ws?token=<JWT>`

**Исходящие (client → server):**
```json
{"type": "send_message", "payload": {"chat_id": "...", "content": "...", "msg_type": "text"}}
{"type": "typing", "payload": {"chat_id": "..."}}
```

**Входящие (server → client):**
```json
{"type": "new_message", "payload": {"id": "...", "chat_id": "...", "sender_id": "...", "type": "text", "content": "...", "created_at": "..."}}
{"type": "typing", "payload": {"chat_id": "...", "user_id": "..."}}
```

`socket.ts` парсит `type` и диспатчит в соответствующий Zustand store.

**Отличия от REST API:** WebSocket использует `msg_type` (JSON tag), REST использует `type` (proto field name). Ответы на сообщения (`reply_to_id`) доступны только через REST `POST /api/messages/send`, WebSocket handler не поддерживает `reply_to_id`.

Auto-reconnect: exponential backoff (1s, 2s, 4s, 8s, max 30s).

### API Client

`api/client.ts` — fetch wrapper:
- Base URL определяется Vite proxy (пустой prefix, запросы идут на тот же origin)
- `Authorization: Bearer <accessToken>` header из authStore
- При 401 — попытка `refreshTokens()`, повтор запроса
- При повторном 401 — logout, redirect на /login

### UI Design

Тёмная тема в стиле Telegram Desktop:
- Фон: `#0e1621` (chat area), `#17212b` (sidebar, headers)
- Акцент: `#2b5278` (active chat, own messages)
- Входящие: `#182533`
- Текст: `#ffffff` (primary), `#6c7883` (secondary)
- Online: `#3e9152`

CSS Variables в `:root` для единообразия.

### Key Screens

**Login/Register:** Центрированная форма на тёмном фоне. Логотип "Груша" сверху. Валидация на клиенте (email format, password length).

**Main (Chat):** Двухколоночный layout — sidebar (320px) + chat area (flex). Sidebar: поиск сверху, список чатов/каналов (переключение табами). Chat area: header (имя, статус, действия), messages (scroll, infinite load вверх), input (текст, attach, voice).

**Message Bubble:** Rounded corners. Входящие — слева, серый фон. Исходящие — справа, синий фон. Файлы — иконка + имя + размер. Реакции — строка emoji с каунтерами под сообщением.

**Modals:** Create direct chat (поиск юзера), create group (имя + выбор участников), create channel (slug, name, description), user profile.

### Vite Config

```ts
export default defineConfig({
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:8080',
      '/ws': {
        target: 'http://localhost:8080',
        ws: true,
      },
    },
  },
})
```

## Scope Exclusions

- Нет мобильной адаптации (desktop-only для демо)
- Нет пуш-уведомлений (только in-app bell)
- Нет редактирования/удаления сообщений (бэкенд не поддерживает)
- Нет аватарок-картинок (цветные кружки с инициалами)
- Нет dark/light theme toggle (только dark)
- JWT хранится в localStorage (для простоты; httpOnly cookies были бы безопаснее, но усложняют flow)
- Оптимистичные обновления не применяются — UI ждёт подтверждения от сервера через WebSocket
- Upload файлов — двухшаговый flow: (1) отправить сообщение с type `file`/`voice`/`image` через REST, (2) загрузить файл через `POST /api/files/upload?message_id=...`
