import { useState, useCallback } from 'react'

interface UseApiState<T> {
  data: T | null
  error: string | null
  loading: boolean
}

/**
 * Простой хук для вызовов API с состоянием loading / error / data.
 */
export function useApi<T>() {
  const [state, setState] = useState<UseApiState<T>>({
    data: null,
    error: null,
    loading: false,
  })

  const execute = useCallback(async (fn: () => Promise<T>) => {
    setState({ data: null, error: null, loading: true })
    try {
      const data = await fn()
      setState({ data, error: null, loading: false })
      return data
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Неизвестная ошибка'
      setState({ data: null, error: message, loading: false })
      return null
    }
  }, [])

  const reset = useCallback(() => {
    setState({ data: null, error: null, loading: false })
  }, [])

  return { ...state, execute, reset }
}
