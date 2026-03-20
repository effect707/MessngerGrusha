# Grusha Frontend Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Telegram-like React web client for the Grusha messenger backend, plus add grpc-gateway to expose REST endpoints.

**Architecture:** Monorepo approach — React app in `web/` connects to Go backend via grpc-gateway (REST) and native WebSocket. Vite dev server proxies requests to Go on :8080. Zustand for state, CSS Modules for styling.

**Tech Stack:** React 19, TypeScript, Vite, Zustand, React Router v7, CSS Modules, grpc-gateway, Buf CLI

**Spec:** `docs/superpowers/specs/2026-03-20-grusha-frontend-design.md`

---

## File Structure

### Backend changes (grpc-gateway)

| Action | File | Responsibility |
|--------|------|----------------|
| Modify | `buf.yaml` | Add googleapis dependency |
| Modify | `buf.gen.yaml` | Add grpc-gateway plugin |
| Modify | `api/proto/auth/auth.proto` | Add HTTP annotations to all RPCs |
| Modify | `api/proto/chat/chat.proto` | Add HTTP annotations to all RPCs |
| Modify | `api/proto/message/message.proto` | Add HTTP annotations to all RPCs |
| Modify | `api/proto/channel/channel.proto` | Add HTTP annotations to all RPCs |
| Modify | `api/proto/user/user.proto` | Add HTTP annotations to all RPCs |
| Modify | `internal/app/app.go` | Register grpc-gateway mux on HTTP server |
| Modify | `go.mod` | Add grpc-gateway dependency |

### Frontend (new)

| Action | File | Responsibility |
|--------|------|----------------|
| Create | `web/package.json` | Dependencies and scripts |
| Create | `web/tsconfig.json` | TypeScript config |
| Create | `web/vite.config.ts` | Vite config with proxy |
| Create | `web/index.html` | HTML entry |
| Create | `web/src/main.tsx` | React entry point |
| Create | `web/src/App.tsx` | Router setup, auth guard |
| Create | `web/src/App.module.css` | Global app styles |
| Create | `web/src/types/index.ts` | TypeScript interfaces (User, Chat, Message, etc.) |
| Create | `web/src/api/client.ts` | Fetch wrapper with JWT + 401 refresh |
| Create | `web/src/api/auth.ts` | Auth API calls |
| Create | `web/src/api/chats.ts` | Chat API calls |
| Create | `web/src/api/messages.ts` | Message API calls |
| Create | `web/src/api/channels.ts` | Channel API calls |
| Create | `web/src/api/users.ts` | User API calls |
| Create | `web/src/api/files.ts` | File upload/download |
| Create | `web/src/api/notifications.ts` | Notification API calls |
| Create | `web/src/ws/socket.ts` | WebSocket client with reconnect |
| Create | `web/src/store/authStore.ts` | Auth state (user, tokens) |
| Create | `web/src/store/chatStore.ts` | Chats state |
| Create | `web/src/store/messageStore.ts` | Messages + typing state |
| Create | `web/src/store/channelStore.ts` | Channels state |
| Create | `web/src/store/notificationStore.ts` | Notifications state |
| Create | `web/src/store/uiStore.ts` | UI state (sidebar, modals) |
| Create | `web/src/hooks/useWebSocket.ts` | WS connect/dispatch hook |
| Create | `web/src/hooks/useInfiniteScroll.ts` | Cursor-based scroll loading |
| Create | `web/src/pages/LoginPage.tsx` | Login form |
| Create | `web/src/pages/LoginPage.module.css` | Login styles |
| Create | `web/src/pages/RegisterPage.tsx` | Register form |
| Create | `web/src/pages/RegisterPage.module.css` | Register styles |
| Create | `web/src/pages/MainPage.tsx` | Two-column layout |
| Create | `web/src/pages/MainPage.module.css` | Main layout styles |
| Create | `web/src/components/common/Avatar.tsx` | Colored circle with initials |
| Create | `web/src/components/common/Avatar.module.css` | Avatar styles |
| Create | `web/src/components/common/OnlineIndicator.tsx` | Green dot indicator |
| Create | `web/src/components/common/OnlineIndicator.module.css` | Indicator styles |
| Create | `web/src/components/common/NotificationBell.tsx` | Bell with unread count |
| Create | `web/src/components/common/NotificationBell.module.css` | Bell styles |
| Create | `web/src/components/common/FilePreview.tsx` | File/image/voice display |
| Create | `web/src/components/common/FilePreview.module.css` | File preview styles |
| Create | `web/src/components/sidebar/Sidebar.tsx` | Left panel container |
| Create | `web/src/components/sidebar/Sidebar.module.css` | Sidebar styles |
| Create | `web/src/components/sidebar/ChatList.tsx` | List of chats |
| Create | `web/src/components/sidebar/ChatList.module.css` | Chat list styles |
| Create | `web/src/components/sidebar/ChatItem.tsx` | Single chat row |
| Create | `web/src/components/sidebar/ChatItem.module.css` | Chat item styles |
| Create | `web/src/components/sidebar/SearchBar.tsx` | Search input |
| Create | `web/src/components/sidebar/SearchBar.module.css` | Search styles |
| Create | `web/src/components/chat/ChatView.tsx` | Chat area container |
| Create | `web/src/components/chat/ChatView.module.css` | Chat view styles |
| Create | `web/src/components/chat/ChatHeader.tsx` | Chat header |
| Create | `web/src/components/chat/ChatHeader.module.css` | Header styles |
| Create | `web/src/components/chat/MessageList.tsx` | Messages with infinite scroll |
| Create | `web/src/components/chat/MessageList.module.css` | Message list styles |
| Create | `web/src/components/chat/MessageBubble.tsx` | Message bubble (text/file/voice) |
| Create | `web/src/components/chat/MessageBubble.module.css` | Bubble styles |
| Create | `web/src/components/chat/MessageInput.tsx` | Input with file attach |
| Create | `web/src/components/chat/MessageInput.module.css` | Input styles |
| Create | `web/src/components/chat/ReactionBar.tsx` | Emoji reactions |
| Create | `web/src/components/chat/ReactionBar.module.css` | Reaction styles |
| Create | `web/src/components/chat/TypingIndicator.tsx` | "X is typing..." |
| Create | `web/src/components/chat/TypingIndicator.module.css` | Typing styles |
| Create | `web/src/components/channel/ChannelView.tsx` | Channel messages view |
| Create | `web/src/components/channel/ChannelView.module.css` | Channel view styles |
| Create | `web/src/components/channel/ChannelList.tsx` | Public channels browser |
| Create | `web/src/components/channel/ChannelList.module.css` | Channel list styles |
| Create | `web/src/components/channel/ChannelHeader.tsx` | Channel header |
| Create | `web/src/components/channel/ChannelHeader.module.css` | Channel header styles |
| Create | `web/src/components/modals/CreateChatModal.tsx` | Create direct chat modal |
| Create | `web/src/components/modals/CreateChatModal.module.css` | Modal styles |
| Create | `web/src/components/modals/CreateGroupModal.tsx` | Create group modal |
| Create | `web/src/components/modals/CreateGroupModal.module.css` | Modal styles |
| Create | `web/src/components/modals/CreateChannelModal.tsx` | Create channel modal |
| Create | `web/src/components/modals/CreateChannelModal.module.css` | Modal styles |
| Create | `web/src/components/modals/UserProfileModal.tsx` | User profile modal |
| Create | `web/src/components/modals/UserProfileModal.module.css` | Profile styles |

---

## Task 1: grpc-gateway — Proto Annotations + Buf Config

**Files:**
- Modify: `buf.yaml`
- Modify: `buf.gen.yaml`
- Modify: `api/proto/auth/auth.proto`
- Modify: `api/proto/chat/chat.proto`
- Modify: `api/proto/message/message.proto`
- Modify: `api/proto/channel/channel.proto`
- Modify: `api/proto/user/user.proto`

- [ ] **Step 1: Add googleapis dep to buf.yaml**

```yaml
version: v2
modules:
  - path: api/proto
deps:
  - buf.build/googleapis/googleapis
lint:
  use:
    - STANDARD
breaking:
  use:
    - FILE
```

- [ ] **Step 2: Add grpc-gateway plugin to buf.gen.yaml**

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

- [ ] **Step 3: Add HTTP annotations to auth.proto**

```protobuf
syntax = "proto3";

package auth;

option go_package = "github.com/effect707/MessngerGrusha/api/gen/auth";

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
  rpc Logout(LogoutRequest) returns (LogoutResponse) {
    option (google.api.http) = {
      post: "/api/auth/logout"
      body: "*"
    };
  }
  rpc LogoutAll(LogoutAllRequest) returns (LogoutAllResponse) {
    option (google.api.http) = {
      post: "/api/auth/logout-all"
    };
  }
  rpc RefreshTokens(RefreshTokensRequest) returns (RefreshTokensResponse) {
    option (google.api.http) = {
      post: "/api/auth/refresh"
      body: "*"
    };
  }
}

// ... existing messages unchanged
```

- [ ] **Step 4: Add HTTP annotations to chat.proto**

```protobuf
import "google/api/annotations.proto";

service ChatService {
  rpc CreateDirectChat(CreateDirectChatRequest) returns (CreateDirectChatResponse) {
    option (google.api.http) = {
      post: "/api/chats/direct"
      body: "*"
    };
  }
  rpc CreateGroupChat(CreateGroupChatRequest) returns (CreateGroupChatResponse) {
    option (google.api.http) = {
      post: "/api/chats/group"
      body: "*"
    };
  }
  rpc GetChat(GetChatRequest) returns (GetChatResponse) {
    option (google.api.http) = {
      get: "/api/chats/{chat_id}"
    };
  }
  rpc GetUserChats(GetUserChatsRequest) returns (GetUserChatsResponse) {
    option (google.api.http) = {
      get: "/api/chats/mine"
    };
  }
  rpc AddMember(AddMemberRequest) returns (AddMemberResponse) {
    option (google.api.http) = {
      post: "/api/chats/{chat_id}/members"
      body: "*"
    };
  }
  rpc RemoveMember(RemoveMemberRequest) returns (RemoveMemberResponse) {
    option (google.api.http) = {
      delete: "/api/chats/{chat_id}/members/{user_id}"
    };
  }
}
```

- [ ] **Step 5: Add HTTP annotations to message.proto**

```protobuf
import "google/api/annotations.proto";

service MessageService {
  rpc SendMessage(SendMessageRequest) returns (SendMessageResponse) {
    option (google.api.http) = {
      post: "/api/messages/send"
      body: "*"
    };
  }
  rpc GetHistory(GetHistoryRequest) returns (GetHistoryResponse) {
    option (google.api.http) = {
      get: "/api/messages/history"
    };
  }
  rpc SearchMessages(SearchMessagesRequest) returns (SearchMessagesResponse) {
    option (google.api.http) = {
      get: "/api/messages/search"
    };
  }
  rpc AddReaction(AddReactionRequest) returns (AddReactionResponse) {
    option (google.api.http) = {
      post: "/api/messages/{message_id}/reactions"
      body: "*"
    };
  }
  rpc RemoveReaction(RemoveReactionRequest) returns (RemoveReactionResponse) {
    option (google.api.http) = {
      delete: "/api/messages/{message_id}/reactions"
    };
  }
  rpc GetReactions(GetReactionsRequest) returns (GetReactionsResponse) {
    option (google.api.http) = {
      get: "/api/messages/{message_id}/reactions"
    };
  }
  rpc GetAttachments(GetAttachmentsRequest) returns (GetAttachmentsResponse) {
    option (google.api.http) = {
      get: "/api/messages/{message_id}/attachments"
    };
  }
}
```

- [ ] **Step 6: Add HTTP annotations to channel.proto**

