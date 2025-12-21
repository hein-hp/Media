<template>
  <RecycleScroller
    class="h-[calc(100vh-8rem)] p-6"
    :items="rows"
    :item-size="itemHeight"
    key-field="rowIndex"
    v-slot="{ item: row }"
  >
    <div class="grid gap-4" :style="gridStyle">
      <div
        v-for="media in row.items"
        :key="media.path"
        class="group relative aspect-square bg-gray-50 rounded-xl overflow-hidden cursor-pointer
               shadow-sm hover:shadow-lg transition-all duration-300 ease-out
               hover:scale-[1.02] border border-gray-100 hover:border-gray-200"
        @click="$emit('select', media.originalIndex)"
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
  </RecycleScroller>
</template>

<script lang="ts" setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { RecycleScroller } from 'vue-virtual-scroller'
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'
import type { MediaInfo } from '@/types'

const props = defineProps<{
  images: MediaInfo[]
}>()

defineEmits<{
  select: [index: number]
}>()

// 响应式列数
const columnsCount = ref(6)
const containerWidth = ref(0)

// 计算每个卡片的高度（包含 gap）
const GAP = 16 // gap-4 = 16px
const PADDING = 24 // p-6 = 24px

const itemHeight = computed(() => {
  // 计算单个卡片宽度
  const availableWidth = containerWidth.value - PADDING * 2
  const cardWidth = (availableWidth - GAP * (columnsCount.value - 1)) / columnsCount.value
  // 卡片是正方形 (aspect-square)，加上 gap
  return cardWidth + GAP
})

// grid 样式
const gridStyle = computed(() => ({
  gridTemplateColumns: `repeat(${columnsCount.value}, minmax(0, 1fr))`
}))

// 将媒体数组按行分组
const rows = computed(() => {
  const result: { rowIndex: number; items: (MediaInfo & { originalIndex: number })[] }[] = []
  
  for (let i = 0; i < props.images.length; i += columnsCount.value) {
    const rowItems = props.images.slice(i, i + columnsCount.value).map((item, idx) => ({
      ...item,
      originalIndex: i + idx
    }))
    result.push({
      rowIndex: Math.floor(i / columnsCount.value),
      items: rowItems
    })
  }
  
  return result
})

// 根据窗口宽度计算列数
function updateColumns() {
  const width = window.innerWidth
  containerWidth.value = width
  
  if (width >= 1280) {      // xl
    columnsCount.value = 6
  } else if (width >= 1024) { // lg
    columnsCount.value = 5
  } else if (width >= 768) {  // md
    columnsCount.value = 4
  } else if (width >= 640) {  // sm
    columnsCount.value = 3
  } else {
    columnsCount.value = 2
  }
}

onMounted(() => {
  updateColumns()
  window.addEventListener('resize', updateColumns)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateColumns)
})

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
