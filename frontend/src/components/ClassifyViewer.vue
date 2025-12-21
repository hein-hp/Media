<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="isOpen && media"
        class="fixed inset-0 z-50 flex flex-col bg-gradient-to-b from-gray-900 to-black">

        <!-- 顶部工具栏 -->
        <div class="flex-shrink-0 h-14 flex items-center justify-between px-6 bg-black/40 backdrop-blur-sm border-b border-white/10">
          <div class="flex items-center gap-4">
            <span class="text-white/90 text-sm font-medium truncate max-w-md">{{ media.name }}</span>
            <span class="text-white/50 text-xs">{{ formatSize(media.size) }}</span>
          </div>
          <div class="flex items-center gap-3">
            <!-- 进度 -->
            <div class="flex items-center gap-2 text-white/60 text-sm">
              <span>{{ currentIndex + 1 }} / {{ total }}</span>
              <div class="w-20 h-1.5 bg-white/20 rounded-full overflow-hidden">
                <div
                  class="h-full bg-emerald-500 rounded-full transition-all duration-300"
                  :style="{ width: `${progressPercent}%` }"
                />
              </div>
            </div>

            <!-- 撤销按钮 -->
            <button
              v-if="undoCount > 0"
              class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-white/10 text-white/80
                     hover:bg-white/20 transition-colors text-sm"
              @click="$emit('undo')">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
              </svg>
              撤销 ({{ undoCount }})
            </button>

            <!-- 关闭按钮 -->
            <button
              class="p-2 rounded-lg text-white/70 hover:text-white hover:bg-white/10 transition-colors"
              @click="$emit('close')">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- 主内容区域 -->
        <div class="flex-1 relative flex items-center justify-center overflow-hidden"
          @mousemove="handleMouseMove"
          @mouseup="handleMouseUp"
          @mouseleave="handleMouseUp"
          @wheel.prevent="handleWheel">

          <!-- 图片预览 -->
          <Transition name="slide" mode="out-in">
            <img
              :key="media.url"
              :src="media.url"
              :alt="media.name"
              class="max-w-full max-h-full object-contain select-none"
              :class="[
                isDragging ? 'cursor-grabbing' : (scale > 1 ? 'cursor-grab' : 'cursor-default'),
                isDragging ? '' : 'transition-transform duration-200'
              ]"
              :style="{
                transform: `translate(${offsetX}px, ${offsetY}px) scale(${scale})`
              }"
              draggable="false"
              @mousedown="handleMouseDown"
            />
          </Transition>

          <!-- 操作反馈提示 -->
          <Transition name="toast">
            <div
              v-if="lastAction"
              class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2
                     px-6 py-3 rounded-xl bg-black/80 text-white text-lg font-medium
                     backdrop-blur-sm pointer-events-none">
              {{ lastAction }}
            </div>
          </Transition>

          <!-- 处理中遮罩 -->
          <div
            v-if="isProcessing"
            class="absolute inset-0 bg-black/40 flex items-center justify-center">
            <div class="w-8 h-8 border-3 border-white/30 border-t-white rounded-full animate-spin" />
          </div>
        </div>

        <!-- 左侧导航 -->
        <button
          v-if="hasPrev"
          class="absolute left-4 top-1/2 -translate-y-1/2 p-3 rounded-full
                 bg-white/10 hover:bg-white/20 text-white/80 hover:text-white
                 transition-all backdrop-blur-sm z-10"
          @click="$emit('prev')">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>

        <!-- 右侧导航 -->
        <button
          v-if="hasNext"
          class="absolute right-4 top-1/2 -translate-y-1/2 p-3 rounded-full
                 bg-white/10 hover:bg-white/20 text-white/80 hover:text-white
                 transition-all backdrop-blur-sm z-10"
          @click="$emit('next')">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>

        <!-- 底部快捷键提示栏 -->
        <div class="flex-shrink-0 px-6 py-4 bg-black/60 backdrop-blur-sm border-t border-white/10">
          <div class="flex items-center justify-center gap-3 flex-wrap">
            <!-- 快捷键按钮 -->
            <button
              v-for="shortcut in validShortcuts"
              :key="shortcut.key"
              class="flex items-center gap-2 px-4 py-2 rounded-xl
                     bg-white/10 hover:bg-white/20 text-white/90
                     transition-all hover:scale-105 active:scale-95"
              @click="$emit('move', shortcut.key)">
              <kbd class="w-7 h-7 flex items-center justify-center rounded-lg
                          bg-white/20 text-white font-bold text-sm uppercase">
                {{ shortcut.key }}
              </kbd>
              <span class="text-sm">{{ shortcut.label }}</span>
            </button>

            <!-- 分隔线 -->
            <div class="w-px h-8 bg-white/20 mx-2" />

            <!-- 跳过按钮 -->
            <button
              class="flex items-center gap-2 px-4 py-2 rounded-xl
                     bg-gray-500/30 hover:bg-gray-500/50 text-white/70
                     transition-all"
              @click="$emit('skip')">
              <kbd class="px-2 py-1 rounded-lg bg-white/20 text-white text-xs">空格</kbd>
              <span class="text-sm">跳过</span>
            </button>

            <!-- 撤销提示 -->
            <div class="flex items-center gap-1 text-white/50 text-xs ml-4">
              <kbd class="px-1.5 py-0.5 rounded bg-white/10">⌘Z</kbd>
              <span>撤销</span>
            </div>
          </div>

          <!-- 分类统计 -->
          <div v-if="Object.keys(stats.categories).length > 0"
            class="flex items-center justify-center gap-4 mt-3 text-xs text-white/50">
            <span>已处理: {{ stats.processed }}</span>
            <span v-for="(count, label) in stats.categories" :key="label">
              {{ label }}: {{ count }}
            </span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script lang="ts" setup>
import { computed, watch } from 'vue'
import type { MediaInfo, ShortcutConfig, ClassifyStats } from '@/types'

const props = defineProps<{
  isOpen: boolean
  media: MediaInfo | null
  currentIndex: number
  total: number
  scale: number
  offsetX: number
  offsetY: number
  isDragging: boolean
  isProcessing: boolean
  lastAction: string
  undoCount: number
  stats: ClassifyStats
  hasNext: boolean
  hasPrev: boolean
  shortcuts: ShortcutConfig[]
}>()

const emit = defineEmits<{
  'close': []
  'prev': []
  'next': []
  'skip': []
  'move': [key: string]
  'undo': []
  'zoom-in': []
  'zoom-out': []
  'start-drag': []
  'end-drag': []
  'drag': [deltaX: number, deltaY: number]
}>()

// 有效的快捷键（已配置目标目录）
const validShortcuts = computed(() => {
  return props.shortcuts.filter(s => s.key && s.targetDir)
})

// 进度百分比
const progressPercent = computed(() => {
  if (props.total === 0) return 0
  return Math.round(((props.total - props.currentIndex) / props.total) * 100)
})

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

// 滚轮缩放
let lastWheelTime = 0
const WHEEL_THROTTLE_MS = 100

function handleWheel(e: WheelEvent) {
  const now = Date.now()
  if (now - lastWheelTime < WHEEL_THROTTLE_MS) return
  lastWheelTime = now

  if (e.deltaY < 0) {
    emit('zoom-in')
  } else {
    emit('zoom-out')
  }
}

// 格式化文件大小
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

.toast-enter-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.toast-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translate(-50%, -50%) scale(0.8);
}

.toast-leave-to {
  opacity: 0;
  transform: translate(-50%, -50%) scale(1.1);
}

kbd {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.border-3 {
  border-width: 3px;
}
</style>