```protobuf
import "google/api/annotations.proto";

service ChannelService {
  rpc CreateChannel(CreateChannelRequest) returns (CreateChannelResponse) {
    option (google.api.http) = {
      post: "/api/channels"
      body: "*"
    };
  }
  rpc GetChannel(GetChannelRequest) returns (GetChannelResponse) {
    option (google.api.http) = {
      get: "/api/channels/{channel_id}"
    };
  }
  rpc UpdateChannel(UpdateChannelRequest) returns (UpdateChannelResponse) {
    option (google.api.http) = {
      put: "/api/channels/{channel_id}"
      body: "*"
    };
  }
  rpc DeleteChannel(DeleteChannelRequest) returns (DeleteChannelResponse) {
    option (google.api.http) = {
      delete: "/api/channels/{channel_id}"
    };
  }
  rpc Subscribe(SubscribeRequest) returns (SubscribeResponse) {
    option (google.api.http) = {
      post: "/api/channels/{channel_id}/subscribe"
    };
  }
  rpc Unsubscribe(UnsubscribeRequest) returns (UnsubscribeResponse) {
    option (google.api.http) = {
      delete: "/api/channels/{channel_id}/subscribe"
    };
  }
  rpc GetPublicChannels(GetPublicChannelsRequest) returns (GetPublicChannelsResponse) {
    option (google.api.http) = {
      get: "/api/channels/public"
    };
  }
  rpc GetMyChannels(GetMyChannelsRequest) returns (GetMyChannelsResponse) {
    option (google.api.http) = {
      get: "/api/channels/mine"
    };
  }
}
```

- [ ] **Step 7: Add HTTP annotations to user.proto**

```protobuf
import "google/api/annotations.proto";

service UserService {
  rpc GetProfile(GetProfileRequest) returns (GetProfileResponse) {
    option (google.api.http) = {
      get: "/api/users/{user_id}"
    };
  }
  rpc GetOnlineStatus(GetOnlineStatusRequest) returns (GetOnlineStatusResponse) {
    option (google.api.http) = {
      post: "/api/users/online-status"
      body: "*"
    };
  }
  rpc GetNotifications(GetNotificationsRequest) returns (GetNotificationsResponse) {
    option (google.api.http) = {
      get: "/api/notifications"
    };
  }
  rpc MarkNotificationRead(MarkNotificationReadRequest) returns (MarkNotificationReadResponse) {
    option (google.api.http) = {
      post: "/api/notifications/{notification_id}/read"
    };
  }
  rpc MarkAllNotificationsRead(MarkAllNotificationsReadRequest) returns (MarkAllNotificationsReadResponse) {
    option (google.api.http) = {
      post: "/api/notifications/read-all"
    };
  }
  rpc GetUnreadCount(GetUnreadCountRequest) returns (GetUnreadCountResponse) {
    option (google.api.http) = {
      get: "/api/notifications/unread-count"
    };
  }
}
```

- [ ] **Step 8: Install grpc-gateway protoc plugin**

Run:
```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

- [ ] **Step 9: Update buf deps and generate code**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha
buf dep update
buf generate
```

Expected: New `*.gw.go` files appear in `api/gen/*/` alongside existing `*.pb.go` and `*_grpc.pb.go` files.

- [ ] **Step 10: Verify generated gateway files exist**

Run:
```bash
find api/gen -name "*.gw.go" | sort
```

Expected output (5 files):
```
api/gen/auth/auth.pb.gw.go
api/gen/channel/channel.pb.gw.go
api/gen/chat/chat.pb.gw.go
api/gen/message/message.pb.gw.go
api/gen/user/user.pb.gw.go
```

- [ ] **Step 11: Commit**

```bash
git add buf.yaml buf.gen.yaml api/proto/ api/gen/
git commit -m "feat: add grpc-gateway HTTP annotations to all proto services"
```

---

## Task 2: grpc-gateway — Register Gateway Mux in app.go

**Files:**
- Modify: `internal/app/app.go`
- Modify: `go.mod` / `go.sum`

- [ ] **Step 1: Add grpc-gateway dependency**

Run:
```bash
go get github.com/grpc-ecosystem/grpc-gateway/v2@latest
```

- [ ] **Step 2: Modify app.go — add imports and gateway registration**

In `internal/app/app.go`, add these imports:

```go
import (
    // ... existing imports ...
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc/credentials/insecure"

    gw_auth "github.com/effect707/MessngerGrusha/api/gen/auth"
    gw_channel "github.com/effect707/MessngerGrusha/api/gen/channel"
    gw_chat "github.com/effect707/MessngerGrusha/api/gen/chat"
    gw_msg "github.com/effect707/MessngerGrusha/api/gen/message"
    gw_user "github.com/effect707/MessngerGrusha/api/gen/user"
)
```

Note: The generated gateway code lives in the same packages as existing pb code (`pb_auth`, etc.), so you can reuse the existing import aliases. The `Register...HandlerFromEndpoint` functions are generated into the `*.pb.gw.go` files in the same package. So the actual imports become:

```go
// No new import aliases needed — pb_auth, pb_chat, etc. already import the correct packages
```

After the line `mux.HandleFunc("/api/files/download", fileHandler.Download)` (line 167), add gateway registration:

```go
    // grpc-gateway
    grpcAddr := fmt.Sprintf("localhost:%d", cfg.GRPC.Port)
    gwMux := runtime.NewServeMux()
    gwOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

    if err := pb_auth.RegisterAuthServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, gwOpts); err != nil {
        return nil, fmt.Errorf("register auth gateway: %w", err)
    }
    if err := pb_chat.RegisterChatServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, gwOpts); err != nil {
        return nil, fmt.Errorf("register chat gateway: %w", err)
    }
    if err := pb_msg.RegisterMessageServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, gwOpts); err != nil {
        return nil, fmt.Errorf("register message gateway: %w", err)
    }
    if err := pb_channel.RegisterChannelServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, gwOpts); err != nil {
        return nil, fmt.Errorf("register channel gateway: %w", err)
    }
    if err := pb_user.RegisterUserServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, gwOpts); err != nil {
        return nil, fmt.Errorf("register user gateway: %w", err)
    }

    mux.Handle("/api/", gwMux)
```

Also add necessary imports:

```go
"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
"google.golang.org/grpc/credentials/insecure"
```

- [ ] **Step 3: Verify build**

Run:
```bash
go build ./cmd/grusha
```

Expected: Successful build, no errors.

- [ ] **Step 4: Run existing tests to ensure nothing breaks**

Run:
```bash
go test -race ./...
```

Expected: All existing tests pass.

- [ ] **Step 5: Commit**

```bash
git add internal/app/app.go go.mod go.sum
git commit -m "feat: register grpc-gateway mux on HTTP server for REST API"
```

---

## Task 3: React App Scaffold

**Files:**
- Create: `web/package.json`, `web/tsconfig.json`, `web/vite.config.ts`, `web/index.html`
- Create: `web/src/main.tsx`, `web/src/App.tsx`, `web/src/App.module.css`
- Create: `web/src/types/index.ts`

- [ ] **Step 1: Create web/package.json**

```json
{
  "name": "grusha-web",
  "private": true,
  "version": "0.0.1",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "^19.1.0",
    "react-dom": "^19.1.0",
    "react-router-dom": "^7.6.0",
    "zustand": "^5.0.0"
  },
  "devDependencies": {
    "@types/react": "^19.1.0",
    "@types/react-dom": "^19.1.0",
    "@vitejs/plugin-react": "^4.4.0",
    "typescript": "~5.8.0",
    "vite": "^6.3.0"
  }
}
```

- [ ] **Step 2: Create web/tsconfig.json**

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "lib": ["ES2023", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "isolatedModules": true,
    "moduleDetection": "force",
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "noUncheckedSideEffectImports": true
  },
  "include": ["src"]
}
```

- [ ] **Step 3: Create web/vite.config.ts**

```ts
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
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

- [ ] **Step 4: Create web/index.html**

```html
<!DOCTYPE html>
<html lang="ru">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Груша</title>
    <style>
      :root {
        --bg-primary: #0e1621;
        --bg-secondary: #17212b;
        --bg-active: #2b5278;
        --bg-incoming: #182533;
        --bg-input: #242f3d;
        --text-primary: #ffffff;
        --text-secondary: #6c7883;
        --accent: #5288c1;
        --online: #3e9152;
        --border: #101921;
        --danger: #e17076;
      }
      * { margin: 0; padding: 0; box-sizing: border-box; }
      body {
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        background: var(--bg-primary);
        color: var(--text-primary);
        overflow: hidden;
        height: 100vh;
      }
      #root { height: 100%; }
    </style>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

- [ ] **Step 5: Create web/src/types/index.ts**

```ts
export interface User {
  id: string
  username: string
  email: string
  display_name: string
  avatar_url: string
  bio: string
  created_at: string
}

export interface Chat {
  id: string
  type: 'direct' | 'group'
  name: string
  avatar_url: string
  created_by: string
  created_at: string
}

export interface ChatMember {
  chat_id: string
  user_id: string
  role: 'admin' | 'member'
  joined_at: string
}

export interface Message {
  id: string
  chat_id: string
  sender_id: string
  type: 'text' | 'image' | 'file' | 'voice' | 'system'
  content: string
  reply_to_id?: string
  is_edited: boolean
  created_at: string
  updated_at: string
}

export interface Channel {
  id: string
  slug: string
  name: string
  description: string
  avatar_url: string
  owner_id: string
  is_private: boolean
  created_at: string
}

export interface Notification {
  id: string
  type: string
  payload: string
  is_read: boolean
  created_at: string
}

export interface Reaction {
  message_id: string
  user_id: string
  emoji: string
  created_at: string
}

export interface Attachment {
  id: string
  message_id: string
  file_name: string
  file_size: number
  mime_type: string
  duration_ms?: number
  created_at: string
}
```

- [ ] **Step 6: Create web/src/main.tsx**

```tsx
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { App } from './App'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </StrictMode>,
)
```

- [ ] **Step 7: Create web/src/App.tsx (placeholder with routes)**

```tsx
import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'

function LoginPlaceholder() {
  return <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh', color: 'var(--text-secondary)' }}>Login Page (TODO)</div>
}

function RegisterPlaceholder() {
  return <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh', color: 'var(--text-secondary)' }}>Register Page (TODO)</div>
}

function MainPlaceholder() {
  return <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh', color: 'var(--text-secondary)' }}>Main Page (TODO)</div>
}

function AuthGuard({ children }: { children: React.ReactNode }) {
  const accessToken = useAuthStore((s) => s.accessToken)
  if (!accessToken) return <Navigate to="/login" replace />
  return <>{children}</>
}

export function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPlaceholder />} />
      <Route path="/register" element={<RegisterPlaceholder />} />
      <Route path="/" element={<AuthGuard><MainPlaceholder /></AuthGuard>} />
    </Routes>
  )
}
```

- [ ] **Step 8: Install deps and verify dev server starts**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha/web
npm install
npx vite --host 127.0.0.1 &
sleep 3
curl -s http://127.0.0.1:5173 | head -5
kill %1
```

Expected: HTML page with `<title>Груша</title>`.

- [ ] **Step 9: Commit**

```bash
git add web/
git commit -m "feat: scaffold React app with Vite, routing, and TypeScript types"
```

---

## Task 4: API Client + Auth Store + Auth API

**Files:**
- Create: `web/src/api/client.ts`
- Create: `web/src/api/auth.ts`
- Create: `web/src/store/authStore.ts`

- [ ] **Step 1: Create web/src/api/client.ts**

