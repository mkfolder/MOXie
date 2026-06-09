const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? 'https://moxie-bhwong.fly.dev'

export class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

function buildRequestInit(options?: RequestInit): RequestInit {
  const is_form_data = options?.body instanceof FormData

  return {
    credentials: 'include',
    ...options,
    headers: is_form_data
      ? (options?.headers as Record<string, string>)
      : { 'Content-Type': 'application/json', ...(options?.headers as Record<string, string>) },
  }
}

async function request(path: string, options?: RequestInit): Promise<Response> {
  const url = `${API_BASE}${path}`
  let init = buildRequestInit(options)
  let res = await fetch(url, init)

  if (res.status !== 401) return res

  const refreshRes = await fetch(`${API_BASE}/auth/refresh`, {
    method: 'POST',
    credentials: 'include',
  })

  if (!refreshRes.ok) {
    throw new ApiError('Unauthorized', 401)
  }

  init = buildRequestInit(options)
  res = await fetch(url, init)

  return res
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

  upload<T>(path: string, formData: FormData): Promise<T> {
    return request(path, {
      method: 'PUT',
      body: formData,
    }).then(handleResponse<T>)
  },

  raw(path: string, options?: RequestInit): Promise<Response> {
    return request(path, options)
  },
}
