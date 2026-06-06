const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? 'http://api.localhost'

export class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

async function request(path: string, options?: RequestInit): Promise<Response> {
  const url = `${API_BASE}${path}`
  const res = await fetch(url, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...options?.headers },
    ...options,
  })

  if (res.status !== 401) return res

  const refreshRes = await fetch(`${API_BASE}/auth/refresh`, {
    method: 'POST',
    credentials: 'include',
  })

  if (!refreshRes.ok) {
    throw new ApiError('Unauthorized', 401)
  }

  const retryRes = await fetch(url, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...options?.headers },
    ...options,
  })

  return retryRes
}

async function handleResponse<T>(res: Response): Promise<T> {
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))

    throw new ApiError(body?.error ?? res.statusText, res.status)
  }
  if (res.status === 204) return undefined as T

  return res.json()
}

export const api = {
  get<T>(path: string): Promise<T> {
    return request(path).then(handleResponse<T>)
  },

  post<T>(path: string, body?: unknown): Promise<T> {
    return request(path, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    }).then(handleResponse<T>)
  },

  put<T>(path: string, body?: unknown): Promise<T> {
    return request(path, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    }).then(handleResponse<T>)
  },

  raw(path: string, options?: RequestInit): Promise<Response> {
    return request(path, options)
  },
}