```ts
let getAccessToken: () => string | null = () => null
let onUnauthorized: () => void = () => {}

export function setAuthCallbacks(
  tokenGetter: () => string | null,
  unauthorizedHandler: () => void,
) {
  getAccessToken = tokenGetter
  onUnauthorized = unauthorizedHandler
}

export async function apiRequest<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const token = getAccessToken()
  const headers: Record<string, string> = {
    ...((options.headers as Record<string, string>) || {}),
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json'
  }

  const res = await fetch(path, { ...options, headers })

  if (res.status === 401) {
    onUnauthorized()
    throw new Error('Unauthorized')
  }

  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `HTTP ${res.status}`)
  }

  if (res.status === 204 || res.headers.get('content-length') === '0') {
    return {} as T
  }

  return res.json()
}
```

- [ ] **Step 2: Create web/src/api/auth.ts**

```ts
import { apiRequest } from './client'
import type { User } from '../types'

interface LoginResponse {
  access_token: string
  refresh_token: string
}

interface RegisterResponse {
  user: User
}

interface RefreshResponse {
  access_token: string
  refresh_token: string
}

export const authApi = {
  register(username: string, email: string, password: string, displayName: string) {
    return apiRequest<RegisterResponse>('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify({ username, email, password, display_name: displayName }),
    })
  },

  login(email: string, password: string) {
    return apiRequest<LoginResponse>('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
  },

  logout(refreshToken: string) {
    return apiRequest<object>('/api/auth/logout', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken }),
    })
  },

  refresh(refreshToken: string) {
    return apiRequest<RefreshResponse>('/api/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken }),
    })
  },
}
```

- [ ] **Step 3: Create web/src/store/authStore.ts**

```ts
import { create } from 'zustand'
import type { User } from '../types'
import { authApi } from '../api/auth'
import { setAuthCallbacks } from '../api/client'

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  login: (email: string, password: string) => Promise<void>
  register: (username: string, email: string, password: string, displayName: string) => Promise<void>
  logout: () => Promise<void>
  restoreSession: () => void
}

export const useAuthStore = create<AuthState>((set, get) => {
  // Wire up the API client callbacks
  setAuthCallbacks(
    () => get().accessToken,
    () => {
      set({ user: null, accessToken: null, refreshToken: null })
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
      window.location.href = '/login'
    },
  )

  return {
    user: null,
    accessToken: null,
    refreshToken: null,

    async login(email, password) {
      const res = await authApi.login(email, password)
      localStorage.setItem('accessToken', res.access_token)
      localStorage.setItem('refreshToken', res.refresh_token)
      set({ accessToken: res.access_token, refreshToken: res.refresh_token })
    },

    async register(username, email, password, displayName) {
      await authApi.register(username, email, password, displayName)
    },

    async logout() {
      const rt = get().refreshToken
      if (rt) {
        await authApi.logout(rt).catch(() => {})
      }
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
      set({ user: null, accessToken: null, refreshToken: null })
    },

    restoreSession() {
      const accessToken = localStorage.getItem('accessToken')
      const refreshToken = localStorage.getItem('refreshToken')
      if (accessToken) {
        set({ accessToken, refreshToken })
      }
    },
  }
})
```

- [ ] **Step 4: Update App.tsx to call restoreSession on mount**

In `web/src/App.tsx`, add:

```tsx
import { useEffect } from 'react'

// Inside App component, before return:
useEffect(() => {
  useAuthStore.getState().restoreSession()
}, [])
```

- [ ] **Step 5: Verify TypeScript compiles**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha/web
npx tsc --noEmit
```

Expected: No errors.

- [ ] **Step 6: Commit**

```bash
git add web/src/api/ web/src/store/authStore.ts web/src/App.tsx
git commit -m "feat: add API client with JWT auth, auth API layer, and auth store"
```

---

## Task 5: Login + Register Pages

**Files:**
- Create: `web/src/pages/LoginPage.tsx`, `web/src/pages/LoginPage.module.css`
- Create: `web/src/pages/RegisterPage.tsx`, `web/src/pages/RegisterPage.module.css`
- Modify: `web/src/App.tsx` — replace placeholders

- [ ] **Step 1: Create web/src/pages/LoginPage.module.css**

```css
.container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  background: var(--bg-primary);
}

.form {
  width: 360px;
  padding: 40px;
  background: var(--bg-secondary);
  border-radius: 12px;
}

.logo {
  text-align: center;
  font-size: 32px;
  margin-bottom: 8px;
}

.title {
  text-align: center;
  font-size: 20px;
  font-weight: 500;
  margin-bottom: 24px;
  color: var(--text-primary);
}

.input {
  width: 100%;
  padding: 12px 16px;
  margin-bottom: 12px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
}

.input:focus {
  border-color: var(--accent);
}

.button {
  width: 100%;
  padding: 12px;
  margin-top: 8px;
  background: var(--accent);
  color: var(--text-primary);
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
}

.button:hover {
  opacity: 0.9;
}

.button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.link {
  display: block;
  text-align: center;
  margin-top: 16px;
  color: var(--accent);
  font-size: 14px;
  text-decoration: none;
}

.error {
  color: var(--danger);
  font-size: 13px;
  text-align: center;
  margin-bottom: 12px;
}
```

- [ ] **Step 2: Create web/src/pages/LoginPage.tsx**

```tsx
import { useState, type FormEvent } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'
import styles from './LoginPage.module.css'

export function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const login = useAuthStore((s) => s.login)
  const navigate = useNavigate()

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await login(email, password)
      navigate('/', { replace: true })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className={styles.container}>
      <form className={styles.form} onSubmit={handleSubmit}>
        <div className={styles.logo}>🍐</div>
        <h1 className={styles.title}>Груша</h1>
        {error && <div className={styles.error}>{error}</div>}
        <input
          className={styles.input}
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <input
          className={styles.input}
          type="password"
          placeholder="Пароль"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button className={styles.button} type="submit" disabled={loading}>
          {loading ? 'Вход...' : 'Войти'}
        </button>
        <Link className={styles.link} to="/register">
          Нет аккаунта? Зарегистрироваться
        </Link>
      </form>
    </div>
  )
}
```

- [ ] **Step 3: Create web/src/pages/RegisterPage.module.css**

Same as `LoginPage.module.css` — copy it:

```bash
cp web/src/pages/LoginPage.module.css web/src/pages/RegisterPage.module.css
```

- [ ] **Step 4: Create web/src/pages/RegisterPage.tsx**

```tsx
import { useState, type FormEvent } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'
import styles from './RegisterPage.module.css'

export function RegisterPage() {
  const [username, setUsername] = useState('')
  const [displayName, setDisplayName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const register = useAuthStore((s) => s.register)
  const navigate = useNavigate()

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await register(username, email, password, displayName)
      navigate('/login', { replace: true })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className={styles.container}>
      <form className={styles.form} onSubmit={handleSubmit}>
        <div className={styles.logo}>🍐</div>
        <h1 className={styles.title}>Регистрация</h1>
        {error && <div className={styles.error}>{error}</div>}
        <input
          className={styles.input}
          type="text"
          placeholder="Имя пользователя"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
        <input
          className={styles.input}
          type="text"
          placeholder="Отображаемое имя"
          value={displayName}
          onChange={(e) => setDisplayName(e.target.value)}
          required
        />
        <input
          className={styles.input}
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <input
          className={styles.input}
          type="password"
          placeholder="Пароль"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          minLength={6}
        />
        <button className={styles.button} type="submit" disabled={loading}>
          {loading ? 'Регистрация...' : 'Зарегистрироваться'}
        </button>
        <Link className={styles.link} to="/login">
          Уже есть аккаунт? Войти
        </Link>
      </form>
    </div>
  )
}
```

- [ ] **Step 5: Update App.tsx — replace placeholders with real pages**

```tsx
import { useEffect } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import { LoginPage } from './pages/LoginPage'
import { RegisterPage } from './pages/RegisterPage'

function MainPlaceholder() {
  return <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh', color: 'var(--text-secondary)' }}>Main Page (TODO)</div>
}

function AuthGuard({ children }: { children: React.ReactNode }) {
  const accessToken = useAuthStore((s) => s.accessToken)
  if (!accessToken) return <Navigate to="/login" replace />
  return <>{children}</>
}

export function App() {
  useEffect(() => {
    useAuthStore.getState().restoreSession()
  }, [])

  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/" element={<AuthGuard><MainPlaceholder /></AuthGuard>} />
    </Routes>
  )
}
```

- [ ] **Step 6: Verify TypeScript compiles**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha/web && npx tsc --noEmit
```

- [ ] **Step 7: Commit**

```bash
git add web/src/pages/ web/src/App.tsx
git commit -m "feat: add Login and Register pages with Telegram dark theme"
```

---

## Task 6: Remaining API Layers + Stores

**Files:**
- Create: `web/src/api/chats.ts`, `web/src/api/messages.ts`, `web/src/api/channels.ts`, `web/src/api/users.ts`, `web/src/api/files.ts`, `web/src/api/notifications.ts`
- Create: `web/src/store/chatStore.ts`, `web/src/store/messageStore.ts`, `web/src/store/channelStore.ts`, `web/src/store/notificationStore.ts`, `web/src/store/uiStore.ts`

- [ ] **Step 1: Create web/src/api/chats.ts**

```ts
import { apiRequest } from './client'
import type { Chat } from '../types'

export const chatsApi = {
  createDirect(recipientId: string) {
    return apiRequest<{ chat: Chat }>('/api/chats/direct', {
      method: 'POST',
      body: JSON.stringify({ recipient_id: recipientId }),
    })
  },

  createGroup(name: string, memberIds: string[]) {
    return apiRequest<{ chat: Chat }>('/api/chats/group', {
      method: 'POST',
      body: JSON.stringify({ name, member_ids: memberIds }),
    })
  },

  getChat(chatId: string) {
    return apiRequest<{ chat: Chat }>(`/api/chats/${chatId}`)
  },

  getUserChats() {
    return apiRequest<{ chats: Chat[] }>('/api/chats/mine')
  },

  addMember(chatId: string, userId: string) {
    return apiRequest<object>(`/api/chats/${chatId}/members`, {
      method: 'POST',
      body: JSON.stringify({ user_id: userId }),
    })
  },

  removeMember(chatId: string, userId: string) {
    return apiRequest<object>(`/api/chats/${chatId}/members/${userId}`, {
      method: 'DELETE',
    })
  },
}
```

- [ ] **Step 2: Create web/src/api/messages.ts**

```ts
import { apiRequest } from './client'
import type { Message, Reaction, Attachment } from '../types'

interface HistoryResponse {
  messages: Message[]
  has_more: boolean
  next_cursor_id?: string
  next_cursor_created_at?: string
}

export const messagesApi = {
  send(chatId: string, type: string, content: string, replyToId?: string) {
    return apiRequest<{ message: Message }>('/api/messages/send', {
      method: 'POST',
      body: JSON.stringify({ chat_id: chatId, type, content, reply_to_id: replyToId }),
    })
  },

  getHistory(chatId: string, limit: number, cursorId?: string, cursorCreatedAt?: string) {
    const params = new URLSearchParams({ chat_id: chatId, limit: String(limit) })
    if (cursorId) params.set('cursor_id', cursorId)
    if (cursorCreatedAt) params.set('cursor_created_at', cursorCreatedAt)
    return apiRequest<HistoryResponse>(`/api/messages/history?${params}`)
  },

  search(chatId: string, query: string, limit: number) {
    const params = new URLSearchParams({ chat_id: chatId, query, limit: String(limit) })
    return apiRequest<{ messages: Message[] }>(`/api/messages/search?${params}`)
  },

  addReaction(messageId: string, emoji: string) {
    return apiRequest<object>(`/api/messages/${messageId}/reactions`, {
      method: 'POST',
      body: JSON.stringify({ emoji }),
    })
  },

  removeReaction(messageId: string, emoji: string) {
    const params = new URLSearchParams({ emoji })
    return apiRequest<object>(`/api/messages/${messageId}/reactions?${params}`, {
      method: 'DELETE',
    })
  },

  getReactions(messageId: string) {
    return apiRequest<{ reactions: Reaction[] }>(`/api/messages/${messageId}/reactions`)
  },

  getAttachments(messageId: string) {
    return apiRequest<{ attachments: Attachment[] }>(`/api/messages/${messageId}/attachments`)
  },
}
```

