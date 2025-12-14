import { ref, onMounted, onUnmounted } from 'vue'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import type { MediaInfo } from '@/types'

/**
 * 媒体列表管理 composable 监听 media-list 事件，管理媒体列表状态
 */
export function useMediaList() {
  const mediaList = ref<MediaInfo[]>([])
  const isLoading = ref(false)

  const handleMediaList = (data: MediaInfo[]) => {
    mediaList.value = data
    isLoading.value = false
  }

  onMounted(() => {
    EventsOn('media-list', handleMediaList)
  })

  onUnmounted(() => {
    EventsOff('media-list')
  })

  return {
    mediaList,
    isLoading
  }
}
