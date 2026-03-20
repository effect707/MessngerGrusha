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