- [ ] **Step 3: Create web/src/api/channels.ts**

```ts
import { apiRequest } from './client'
import type { Channel } from '../types'

export const channelsApi = {
  create(slug: string, name: string, description: string, isPrivate: boolean) {
    return apiRequest<{ channel: Channel }>('/api/channels', {
      method: 'POST',
      body: JSON.stringify({ slug, name, description, is_private: isPrivate }),
    })
  },

  getChannel(channelId: string) {
    return apiRequest<{ channel: Channel }>(`/api/channels/${channelId}`)
  },

  update(channelId: string, name: string, description: string, isPrivate: boolean) {
    return apiRequest<{ channel: Channel }>(`/api/channels/${channelId}`, {
      method: 'PUT',
      body: JSON.stringify({ name, description, is_private: isPrivate }),
    })
  },

  delete(channelId: string) {
    return apiRequest<object>(`/api/channels/${channelId}`, { method: 'DELETE' })
  },

  subscribe(channelId: string) {
    return apiRequest<object>(`/api/channels/${channelId}/subscribe`, { method: 'POST' })
  },

  unsubscribe(channelId: string) {
    return apiRequest<object>(`/api/channels/${channelId}/subscribe`, { method: 'DELETE' })
  },

  getPublic(limit: number) {
    return apiRequest<{ channels: Channel[] }>(`/api/channels/public?limit=${limit}`)
  },

  getMine() {
    return apiRequest<{ channels: Channel[] }>('/api/channels/mine')
  },
}
```

- [ ] **Step 4: Create web/src/api/users.ts**

```ts
import { apiRequest } from './client'
import type { User } from '../types'

export const usersApi = {
  getProfile(userId: string) {
    return apiRequest<{ user: User }>(`/api/users/${userId}`)
  },

  getOnlineStatus(userIds: string[]) {
    return apiRequest<{ statuses: Record<string, boolean> }>('/api/users/online-status', {
      method: 'POST',
      body: JSON.stringify({ user_ids: userIds }),
    })
  },
}
```

- [ ] **Step 5: Create web/src/api/files.ts**

```ts
import type { Attachment } from '../types'

export const filesApi = {
  async upload(
    messageId: string,
    file: File,
    durationMs?: number,
    token?: string,
  ): Promise<Attachment> {
    const form = new FormData()
    form.append('file', file)
    form.append('message_id', messageId)
    if (durationMs !== undefined) {
      form.append('duration_ms', String(durationMs))
    }

    const headers: Record<string, string> = {}
    if (token) headers['Authorization'] = `Bearer ${token}`

    const res = await fetch('/api/files/upload', {
      method: 'POST',
      headers,
      body: form,
    })

    if (!res.ok) throw new Error('Upload failed')
    return res.json()
  },

  downloadUrl(attachmentId: string, token: string) {
    return `/api/files/download?id=${attachmentId}&token=${token}`
  },
}
```

- [ ] **Step 6: Create web/src/api/notifications.ts**

```ts
import { apiRequest } from './client'
import type { Notification } from '../types'

export const notificationsApi = {
  getAll(limit: number, unreadOnly: boolean = false) {
    const params = new URLSearchParams({ limit: String(limit) })
    if (unreadOnly) params.set('unread_only', 'true')
    return apiRequest<{ notifications: Notification[] }>(`/api/notifications?${params}`)
  },

  markRead(notificationId: string) {
    return apiRequest<object>(`/api/notifications/${notificationId}/read`, { method: 'POST' })
  },

  markAllRead() {
    return apiRequest<object>('/api/notifications/read-all', { method: 'POST' })
  },

  getUnreadCount() {
    return apiRequest<{ count: number }>('/api/notifications/unread-count')
  },
}
```

- [ ] **Step 7: Create web/src/store/chatStore.ts**

```ts
import { create } from 'zustand'
import type { Chat } from '../types'
import { chatsApi } from '../api/chats'

interface ChatState {
  chats: Chat[]
  activeChatId: string | null
  fetchChats: () => Promise<void>
  setActiveChat: (chatId: string | null) => void
  createDirect: (recipientId: string) => Promise<Chat>
  createGroup: (name: string, memberIds: string[]) => Promise<Chat>
}

export const useChatStore = create<ChatState>((set) => ({
  chats: [],
  activeChatId: null,

  async fetchChats() {
    const res = await chatsApi.getUserChats()
    set({ chats: res.chats || [] })
  },

  setActiveChat(chatId) {
    set({ activeChatId: chatId })
  },

  async createDirect(recipientId) {
    const res = await chatsApi.createDirect(recipientId)
    set((s) => ({ chats: [res.chat, ...s.chats], activeChatId: res.chat.id }))
    return res.chat
  },

  async createGroup(name, memberIds) {
    const res = await chatsApi.createGroup(name, memberIds)
    set((s) => ({ chats: [res.chat, ...s.chats], activeChatId: res.chat.id }))
    return res.chat
  },
}))
```

- [ ] **Step 8: Create web/src/store/messageStore.ts**

```ts
import { create } from 'zustand'
import type { Message } from '../types'
import { messagesApi } from '../api/messages'

interface Cursor {
  id: string
  createdAt: string
}

interface MessageState {
  messages: Record<string, Message[]>
  cursors: Record<string, Cursor>
  hasMore: Record<string, boolean>
  typing: Record<string, string[]>
  fetchHistory: (chatId: string) => Promise<void>
  addMessage: (message: Message) => void
  setTyping: (chatId: string, userId: string) => void
  clearTyping: (chatId: string, userId: string) => void
}

export const useMessageStore = create<MessageState>((set, get) => ({
  messages: {},
  cursors: {},
  hasMore: {},
  typing: {},

  async fetchHistory(chatId) {
    const cursor = get().cursors[chatId]
    const res = await messagesApi.getHistory(chatId, 30, cursor?.id, cursor?.createdAt)
    const incoming = res.messages || []
    set((s) => ({
      messages: {
        ...s.messages,
        [chatId]: [...(s.messages[chatId] || []), ...incoming],
      },
      hasMore: { ...s.hasMore, [chatId]: res.has_more },
      cursors: res.next_cursor_id
        ? {
            ...s.cursors,
            [chatId]: { id: res.next_cursor_id, createdAt: res.next_cursor_created_at! },
          }
        : s.cursors,
    }))
  },

  addMessage(message) {
    set((s) => ({
      messages: {
        ...s.messages,
        [message.chat_id]: [message, ...(s.messages[message.chat_id] || [])],
      },
    }))
  },

  setTyping(chatId, userId) {
    set((s) => {
      const current = s.typing[chatId] || []
      if (current.includes(userId)) return s
      return { typing: { ...s.typing, [chatId]: [...current, userId] } }
    })
    // Auto-clear after 3 seconds
    setTimeout(() => get().clearTyping(chatId, userId), 3000)
  },

  clearTyping(chatId, userId) {
    set((s) => ({
      typing: {
        ...s.typing,
        [chatId]: (s.typing[chatId] || []).filter((id) => id !== userId),
      },
    }))
  },
}))
```

- [ ] **Step 9: Create web/src/store/channelStore.ts**

```ts
import { create } from 'zustand'
import type { Channel } from '../types'
import { channelsApi } from '../api/channels'

interface ChannelState {
  channels: Channel[]
  publicChannels: Channel[]
  activeChannelId: string | null
  fetchMyChannels: () => Promise<void>
  fetchPublicChannels: () => Promise<void>
  setActiveChannel: (id: string | null) => void
  subscribe: (channelId: string) => Promise<void>
  unsubscribe: (channelId: string) => Promise<void>
}

export const useChannelStore = create<ChannelState>((set) => ({
  channels: [],
  publicChannels: [],
  activeChannelId: null,

  async fetchMyChannels() {
    const res = await channelsApi.getMine()
    set({ channels: res.channels || [] })
  },

  async fetchPublicChannels() {
    const res = await channelsApi.getPublic(50)
    set({ publicChannels: res.channels || [] })
  },

  setActiveChannel(id) {
    set({ activeChannelId: id })
  },

  async subscribe(channelId) {
    await channelsApi.subscribe(channelId)
  },

  async unsubscribe(channelId) {
    await channelsApi.unsubscribe(channelId)
    set((s) => ({ channels: s.channels.filter((c) => c.id !== channelId) }))
  },
}))
```

- [ ] **Step 10: Create web/src/store/notificationStore.ts**

```ts
import { create } from 'zustand'
import type { Notification } from '../types'
import { notificationsApi } from '../api/notifications'

interface NotificationState {
  notifications: Notification[]
  unreadCount: number
  fetchNotifications: (limit?: number) => Promise<void>
  fetchUnreadCount: () => Promise<void>
  markRead: (id: string) => Promise<void>
  markAllRead: () => Promise<void>
}

export const useNotificationStore = create<NotificationState>((set) => ({
  notifications: [],
  unreadCount: 0,

  async fetchNotifications(limit = 20) {
    const res = await notificationsApi.getAll(limit)
    set({ notifications: res.notifications || [] })
  },

  async fetchUnreadCount() {
    const res = await notificationsApi.getUnreadCount()
    set({ unreadCount: Number(res.count) || 0 })
  },

  async markRead(id) {
    await notificationsApi.markRead(id)
    set((s) => ({
      notifications: s.notifications.map((n) => (n.id === id ? { ...n, is_read: true } : n)),
      unreadCount: Math.max(0, s.unreadCount - 1),
    }))
  },

  async markAllRead() {
    await notificationsApi.markAllRead()
    set((s) => ({
      notifications: s.notifications.map((n) => ({ ...n, is_read: true })),
      unreadCount: 0,
    }))
  },
}))
```

- [ ] **Step 11: Create web/src/store/uiStore.ts**

```ts
import { create } from 'zustand'

type SidebarTab = 'chats' | 'channels'

interface UIState {
  sidebarTab: SidebarTab
  searchQuery: string
  activeModal: string | null
  setSidebarTab: (tab: SidebarTab) => void
  setSearchQuery: (q: string) => void
  openModal: (modal: string) => void
  closeModal: () => void
}

export const useUIStore = create<UIState>((set) => ({
  sidebarTab: 'chats',
  searchQuery: '',
  activeModal: null,

  setSidebarTab(tab) { set({ sidebarTab: tab }) },
  setSearchQuery(q) { set({ searchQuery: q }) },
  openModal(modal) { set({ activeModal: modal }) },
  closeModal() { set({ activeModal: null }) },
}))
```

- [ ] **Step 12: Verify TypeScript compiles**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha/web && npx tsc --noEmit
```

- [ ] **Step 13: Commit**

```bash
git add web/src/api/ web/src/store/
git commit -m "feat: add all API layers and Zustand stores for chats, messages, channels, notifications"
```

---

## Task 7: WebSocket Client + Hook

**Files:**
- Create: `web/src/ws/socket.ts`
- Create: `web/src/hooks/useWebSocket.ts`

- [ ] **Step 1: Create web/src/ws/socket.ts**

```ts
type MessageHandler = (data: unknown) => void

export class GrushaSocket {
  private ws: WebSocket | null = null
  private handlers = new Map<string, MessageHandler[]>()
  private reconnectDelay = 1000
  private maxReconnectDelay = 30000
  private token: string
  private closed = false

  constructor(token: string) {
    this.token = token
  }

