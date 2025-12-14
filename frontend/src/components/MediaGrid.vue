<template>
  <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4 p-6">
    <div
      v-for="(media, index) in images"
      :key="media.path"
      class="group relative aspect-square bg-gray-50 rounded-xl overflow-hidden cursor-pointer
             shadow-sm hover:shadow-lg transition-all duration-300 ease-out
             hover:scale-[1.02] border border-gray-100 hover:border-gray-200"
      @click="$emit('select', index)"
    >
      <!-- 图片 -->
      <img
        v-if="media.type === 'image'"
        :src="media.url"
        :alt="media.name"
        class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
        loading="lazy"
      />

      <!-- 视频缩略图 -->
      <div v-else-if="media.type === 'video'" class="relative w-full h-full bg-gray-100">
        <video
          :src="media.url"
          class="w-full h-full object-cover"
          preload="metadata"
          muted
        ></video>
        <!-- 视频播放图标 -->
        <div class="absolute inset-0 flex items-center justify-center">
          <div class="w-12 h-12 rounded-full bg-black/50 flex items-center justify-center
                      group-hover:bg-black/70 transition-colors">
            <svg class="w-6 h-6 text-white ml-1" fill="currentColor" viewBox="0 0 24 24">
              <path d="M8 5v14l11-7z" />
            </svg>
          </div>
        </div>
      </div>

      <!-- 悬浮信息层 -->
      <div class="absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/60 via-black/20 to-transparent
                  opacity-0 group-hover:opacity-100 transition-opacity duration-300 p-3 pt-8">
        <div class="flex items-center gap-2">
          <!-- 媒体类型图标 -->
          <span v-if="media.type === 'video'" class="text-white/80">
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
              <path d="M18 4l2 4h-3l-2-4h-2l2 4h-3l-2-4H8l2 4H7L5 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V4h-4z"/>
            </svg>
          </span>
          <p class="text-white text-sm font-medium truncate flex-1">{{ media.name }}</p>
        </div>
        <p class="text-white/70 text-xs">{{ formatSize(media.size) }}</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import type { MediaInfo } from '@/types'

defineProps<{
  images: MediaInfo[]
}>()

defineEmits<{
  select: [index: number]
}>()

/**
 * 格式化文件大小
 */
function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}
</script>
