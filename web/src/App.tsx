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