  connect() {
    this.closed = false
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const url = `${protocol}//${window.location.host}/ws?token=${this.token}`
    this.ws = new WebSocket(url)

    this.ws.onopen = () => {
      this.reconnectDelay = 1000
    }

    this.ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data) as { type: string; payload: unknown }
        const handlers = this.handlers.get(msg.type)
        if (handlers) {
          handlers.forEach((h) => h(msg.payload))
        }
      } catch {
        // ignore malformed messages
      }
    }

    this.ws.onclose = () => {
      if (!this.closed) {
        setTimeout(() => this.connect(), this.reconnectDelay)
        this.reconnectDelay = Math.min(this.reconnectDelay * 2, this.maxReconnectDelay)
      }
    }
  }

  on(type: string, handler: MessageHandler) {
    const existing = this.handlers.get(type) || []
    this.handlers.set(type, [...existing, handler])
  }

  off(type: string, handler: MessageHandler) {
    const existing = this.handlers.get(type) || []
    this.handlers.set(type, existing.filter((h) => h !== handler))
  }

  send(type: string, payload: unknown) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ type, payload }))
    }
  }

  disconnect() {
    this.closed = true
    this.ws?.close()
    this.ws = null
  }
}
```

- [ ] **Step 2: Create web/src/hooks/useWebSocket.ts**

```ts
import { useEffect, useRef } from 'react'
import { GrushaSocket } from '../ws/socket'
import { useAuthStore } from '../store/authStore'
import { useMessageStore } from '../store/messageStore'
import type { Message } from '../types'

export function useWebSocket() {
  const socketRef = useRef<GrushaSocket | null>(null)
  const accessToken = useAuthStore((s) => s.accessToken)
  const addMessage = useMessageStore((s) => s.addMessage)
  const setTyping = useMessageStore((s) => s.setTyping)

  useEffect(() => {
    if (!accessToken) return

    const socket = new GrushaSocket(accessToken)
    socketRef.current = socket

    socket.on('new_message', (payload) => {
      addMessage(payload as Message)
    })

    socket.on('typing', (payload) => {
      const p = payload as { chat_id: string; user_id: string }
      setTyping(p.chat_id, p.user_id)
    })

    socket.connect()

    return () => {
      socket.disconnect()
      socketRef.current = null
    }
  }, [accessToken, addMessage, setTyping])

  return socketRef
}
```

- [ ] **Step 3: Create web/src/hooks/useInfiniteScroll.ts**

```ts
import { useCallback, useRef } from 'react'

export function useInfiniteScroll(
  onLoadMore: () => Promise<void>,
  hasMore: boolean,
) {
  const loading = useRef(false)
  const containerRef = useRef<HTMLDivElement>(null)

  const handleScroll = useCallback(async () => {
    const el = containerRef.current
    if (!el || loading.current || !hasMore) return

    // Load more when scrolled near top (messages load upward)
    if (el.scrollTop < 100) {
      loading.current = true
      const prevHeight = el.scrollHeight
      await onLoadMore()
      // Maintain scroll position after prepending messages
      const newHeight = el.scrollHeight
      el.scrollTop = newHeight - prevHeight
      loading.current = false
    }
  }, [onLoadMore, hasMore])

  return { containerRef, handleScroll }
}
```

- [ ] **Step 4: Verify TypeScript compiles**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha/web && npx tsc --noEmit
```

- [ ] **Step 5: Commit**

```bash
git add web/src/ws/ web/src/hooks/
git commit -m "feat: add WebSocket client with auto-reconnect and hooks"
```

---

## Task 8: Common Components

**Files:**
- Create: `web/src/components/common/Avatar.tsx` + `.module.css`
- Create: `web/src/components/common/OnlineIndicator.tsx` + `.module.css`
- Create: `web/src/components/common/NotificationBell.tsx` + `.module.css`
- Create: `web/src/components/common/FilePreview.tsx` + `.module.css`

- [ ] **Step 1: Create Avatar component**

`web/src/components/common/Avatar.module.css`:
```css
.avatar {
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 500;
  flex-shrink: 0;
  color: #fff;
}
```

`web/src/components/common/Avatar.tsx`:
```tsx
import styles from './Avatar.module.css'

const COLORS = ['#5288c1', '#e17076', '#67a551', '#e4ae3a', '#7b72e9', '#ee7aae', '#6ec9cb', '#faa774']

function hashColor(str: string): string {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash)
  }
  return COLORS[Math.abs(hash) % COLORS.length]
}

interface Props {
  name: string
  size?: number
}

export function Avatar({ name, size = 48 }: Props) {
  const initials = name
    .split(' ')
    .map((w) => w[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)

  return (
    <div
      className={styles.avatar}
      style={{
        width: size,
        height: size,
        fontSize: size * 0.38,
        background: hashColor(name),
      }}
    >
      {initials || '?'}
    </div>
  )
}
```

- [ ] **Step 2: Create OnlineIndicator component**

`web/src/components/common/OnlineIndicator.module.css`:
```css
.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: var(--online);
  border: 2px solid var(--bg-secondary);
  position: absolute;
  bottom: 0;
  right: 0;
}
```

`web/src/components/common/OnlineIndicator.tsx`:
```tsx
import styles from './OnlineIndicator.module.css'

export function OnlineIndicator() {
  return <div className={styles.dot} />
}
```

- [ ] **Step 3: Create NotificationBell component**

`web/src/components/common/NotificationBell.module.css`:
```css
.bell {
  position: relative;
  cursor: pointer;
  font-size: 20px;
  color: var(--text-secondary);
}

.badge {
  position: absolute;
  top: -6px;
  right: -8px;
  background: var(--danger);
  color: #fff;
  font-size: 11px;
  min-width: 18px;
  height: 18px;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 4px;
}
```

`web/src/components/common/NotificationBell.tsx`:
```tsx
import { useNotificationStore } from '../../store/notificationStore'
import styles from './NotificationBell.module.css'

interface Props {
  onClick: () => void
}

export function NotificationBell({ onClick }: Props) {
  const unreadCount = useNotificationStore((s) => s.unreadCount)

  return (
    <div className={styles.bell} onClick={onClick}>
      🔔
      {unreadCount > 0 && <span className={styles.badge}>{unreadCount > 99 ? '99+' : unreadCount}</span>}
    </div>
  )
}
```

- [ ] **Step 4: Create FilePreview component**

`web/src/components/common/FilePreview.module.css`:
```css
.file {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;
}

.icon {
  width: 40px;
  height: 40px;
  background: var(--accent);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  flex-shrink: 0;
}

.info {
  min-width: 0;
}

.name {
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.size {
  font-size: 11px;
  color: var(--text-secondary);
}
```

`web/src/components/common/FilePreview.tsx`:
```tsx
import type { Attachment } from '../../types'
import styles from './FilePreview.module.css'

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function getIcon(mimeType: string): string {
  if (mimeType.startsWith('image/')) return '🖼️'
  if (mimeType.startsWith('audio/')) return '🎵'
  if (mimeType.startsWith('video/')) return '🎬'
  return '📄'
}

interface Props {
  attachment: Attachment
  downloadUrl?: string
}

export function FilePreview({ attachment, downloadUrl }: Props) {
  const content = (
    <div className={styles.file}>
      <div className={styles.icon}>{getIcon(attachment.mime_type)}</div>
      <div className={styles.info}>
        <div className={styles.name}>{attachment.file_name}</div>
        <div className={styles.size}>{formatSize(attachment.file_size)}</div>
      </div>
    </div>
  )

  if (downloadUrl) {
    return <a href={downloadUrl} target="_blank" rel="noopener noreferrer" style={{ textDecoration: 'none', color: 'inherit' }}>{content}</a>
  }

  return content
}
```

- [ ] **Step 5: Verify TypeScript compiles**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha/web && npx tsc --noEmit
```

- [ ] **Step 6: Commit**

```bash
git add web/src/components/common/
git commit -m "feat: add common components — Avatar, OnlineIndicator, NotificationBell, FilePreview"
```

---

## Task 9: Sidebar Components

**Files:**
- Create: `web/src/components/sidebar/Sidebar.tsx` + `.module.css`
- Create: `web/src/components/sidebar/ChatItem.tsx` + `.module.css`
- Create: `web/src/components/sidebar/ChatList.tsx` + `.module.css`
- Create: `web/src/components/sidebar/SearchBar.tsx` + `.module.css`

- [ ] **Step 1: Create SearchBar**

`web/src/components/sidebar/SearchBar.module.css`:
```css
.container {
  padding: 8px 12px;
  border-bottom: 1px solid var(--border);
}

.row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.menu {
  width: 32px;
  height: 32px;
  background: none;
  border: none;
  color: var(--text-secondary);
  font-size: 18px;
  cursor: pointer;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.menu:hover {
  background: var(--bg-input);
}

.input {
  flex: 1;
  background: var(--bg-input);
  border: none;
  border-radius: 18px;
  padding: 8px 14px;
  font-size: 13px;
  color: var(--text-primary);
  outline: none;
}

.input::placeholder {
  color: var(--text-secondary);
}
```

`web/src/components/sidebar/SearchBar.tsx`:
```tsx
import { useUIStore } from '../../store/uiStore'
import styles from './SearchBar.module.css'

export function SearchBar() {
  const searchQuery = useUIStore((s) => s.searchQuery)
  const setSearchQuery = useUIStore((s) => s.setSearchQuery)

  return (
    <div className={styles.container}>
      <div className={styles.row}>
        <button className={styles.menu}>☰</button>
        <input
          className={styles.input}
          placeholder="Поиск"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>
    </div>
  )
}
```

- [ ] **Step 2: Create ChatItem**

`web/src/components/sidebar/ChatItem.module.css`:
```css
.item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  cursor: pointer;
}

.item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.active {
  background: var(--bg-active) !important;
}

.info {
  flex: 1;
  min-width: 0;
}

