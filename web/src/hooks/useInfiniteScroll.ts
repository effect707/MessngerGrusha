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

    if (el.scrollTop < 100) {
      loading.current = true
      const prevHeight = el.scrollHeight
      await onLoadMore()

      const newHeight = el.scrollHeight
      el.scrollTop = newHeight - prevHeight
      loading.current = false
    }
  }, [onLoadMore, hasMore])

  return { containerRef, handleScroll }
}
