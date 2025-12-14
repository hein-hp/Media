<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="isOpen && media"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/90 backdrop-blur-sm">
        <!-- 顶部工具栏 -->
        <div class="absolute top-0 inset-x-0 z-10 h-16 flex items-center justify-between px-6
                    bg-gradient-to-b from-black/50 to-transparent">
          <div class="flex items-center gap-4">
            <!-- 媒体类型图标 -->
            <span v-if="isVideo(media)" class="text-white/70">
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                <path
                  d="M18 4l2 4h-3l-2-4h-2l2 4h-3l-2-4H8l2 4H7L5 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V4h-4z" />
              </svg>
            </span>
            <span class="text-white/90 text-sm font-medium">{{ media.name }}</span>
            <span class="text-white/50 text-xs">{{ formatSize(media.size) }}</span>
          </div>
          <div class="flex items-center gap-2">
            <!-- 媒体缩放控制 -->
            <template v-if="isImage(media)">
              <button class="p-2 rounded-lg text-white/70 hover:text-white hover:bg-white/10 transition-colors"
                title="缩小" @click="$emit('zoom-out')">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
                </svg>
              </button>
              <span class="text-white/70 text-sm min-w-[4rem] text-center">{{ Math.round(scale * 100) }}%</span>
              <button class="p-2 rounded-lg text-white/70 hover:text-white hover:bg-white/10 transition-colors"
                title="放大" @click="$emit('zoom-in')">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
              </button>
              <button class="p-2 rounded-lg text-white/70 hover:text-white hover:bg-white/10 transition-colors"
                title="重置" @click="$emit('reset-zoom')">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4" />
                </svg>
              </button>
            </template>
            <!-- 关闭按钮 -->
            <button class="ml-4 p-2 rounded-lg text-white/70 hover:text-white hover:bg-white/10 transition-colors"
              title="关闭" @click="$emit('close')">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- 主内容区域 -->
        <div class="relative flex items-center justify-center w-full h-full px-20 py-20 overflow-hidden"
          @mousemove="handleMouseMove" @mouseup="handleMouseUp" @mouseleave="handleMouseUp"
          @wheel.prevent="handleWheel">
          <Transition name="slide" mode="out-in">
            <!-- 媒体预览 -->
            <img v-if="isImage(media)" :key="`img-${media.url}`" :src="media.url" :alt="media.name"
              class="max-w-full max-h-full object-contain select-none" :class="[
                isDragging ? 'cursor-grabbing' : (scale > 1 ? 'cursor-grab' : 'cursor-default'),
                isDragging ? '' : 'transition-transform duration-200'
              ]" :style="{
                transform: `translate(${offsetX}px, ${offsetY}px) scale(${scale})`
              }" draggable="false" @mousedown="handleMouseDown" />

            <!-- 视频预览 -->
            <video v-else-if="isVideo(media)" :key="`video-${media.url}`" :src="media.url"
              class="max-w-full max-h-full object-contain" controls autoplay />
          </Transition>
        </div>

        <!-- 左侧导航按钮 -->
        <button v-if="hasPrev" class="absolute left-4 top-1/2 z-10 -translate-y-1/2 p-3 rounded-full
                 bg-white/10 hover:bg-white/20 text-white/80 hover:text-white
                 transition-all duration-200 backdrop-blur-sm" @click="$emit('prev')">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>

        <!-- 右侧导航按钮 -->
        <button v-if="hasNext" class="absolute right-4 top-1/2 z-10 -translate-y-1/2 p-3 rounded-full
                 bg-white/10 hover:bg-white/20 text-white/80 hover:text-white
                 transition-all duration-200 backdrop-blur-sm" @click="$emit('next')">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>

        <!-- 底部页码指示器 -->
        <div class="absolute bottom-6 left-1/2 z-10 -translate-x-1/2 px-4 py-2 rounded-full
                    bg-black/40 backdrop-blur-sm">
          <span class="text-white/90 text-sm font-medium">
            {{ currentIndex + 1 }} / {{ total }}
          </span>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script lang="ts" setup>
import { watch } from 'vue'
import { isImage, isVideo, type MediaInfo } from '@/types'

const props = defineProps<{
  isOpen: boolean
  media: MediaInfo | null
  currentIndex: number
  total: number
  scale: number
  offsetX: number
  offsetY: number
  isDragging: boolean
  hasNext: boolean
  hasPrev: boolean
}>()

const emit = defineEmits<{
  'close': []
  'prev': []
  'next': []
  'zoom-in': []
  'zoom-out': []
  'reset-zoom': []
  'start-drag': []
  'end-drag': []
  'drag': [deltaX: number, deltaY: number]
}>()

// 拖拽相关
let lastX = 0
let lastY = 0

function handleMouseDown(e: MouseEvent) {
  if (props.scale > 1) {
    e.preventDefault()
    lastX = e.clientX
    lastY = e.clientY
    emit('start-drag')
  }
}

function handleMouseMove(e: MouseEvent) {
  if (props.isDragging) {
    const deltaX = e.clientX - lastX
    const deltaY = e.clientY - lastY
    lastX = e.clientX
    lastY = e.clientY
    emit('drag', deltaX, deltaY)
  }
}

function handleMouseUp() {
  if (props.isDragging) {
    emit('end-drag')
  }
}

// 滚轮节流控制
let lastWheelTime = 0
const WHEEL_THROTTLE_MS = 100 // 节流间隔

function handleWheel(e: WheelEvent) {
  if (!props.media || !isImage(props.media)) {
    return
  }
  const now = Date.now()
  if (now - lastWheelTime < WHEEL_THROTTLE_MS) {
    return
  }
  lastWheelTime = now
  if (e.deltaY < 0) {
    emit('zoom-in')
  } else {
    emit('zoom-out')
  }
}

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

// 打开时禁止 body 滚动
watch(() => props.isOpen, (isOpen) => {
  document.body.style.overflow = isOpen ? 'hidden' : ''
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.slide-enter-from {
  opacity: 0;
  transform: scale(0.95);
}

.slide-leave-to {
  opacity: 0;
  transform: scale(1.05);
}

/* 视频播放器样式 */
video::-webkit-media-controls {
  background: transparent;
}
</style>