.top {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.name {
  font-size: 14px;
  font-weight: 500;
}

.time {
  font-size: 12px;
  color: var(--text-secondary);
  flex-shrink: 0;
}

.preview {
  font-size: 13px;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
```

`web/src/components/sidebar/ChatItem.tsx`:
```tsx
import type { Chat } from '../../types'
import { Avatar } from '../common/Avatar'
import styles from './ChatItem.module.css'

interface Props {
  chat: Chat
  isActive: boolean
  onClick: () => void
}

export function ChatItem({ chat, isActive, onClick }: Props) {
  return (
    <div
      className={`${styles.item} ${isActive ? styles.active : ''}`}
      onClick={onClick}
    >
      <Avatar name={chat.name || 'Chat'} size={48} />
      <div className={styles.info}>
        <div className={styles.top}>
          <span className={styles.name}>{chat.name || 'Direct'}</span>
        </div>
      </div>
    </div>
  )
}
```

- [ ] **Step 3: Create ChatList**

`web/src/components/sidebar/ChatList.module.css`:
```css
.list {
  flex: 1;
  overflow-y: auto;
}

.empty {
  padding: 24px;
  text-align: center;
  color: var(--text-secondary);
  font-size: 14px;
}
```

`web/src/components/sidebar/ChatList.tsx`:
```tsx
import { useChatStore } from '../../store/chatStore'
import { ChatItem } from './ChatItem'
import styles from './ChatList.module.css'

export function ChatList() {
  const chats = useChatStore((s) => s.chats)
  const activeChatId = useChatStore((s) => s.activeChatId)
  const setActiveChat = useChatStore((s) => s.setActiveChat)

  if (chats.length === 0) {
    return <div className={styles.empty}>Нет чатов</div>
  }

  return (
    <div className={styles.list}>
      {chats.map((chat) => (
        <ChatItem
          key={chat.id}
          chat={chat}
          isActive={chat.id === activeChatId}
          onClick={() => setActiveChat(chat.id)}
        />
      ))}
    </div>
  )
}
```

- [ ] **Step 4: Create Sidebar**

`web/src/components/sidebar/Sidebar.module.css`:
```css
.sidebar {
  width: 320px;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  height: 100%;
}

.tabs {
  display: flex;
  border-bottom: 1px solid var(--border);
}

.tab {
  flex: 1;
  padding: 10px;
  text-align: center;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  border: none;
  background: none;
  border-bottom: 2px solid transparent;
}

.tabActive {
  color: var(--accent);
  border-bottom-color: var(--accent);
}

.actions {
  padding: 8px 12px;
  border-top: 1px solid var(--border);
}

.newChatBtn {
  width: 100%;
  padding: 8px;
  background: var(--accent);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  cursor: pointer;
}

.newChatBtn:hover {
  opacity: 0.9;
}
```

`web/src/components/sidebar/Sidebar.tsx`:
```tsx
import { useUIStore } from '../../store/uiStore'
import { SearchBar } from './SearchBar'
import { ChatList } from './ChatList'
import styles from './Sidebar.module.css'

export function Sidebar() {
  const sidebarTab = useUIStore((s) => s.sidebarTab)
  const setSidebarTab = useUIStore((s) => s.setSidebarTab)
  const openModal = useUIStore((s) => s.openModal)

  return (
    <div className={styles.sidebar}>
      <SearchBar />
      <div className={styles.tabs}>
        <button
          className={`${styles.tab} ${sidebarTab === 'chats' ? styles.tabActive : ''}`}
          onClick={() => setSidebarTab('chats')}
        >
          Чаты
        </button>
        <button
          className={`${styles.tab} ${sidebarTab === 'channels' ? styles.tabActive : ''}`}
          onClick={() => setSidebarTab('channels')}
        >
          Каналы
        </button>
      </div>
      <ChatList />
      <div className={styles.actions}>
        <button className={styles.newChatBtn} onClick={() => openModal('createChat')}>
          + Новый чат
        </button>
      </div>
    </div>
  )
}
```

- [ ] **Step 5: Verify TypeScript compiles**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha/web && npx tsc --noEmit
```

- [ ] **Step 6: Commit**

```bash
git add web/src/components/sidebar/
git commit -m "feat: add sidebar components — SearchBar, ChatItem, ChatList, Sidebar"
```

---

## Task 10: Chat Components (MessageBubble, MessageList, MessageInput, ChatHeader, ChatView)

**Files:**
- Create all chat/ component files with their CSS modules

- [ ] **Step 1: Create TypingIndicator**

`web/src/components/chat/TypingIndicator.module.css`:
```css
.typing {
  font-size: 12px;
  color: var(--text-secondary);
  padding: 4px 16px;
  font-style: italic;
}
```

`web/src/components/chat/TypingIndicator.tsx`:
```tsx
import styles from './TypingIndicator.module.css'

interface Props {
  userIds: string[]
}

export function TypingIndicator({ userIds }: Props) {
  if (userIds.length === 0) return null

  const text = userIds.length === 1
    ? 'печатает...'
    : `${userIds.length} печатают...`

  return <div className={styles.typing}>{text}</div>
}
```

- [ ] **Step 2: Create ReactionBar**

`web/src/components/chat/ReactionBar.module.css`:
```css
.reactions {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  margin-top: 4px;
}

.reaction {
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 2px 6px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  font-size: 12px;
  cursor: pointer;
  border: none;
  color: var(--text-primary);
}

.reaction:hover {
  background: rgba(255, 255, 255, 0.15);
}

.mine {
  background: rgba(82, 136, 193, 0.3);
}

.count {
  font-size: 11px;
  color: var(--text-secondary);
}
```

`web/src/components/chat/ReactionBar.tsx`:
```tsx
import type { Reaction } from '../../types'
import styles from './ReactionBar.module.css'

interface Props {
  reactions: Reaction[]
  currentUserId: string
  onToggle: (emoji: string) => void
}

export function ReactionBar({ reactions, currentUserId, onToggle }: Props) {
  if (reactions.length === 0) return null

  // Group reactions by emoji
  const grouped = reactions.reduce<Record<string, { count: number; mine: boolean }>>((acc, r) => {
    if (!acc[r.emoji]) acc[r.emoji] = { count: 0, mine: false }
    acc[r.emoji].count++
    if (r.user_id === currentUserId) acc[r.emoji].mine = true
    return acc
  }, {})

  return (
    <div className={styles.reactions}>
      {Object.entries(grouped).map(([emoji, data]) => (
        <button
          key={emoji}
          className={`${styles.reaction} ${data.mine ? styles.mine : ''}`}
          onClick={() => onToggle(emoji)}
        >
          {emoji} <span className={styles.count}>{data.count}</span>
        </button>
      ))}
    </div>
  )
}
```

- [ ] **Step 3: Create MessageBubble**

`web/src/components/chat/MessageBubble.module.css`:
```css
.wrapper {
  display: flex;
  margin-bottom: 4px;
}

.incoming {
  justify-content: flex-start;
}

.outgoing {
  justify-content: flex-end;
}

.bubble {
  max-width: 65%;
  padding: 8px 12px;
  position: relative;
}

.bubbleIncoming {
  background: var(--bg-incoming);
  border-radius: 0 12px 12px 12px;
}

.bubbleOutgoing {
  background: var(--bg-active);
  border-radius: 12px 0 12px 12px;
}

.content {
  font-size: 13px;
  line-height: 1.4;
  word-break: break-word;
}

.meta {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 4px;
  margin-top: 2px;
}

.time {
  font-size: 11px;
  color: var(--text-secondary);
}

.contextMenu {
  position: absolute;
  top: -4px;
  right: -4px;
  opacity: 0;
  transition: opacity 0.15s;
  cursor: pointer;
  font-size: 14px;
  padding: 4px;
}

.wrapper:hover .contextMenu {
  opacity: 1;
}
```

`web/src/components/chat/MessageBubble.tsx`:
```tsx
import type { Message } from '../../types'
import styles from './MessageBubble.module.css'

interface Props {
  message: Message
  isOwn: boolean
  onReactionClick?: (messageId: string) => void
}

export function MessageBubble({ message, isOwn, onReactionClick }: Props) {
  const time = new Date(message.created_at).toLocaleTimeString('ru-RU', {
    hour: '2-digit',
    minute: '2-digit',
  })

  return (
    <div className={`${styles.wrapper} ${isOwn ? styles.outgoing : styles.incoming}`}>
      <div className={`${styles.bubble} ${isOwn ? styles.bubbleOutgoing : styles.bubbleIncoming}`}>
        <div className={styles.content}>{message.content}</div>
        <div className={styles.meta}>
          <span className={styles.time}>{time}</span>
        </div>
        <span
          className={styles.contextMenu}
          onClick={() => onReactionClick?.(message.id)}
        >
          😀
        </span>
      </div>
    </div>
  )
}
```

- [ ] **Step 4: Create MessageList**

`web/src/components/chat/MessageList.module.css`:
```css
.container {
  flex: 1;
  overflow-y: auto;
  padding: 12px 16px;
  display: flex;
  flex-direction: column-reverse;
  gap: 2px;
  background: var(--bg-primary);
}

.loading {
  text-align: center;
  padding: 12px;
  color: var(--text-secondary);
  font-size: 13px;
}

.empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--text-secondary);
  font-size: 14px;
}
```

`web/src/components/chat/MessageList.tsx`:
```tsx
import { useEffect } from 'react'
import { useMessageStore } from '../../store/messageStore'
import { useAuthStore } from '../../store/authStore'
import { useInfiniteScroll } from '../../hooks/useInfiniteScroll'
import { MessageBubble } from './MessageBubble'
import { TypingIndicator } from './TypingIndicator'
import styles from './MessageList.module.css'

interface Props {
  chatId: string
}

export function MessageList({ chatId }: Props) {
  const messages = useMessageStore((s) => s.messages[chatId] || [])
  const hasMore = useMessageStore((s) => s.hasMore[chatId] ?? true)
  const typing = useMessageStore((s) => s.typing[chatId] || [])
  const fetchHistory = useMessageStore((s) => s.fetchHistory)
  const user = useAuthStore((s) => s.user)

  useEffect(() => {
    if (messages.length === 0) {
      fetchHistory(chatId)
    }
  }, [chatId, messages.length, fetchHistory])

  const { containerRef, handleScroll } = useInfiniteScroll(
    () => fetchHistory(chatId),
    hasMore,
  )

  return (
    <>
      <div
        ref={containerRef}
        className={styles.container}
        onScroll={handleScroll}
      >
        {messages.length === 0 && !hasMore && (
          <div className={styles.empty}>Нет сообщений</div>
        )}
        {messages.map((msg) => (
          <MessageBubble
            key={msg.id}
            message={msg}
            isOwn={msg.sender_id === user?.id}
          />
        ))}
      </div>
      <TypingIndicator userIds={typing} />
    </>
  )
}
```

- [ ] **Step 5: Create MessageInput**

`web/src/components/chat/MessageInput.module.css`:
```css
.container {
  padding: 8px 12px;
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  gap: 8px;
}

.attach {
  background: none;
  border: none;
  color: var(--text-secondary);
  font-size: 20px;
  cursor: pointer;
  padding: 4px;
}

.attach:hover {
  color: var(--text-primary);
}

.input {
  flex: 1;
  background: var(--bg-input);
  border: none;
  border-radius: 18px;
  padding: 10px 16px;
  font-size: 13px;
  color: var(--text-primary);
  outline: none;
  resize: none;
}

.input::placeholder {
  color: var(--text-secondary);
}

.send {
  background: none;
  border: none;
  color: var(--accent);
  font-size: 20px;
  cursor: pointer;
  padding: 4px;
}

.send:disabled {
  color: var(--text-secondary);
  cursor: default;
}
```

`web/src/components/chat/MessageInput.tsx`:
```tsx
import { useState, useRef, type KeyboardEvent } from 'react'
import styles from './MessageInput.module.css'

interface Props {
  onSend: (content: string) => void
  onTyping: () => void
  onFileSelect: (file: File) => void
}

export function MessageInput({ onSend, onTyping, onFileSelect }: Props) {
  const [text, setText] = useState('')
  const fileRef = useRef<HTMLInputElement>(null)

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    } else {
      onTyping()
    }
  }

  function handleSend() {
    const trimmed = text.trim()
    if (!trimmed) return
    onSend(trimmed)
    setText('')
  }

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (file) onFileSelect(file)
    e.target.value = ''
  }

  return (
    <div className={styles.container}>
      <button className={styles.attach} onClick={() => fileRef.current?.click()}>
        📎
      </button>
      <input type="file" ref={fileRef} style={{ display: 'none' }} onChange={handleFileChange} />
      <input
        className={styles.input}
        placeholder="Сообщение"
        value={text}
        onChange={(e) => setText(e.target.value)}
        onKeyDown={handleKeyDown}
      />
      <button className={styles.send} onClick={handleSend} disabled={!text.trim()}>
        ➤
      </button>
    </div>
  )
}
```

- [ ] **Step 6: Create ChatHeader**

`web/src/components/chat/ChatHeader.module.css`:
```css
.header {
  padding: 8px 16px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  gap: 10px;
}

.info {
  flex: 1;
}

.name {
  font-size: 14px;
  font-weight: 500;
}

.status {
  font-size: 12px;
  color: var(--text-secondary);
}

.online {
  color: var(--online);
}
```

`web/src/components/chat/ChatHeader.tsx`:
```tsx
import type { Chat } from '../../types'
import { Avatar } from '../common/Avatar'
import styles from './ChatHeader.module.css'

interface Props {
  chat: Chat
}

export function ChatHeader({ chat }: Props) {
  return (
    <div className={styles.header}>
      <Avatar name={chat.name || 'Chat'} size={40} />
      <div className={styles.info}>
        <div className={styles.name}>{chat.name || 'Direct'}</div>
        <div className={styles.status}>{chat.type === 'group' ? 'группа' : ''}</div>
      </div>
    </div>
  )
}
```

- [ ] **Step 7: Create ChatView**

`web/src/components/chat/ChatView.module.css`:
```css
.container {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  font-size: 15px;
  background: var(--bg-primary);
}
```

`web/src/components/chat/ChatView.tsx`:
```tsx
import { useRef, useCallback } from 'react'
import { useChatStore } from '../../store/chatStore'
import { useAuthStore } from '../../store/authStore'
import { ChatHeader } from './ChatHeader'
import { MessageList } from './MessageList'
import { MessageInput } from './MessageInput'
import type { GrushaSocket } from '../../ws/socket'
import styles from './ChatView.module.css'

interface Props {
  socketRef: React.RefObject<GrushaSocket | null>
}

export function ChatView({ socketRef }: Props) {
  const activeChatId = useChatStore((s) => s.activeChatId)
  const chats = useChatStore((s) => s.chats)
  const accessToken = useAuthStore((s) => s.accessToken)
  const typingTimeout = useRef<ReturnType<typeof setTimeout>>()

  const chat = chats.find((c) => c.id === activeChatId)

  const handleSend = useCallback((content: string) => {
    if (!activeChatId) return
    socketRef.current?.send('send_message', {
      chat_id: activeChatId,
      content,
      msg_type: 'text',
    })
  }, [activeChatId, socketRef])

  const handleTyping = useCallback(() => {
    if (!activeChatId || typingTimeout.current) return
    socketRef.current?.send('typing', { chat_id: activeChatId })
    typingTimeout.current = setTimeout(() => {
      typingTimeout.current = undefined
    }, 2000)
  }, [activeChatId, socketRef])

  const handleFileSelect = useCallback((_file: File) => {
    // TODO: implement file upload flow
  }, [])

  if (!activeChatId || !chat) {
    return (
      <div className={styles.empty}>
        Выберите чат для начала общения
      </div>
    )
  }

  return (
    <div className={styles.container}>
      <ChatHeader chat={chat} />
      <MessageList chatId={activeChatId} />
      <MessageInput
        onSend={handleSend}
        onTyping={handleTyping}
        onFileSelect={handleFileSelect}
      />
    </div>
  )
}
```

- [ ] **Step 8: Verify TypeScript compiles**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha/web && npx tsc --noEmit
```

- [ ] **Step 9: Commit**

```bash
git add web/src/components/chat/
git commit -m "feat: add chat components — MessageBubble, MessageList, MessageInput, ChatHeader, ChatView"
```

---

## Task 11: Channel Components

**Files:**
- Create: `web/src/components/channel/ChannelView.tsx` + `.module.css`
- Create: `web/src/components/channel/ChannelList.tsx` + `.module.css`
- Create: `web/src/components/channel/ChannelHeader.tsx` + `.module.css`

- [ ] **Step 1: Create ChannelHeader**

`web/src/components/channel/ChannelHeader.module.css`:
```css
.header {
  padding: 8px 16px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  gap: 10px;
}

.info { flex: 1; }
.name { font-size: 14px; font-weight: 500; }
.desc { font-size: 12px; color: var(--text-secondary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.unsub {
  background: none;
  border: 1px solid var(--danger);
  color: var(--danger);
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
}
```

`web/src/components/channel/ChannelHeader.tsx`:
```tsx
import type { Channel } from '../../types'
import { Avatar } from '../common/Avatar'
import { useChannelStore } from '../../store/channelStore'
import styles from './ChannelHeader.module.css'

interface Props {
  channel: Channel
}

export function ChannelHeader({ channel }: Props) {
  const unsubscribe = useChannelStore((s) => s.unsubscribe)

  return (
    <div className={styles.header}>
      <Avatar name={channel.name} size={40} />
      <div className={styles.info}>
        <div className={styles.name}>{channel.name}</div>
        <div className={styles.desc}>{channel.description}</div>
      </div>
      <button className={styles.unsub} onClick={() => unsubscribe(channel.id)}>
        Отписаться
      </button>
    </div>
  )
}
```

- [ ] **Step 2: Create ChannelView (placeholder — channels use separate message flow)**

`web/src/components/channel/ChannelView.module.css`:
```css
.container { flex: 1; display: flex; flex-direction: column; height: 100%; }
.empty { flex: 1; display: flex; align-items: center; justify-content: center; color: var(--text-secondary); font-size: 15px; background: var(--bg-primary); }
```

`web/src/components/channel/ChannelView.tsx`:
```tsx
import { useChannelStore } from '../../store/channelStore'
import { ChannelHeader } from './ChannelHeader'
import styles from './ChannelView.module.css'

export function ChannelView() {
  const activeChannelId = useChannelStore((s) => s.activeChannelId)
  const channels = useChannelStore((s) => s.channels)
  const channel = channels.find((c) => c.id === activeChannelId)

  if (!activeChannelId || !channel) {
    return <div className={styles.empty}>Выберите канал</div>
  }

  return (
    <div className={styles.container}>
      <ChannelHeader channel={channel} />
      <div className={styles.empty}>Сообщения канала</div>
    </div>
  )
}
```

- [ ] **Step 3: Create ChannelList**

`web/src/components/channel/ChannelList.module.css`:
```css
.list { flex: 1; overflow-y: auto; }
.item { display: flex; align-items: center; gap: 10px; padding: 8px 12px; cursor: pointer; }
.item:hover { background: rgba(255,255,255,0.05); }
.active { background: var(--bg-active) !important; }
.info { flex: 1; min-width: 0; }
.name { font-size: 14px; font-weight: 500; }
.desc { font-size: 12px; color: var(--text-secondary); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.empty { padding: 24px; text-align: center; color: var(--text-secondary); font-size: 14px; }
```

`web/src/components/channel/ChannelList.tsx`:
```tsx
import { useChannelStore } from '../../store/channelStore'
import { Avatar } from '../common/Avatar'
import styles from './ChannelList.module.css'

export function ChannelList() {
  const channels = useChannelStore((s) => s.channels)
  const activeChannelId = useChannelStore((s) => s.activeChannelId)
  const setActiveChannel = useChannelStore((s) => s.setActiveChannel)

  if (channels.length === 0) {
    return <div className={styles.empty}>Нет подписок на каналы</div>
  }

  return (
    <div className={styles.list}>
      {channels.map((ch) => (
        <div
          key={ch.id}
          className={`${styles.item} ${ch.id === activeChannelId ? styles.active : ''}`}
          onClick={() => setActiveChannel(ch.id)}
        >
          <Avatar name={ch.name} size={48} />
          <div className={styles.info}>
            <div className={styles.name}>{ch.name}</div>
            <div className={styles.desc}>{ch.description}</div>
          </div>
        </div>
      ))}
    </div>
  )
}
```

- [ ] **Step 4: Verify TypeScript compiles**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha/web && npx tsc --noEmit
```

- [ ] **Step 5: Commit**

```bash
git add web/src/components/channel/
git commit -m "feat: add channel components — ChannelHeader, ChannelView, ChannelList"
```

---

## Task 12: Modal Components

**Files:**
- Create: `web/src/components/modals/CreateChatModal.tsx` + `.module.css`
- Create: `web/src/components/modals/CreateGroupModal.tsx` + `.module.css`
- Create: `web/src/components/modals/CreateChannelModal.tsx` + `.module.css`
- Create: `web/src/components/modals/UserProfileModal.tsx` + `.module.css`

- [ ] **Step 1: Create shared modal styles**

`web/src/components/modals/CreateChatModal.module.css`:
```css
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.modal {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 24px;
  width: 380px;
  max-height: 80vh;
  overflow-y: auto;
}

.title {
  font-size: 18px;
  font-weight: 500;
  margin-bottom: 16px;
}

.input {
  width: 100%;
  padding: 10px 14px;
  margin-bottom: 12px;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
}

.input:focus {
  border-color: var(--accent);
}

.buttons {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
}

.btnPrimary {
  padding: 8px 20px;
  background: var(--accent);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
}

.btnSecondary {
  padding: 8px 20px;
  background: none;
  color: var(--text-secondary);
  border: 1px solid var(--border);
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
}

.error {
  color: var(--danger);
  font-size: 13px;
  margin-bottom: 8px;
}
```

- [ ] **Step 2: Create CreateChatModal**

`web/src/components/modals/CreateChatModal.tsx`:
```tsx
import { useState } from 'react'
import { useChatStore } from '../../store/chatStore'
import { useUIStore } from '../../store/uiStore'
import styles from './CreateChatModal.module.css'

export function CreateChatModal() {
  const [recipientId, setRecipientId] = useState('')
  const [error, setError] = useState('')
  const createDirect = useChatStore((s) => s.createDirect)
  const closeModal = useUIStore((s) => s.closeModal)

  async function handleCreate() {
    if (!recipientId.trim()) return
    setError('')
    try {
      await createDirect(recipientId.trim())
      closeModal()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error')
    }
  }

  return (
    <div className={styles.overlay} onClick={closeModal}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <h2 className={styles.title}>Новый чат</h2>
        {error && <div className={styles.error}>{error}</div>}
        <input
          className={styles.input}
          placeholder="ID пользователя"
          value={recipientId}
          onChange={(e) => setRecipientId(e.target.value)}
        />
        <div className={styles.buttons}>
          <button className={styles.btnSecondary} onClick={closeModal}>Отмена</button>
          <button className={styles.btnPrimary} onClick={handleCreate}>Создать</button>
        </div>
      </div>
    </div>
  )
}
```

- [ ] **Step 3: Create CreateGroupModal**

Copy `CreateChatModal.module.css` to `CreateGroupModal.module.css`.

`web/src/components/modals/CreateGroupModal.tsx`:
```tsx
import { useState } from 'react'
import { useChatStore } from '../../store/chatStore'
import { useUIStore } from '../../store/uiStore'
import styles from './CreateGroupModal.module.css'

export function CreateGroupModal() {
  const [name, setName] = useState('')
  const [memberIds, setMemberIds] = useState('')
  const [error, setError] = useState('')
  const createGroup = useChatStore((s) => s.createGroup)
  const closeModal = useUIStore((s) => s.closeModal)

  async function handleCreate() {
    if (!name.trim()) return
    setError('')
    try {
      const ids = memberIds.split(',').map((s) => s.trim()).filter(Boolean)
      await createGroup(name.trim(), ids)
      closeModal()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error')
    }
  }

  return (
    <div className={styles.overlay} onClick={closeModal}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <h2 className={styles.title}>Новая группа</h2>
        {error && <div className={styles.error}>{error}</div>}
        <input className={styles.input} placeholder="Название группы" value={name} onChange={(e) => setName(e.target.value)} />
        <input className={styles.input} placeholder="ID участников (через запятую)" value={memberIds} onChange={(e) => setMemberIds(e.target.value)} />
        <div className={styles.buttons}>
          <button className={styles.btnSecondary} onClick={closeModal}>Отмена</button>
          <button className={styles.btnPrimary} onClick={handleCreate}>Создать</button>
        </div>
      </div>
    </div>
  )
}
```

- [ ] **Step 4: Create CreateChannelModal**

Copy `CreateChatModal.module.css` to `CreateChannelModal.module.css`.

`web/src/components/modals/CreateChannelModal.tsx`:
```tsx
import { useState } from 'react'
import { channelsApi } from '../../api/channels'
import { useChannelStore } from '../../store/channelStore'
import { useUIStore } from '../../store/uiStore'
import styles from './CreateChannelModal.module.css'

export function CreateChannelModal() {
  const [slug, setSlug] = useState('')
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [error, setError] = useState('')
  const closeModal = useUIStore((s) => s.closeModal)
  const fetchMyChannels = useChannelStore((s) => s.fetchMyChannels)

  async function handleCreate() {
    if (!slug.trim() || !name.trim()) return
    setError('')
    try {
      await channelsApi.create(slug.trim(), name.trim(), description.trim(), false)
      await fetchMyChannels()
      closeModal()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error')
    }
  }

  return (
    <div className={styles.overlay} onClick={closeModal}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <h2 className={styles.title}>Новый канал</h2>
        {error && <div className={styles.error}>{error}</div>}
        <input className={styles.input} placeholder="Slug (@channel_name)" value={slug} onChange={(e) => setSlug(e.target.value)} />
        <input className={styles.input} placeholder="Название" value={name} onChange={(e) => setName(e.target.value)} />
        <input className={styles.input} placeholder="Описание" value={description} onChange={(e) => setDescription(e.target.value)} />
        <div className={styles.buttons}>
          <button className={styles.btnSecondary} onClick={closeModal}>Отмена</button>
          <button className={styles.btnPrimary} onClick={handleCreate}>Создать</button>
        </div>
      </div>
    </div>
  )
}
```

- [ ] **Step 5: Create UserProfileModal**

`web/src/components/modals/UserProfileModal.module.css` — copy from CreateChatModal.module.css.

`web/src/components/modals/UserProfileModal.tsx`:
```tsx
import { useAuthStore } from '../../store/authStore'
import { useUIStore } from '../../store/uiStore'
import { Avatar } from '../common/Avatar'
import styles from './UserProfileModal.module.css'

export function UserProfileModal() {
  const user = useAuthStore((s) => s.user)
  const logout = useAuthStore((s) => s.logout)
  const closeModal = useUIStore((s) => s.closeModal)

  return (
    <div className={styles.overlay} onClick={closeModal}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <h2 className={styles.title}>Профиль</h2>
        {user ? (
          <div style={{ textAlign: 'center' }}>
            <Avatar name={user.display_name || user.username} size={80} />
            <div style={{ marginTop: 12, fontSize: 18, fontWeight: 500 }}>{user.display_name}</div>
            <div style={{ color: 'var(--text-secondary)', fontSize: 14 }}>@{user.username}</div>
            <div style={{ color: 'var(--text-secondary)', fontSize: 13, marginTop: 4 }}>{user.email}</div>
            {user.bio && <div style={{ marginTop: 12, fontSize: 13 }}>{user.bio}</div>}
          </div>
        ) : (
          <div style={{ color: 'var(--text-secondary)' }}>Не авторизован</div>
        )}
        <div className={styles.buttons}>
          <button className={styles.btnSecondary} onClick={closeModal}>Закрыть</button>
          <button className={styles.btnPrimary} style={{ background: 'var(--danger)' }} onClick={() => { logout(); closeModal() }}>
            Выйти
          </button>
        </div>
      </div>
    </div>
  )
}
```

- [ ] **Step 6: Verify TypeScript compiles**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha/web && npx tsc --noEmit
```

- [ ] **Step 7: Commit**

```bash
git add web/src/components/modals/
git commit -m "feat: add modal components — CreateChat, CreateGroup, CreateChannel, UserProfile"
```

---

## Task 13: MainPage — Wire Everything Together

**Files:**
- Create: `web/src/pages/MainPage.tsx`, `web/src/pages/MainPage.module.css`
- Modify: `web/src/App.tsx` — replace placeholder, wire WebSocket
- Modify: `web/src/components/sidebar/Sidebar.tsx` — support channels tab

- [ ] **Step 1: Create MainPage styles**

`web/src/pages/MainPage.module.css`:
```css
.container {
  display: flex;
  height: 100vh;
}
```

- [ ] **Step 2: Create MainPage**

`web/src/pages/MainPage.tsx`:
```tsx
import { useEffect } from 'react'
import { useChatStore } from '../store/chatStore'
import { useChannelStore } from '../store/channelStore'
import { useNotificationStore } from '../store/notificationStore'
import { useAuthStore } from '../store/authStore'
import { useUIStore } from '../store/uiStore'
import { useWebSocket } from '../hooks/useWebSocket'
import { usersApi } from '../api/users'
import { Sidebar } from '../components/sidebar/Sidebar'
import { ChatView } from '../components/chat/ChatView'
import { ChannelView } from '../components/channel/ChannelView'
import { CreateChatModal } from '../components/modals/CreateChatModal'
import { CreateGroupModal } from '../components/modals/CreateGroupModal'
import { CreateChannelModal } from '../components/modals/CreateChannelModal'
import { UserProfileModal } from '../components/modals/UserProfileModal'
import styles from './MainPage.module.css'

export function MainPage() {
  const socketRef = useWebSocket()
  const fetchChats = useChatStore((s) => s.fetchChats)
  const fetchMyChannels = useChannelStore((s) => s.fetchMyChannels)
  const fetchUnreadCount = useNotificationStore((s) => s.fetchUnreadCount)
  const sidebarTab = useUIStore((s) => s.sidebarTab)
  const activeModal = useUIStore((s) => s.activeModal)
  const accessToken = useAuthStore((s) => s.accessToken)

  // Fetch user profile on mount
  useEffect(() => {
    const state = useAuthStore.getState()
    if (state.accessToken && !state.user) {
      // We need userId from JWT — decode it
      try {
        const payload = JSON.parse(atob(state.accessToken.split('.')[1]))
        if (payload.user_id) {
          usersApi.getProfile(payload.user_id).then((res) => {
            useAuthStore.setState({ user: res.user })
          })
        }
      } catch { /* ignore */ }
    }
  }, [accessToken])

  useEffect(() => {
    fetchChats()
    fetchMyChannels()
    fetchUnreadCount()
  }, [fetchChats, fetchMyChannels, fetchUnreadCount])

  return (
    <div className={styles.container}>
      <Sidebar />
      {sidebarTab === 'chats' ? (
        <ChatView socketRef={socketRef} />
      ) : (
        <ChannelView />
      )}

      {activeModal === 'createChat' && <CreateChatModal />}
      {activeModal === 'createGroup' && <CreateGroupModal />}
      {activeModal === 'createChannel' && <CreateChannelModal />}
      {activeModal === 'profile' && <UserProfileModal />}
    </div>
  )
}
```

- [ ] **Step 3: Update App.tsx — final version**

```tsx
import { useEffect } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import { LoginPage } from './pages/LoginPage'
import { RegisterPage } from './pages/RegisterPage'
import { MainPage } from './pages/MainPage'

function AuthGuard({ children }: { children: React.ReactNode }) {
  const accessToken = useAuthStore((s) => s.accessToken)
  if (!accessToken) return <Navigate to="/login" replace />
  return <>{children}</>
}

export function App() {
  useEffect(() => {
    useAuthStore.getState().restoreSession()
  }, [])

  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/" element={<AuthGuard><MainPage /></AuthGuard>} />
    </Routes>
  )
}
```

- [ ] **Step 4: Update Sidebar to support channels tab and new-channel button**

Update `web/src/components/sidebar/Sidebar.tsx` to conditionally render ChatList or ChannelList based on sidebarTab:

```tsx
import { useUIStore } from '../../store/uiStore'
import { SearchBar } from './SearchBar'
import { ChatList } from './ChatList'
import { ChannelList } from '../../components/channel/ChannelList'
import { NotificationBell } from '../common/NotificationBell'
import styles from './Sidebar.module.css'

export function Sidebar() {
  const sidebarTab = useUIStore((s) => s.sidebarTab)
  const setSidebarTab = useUIStore((s) => s.setSidebarTab)
  const openModal = useUIStore((s) => s.openModal)

  return (
    <div className={styles.sidebar}>
      <SearchBar />
      <div className={styles.tabs}>
        <button
          className={`${styles.tab} ${sidebarTab === 'chats' ? styles.tabActive : ''}`}
          onClick={() => setSidebarTab('chats')}
        >
          Чаты
        </button>
        <button
          className={`${styles.tab} ${sidebarTab === 'channels' ? styles.tabActive : ''}`}
          onClick={() => setSidebarTab('channels')}
        >
          Каналы
        </button>
        <NotificationBell onClick={() => openModal('profile')} />
      </div>
      {sidebarTab === 'chats' ? <ChatList /> : <ChannelList />}
      <div className={styles.actions}>
        <button
          className={styles.newChatBtn}
          onClick={() => openModal(sidebarTab === 'chats' ? 'createChat' : 'createChannel')}
        >
          + {sidebarTab === 'chats' ? 'Новый чат' : 'Новый канал'}
        </button>
      </div>
    </div>
  )
}
```

- [ ] **Step 5: Verify TypeScript compiles**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha/web && npx tsc --noEmit
```

- [ ] **Step 6: Verify dev server starts and shows login page**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha/web
npx vite --host 127.0.0.1 &
sleep 3
curl -s http://127.0.0.1:5173 | grep -o 'Груша'
kill %1
```

Expected: `Груша`

- [ ] **Step 7: Commit**

```bash
git add web/src/pages/MainPage.tsx web/src/pages/MainPage.module.css web/src/App.tsx web/src/components/sidebar/Sidebar.tsx
git commit -m "feat: wire MainPage with sidebar, chat view, channels, modals, and WebSocket"
```

---

## Task 14: Add .gitignore + Makefile Targets for Frontend

**Files:**
- Modify: `.gitignore` — add web/node_modules, web/dist, .superpowers
- Modify: `Makefile` — add web-install, web-dev, web-build targets

- [ ] **Step 1: Update .gitignore**

Append to `.gitignore`:

```
# Frontend
web/node_modules/
web/dist/

# Superpowers brainstorm sessions
.superpowers/
```

- [ ] **Step 2: Add Makefile targets**

Append to `Makefile`:

```makefile
# Frontend
web-install:
	cd web && npm install

web-dev:
	cd web && npm run dev

web-build:
	cd web && npm run build
```

- [ ] **Step 3: Commit**

```bash
git add .gitignore Makefile
git commit -m "chore: add frontend gitignore entries and Makefile targets"
```

---

## Task 15: End-to-End Smoke Test

- [ ] **Step 1: Start infrastructure**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha
make docker-up
```

Expected: Postgres, Redis, MinIO containers running.

- [ ] **Step 2: Run migrations**

Run:
```bash
DATABASE_URL="postgres://grusha:grusha_secret@localhost:5432/grusha?sslmode=disable" make migrate-up
```

- [ ] **Step 3: Start Go backend**

Run:
```bash
make run &
```

Wait for "gRPC server starting" and "HTTP/WS server starting" log messages.

- [ ] **Step 4: Verify grpc-gateway responds**

Run:
```bash
curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test1","email":"test1@example.com","password":"password123","display_name":"Test User 1"}'
```

Expected: JSON response with `user` object.

- [ ] **Step 5: Verify login**

Run:
```bash
curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test1@example.com","password":"password123"}'
```

Expected: JSON with `access_token` and `refresh_token`.

- [ ] **Step 6: Start frontend and verify it loads**

Run:
```bash
cd /Users/tec-5/GolandProjects/пэтпроект/MessngerGrusha/web
npm run dev &
sleep 3
curl -s http://localhost:5173 | grep -o 'Груша'
```

Expected: `Груша` found in page HTML.

- [ ] **Step 7: Stop all processes**

```bash
kill %1 %2 2>/dev/null
make docker-down
```

- [ ] **Step 8: Final commit if any fixes were needed**

```bash
git add -A
git status
# Only commit if there are changes
git diff --staged --quiet || git commit -m "fix: smoke test fixes"
```
