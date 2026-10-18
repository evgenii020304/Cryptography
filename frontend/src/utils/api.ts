const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

async function parseError(res: Response): Promise<string> {
  const text = await res.text()
  try {
    const j = JSON.parse(text) as { error?: string }
    if (j.error) return j.error
  } catch {
    /* not json */
  }
  return text || `HTTP ${res.status}`
}

export async function apiPost<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    throw new Error(await parseError(res))
  }

  return res.json() as Promise<T>
}

export async function apiGet<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`)

  if (!res.ok) {
    throw new Error(await parseError(res))
  }

  return res.json() as Promise<T>
}