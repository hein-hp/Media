<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100">
    <Header>
      <div class="flex items-center gap-3">
        <span>快捷分类</span>
        <span v-if="mediaList.length > 0" class="text-sm text-gray-500 font-normal">
          {{ mediaList.length }} 张图片待处理
        </span>
      </div>
    </Header>

    <!-- 主内容区域 -->
    <main class="pb-16 px-6 pt-4">
      <!-- 工具栏 -->
      <div class="flex items-center justify-between mb-6">
        <div class="flex items-center gap-3">
          <!-- 开始分类按钮 -->
          <button
            v-if="mediaList.length > 0"
            class="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-emerald-500 text-white
                   hover:bg-emerald-600 transition-all shadow-lg shadow-emerald-500/25
                   hover:shadow-emerald-500/40 hover:-translate-y-0.5"
            @click="startClassify">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            开始分类
          </button>
        </div>

        <!-- 设置按钮 -->
        <button
          class="flex items-center gap-2 px-4 py-2 rounded-xl bg-white text-gray-600
                 hover:bg-gray-50 border border-gray-200 transition-all shadow-sm"
          @click="showSettings = true">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          快捷键设置
        </button>
      </div>

      <!-- 快捷键预览 -->
      <div v-if="shortcuts.length > 0" class="mb-6">
        <h3 class="text-sm font-medium text-gray-500 mb-3">已配置的快捷键</h3>
        <div class="flex flex-wrap gap-2">
          <div
            v-for="shortcut in validShortcuts"
            :key="shortcut.key"
            class="flex items-center gap-2 px-3 py-2 rounded-lg bg-white border border-gray-200 shadow-sm">
            <kbd class="w-7 h-7 flex items-center justify-center rounded-lg
                        bg-gray-100 text-gray-700 font-bold text-sm uppercase">
              {{ shortcut.key }}
            </kbd>
            <span class="text-sm text-gray-600">{{ shortcut.label }}</span>
            <span class="text-xs text-gray-400 truncate max-w-32">{{ formatPath(shortcut.targetDir) }}</span>
          </div>
        </div>
      </div>

      <!-- 图片网格预览 -->
      <MediaGrid
        v-if="mediaList.length > 0"
        :images="mediaList"
        @select="handleSelect"
      />

      <!-- 空状态 -->
      <div v-else class="flex flex-col items-center justify-center py-20">
        <div class="w-24 h-24 mb-6 rounded-full bg-gray-100 flex items-center justify-center">
          <svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
              d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-700 mb-2">请先选择文件夹</h3>
        <p class="text-gray-500 text-sm">通过菜单「文件 → 选择文件夹」选择要分类的图片目录</p>
      </div>
    </main>

    <!-- 底部栏 -->
    <Footer
      :status="getStatusText()"
      :status-type="mediaList.length > 0 ? 'success' : 'info'">
      {{ selectedDir }}
    </Footer>

    <!-- 分类预览器 -->
    <ClassifyViewer
      :is-open="viewer.isOpen.value"
      :media="viewer.currentMedia.value"
      :current-index="viewer.currentIndex.value"
      :total="mediaList.length"
      :scale="viewer.scale.value"
      :offset-x="viewer.offsetX.value"
      :offset-y="viewer.offsetY.value"
      :is-dragging="viewer.isDragging.value"
      :is-processing="viewer.isProcessing.value"
      :last-action="viewer.lastAction.value"
      :undo-count="viewer.undoCount.value"
      :stats="viewer.stats.value"
      :has-next="viewer.hasNext.value"
      :has-prev="viewer.hasPrev.value"
      :shortcuts="shortcuts"
      @close="viewer.close"
      @prev="viewer.prev"
      @next="viewer.next"
      @skip="viewer.skip"
      @move="handleMove"
      @undo="handleUndo"
      @zoom-in="viewer.zoomIn"
      @zoom-out="viewer.zoomOut"
      @start-drag="viewer.startDrag"
      @end-drag="viewer.endDrag"
      @drag="viewer.drag"
    />

    <!-- 快捷键设置面板 -->
    <ShortcutSettings
      :is-open="showSettings"
      :shortcuts="shortcuts"
      @close="showSettings = false"
      @saved="handleShortcutsSaved"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { MediaGrid, ClassifyViewer, ShortcutSettings } from '@/components'
import { useMediaList, useSelectedDir, useClassifyViewer } from '@/composables'
import { Footer, Header } from '@/layout'
import type { ShortcutConfig } from '@/types'
import { GetShortcuts } from '../../wailsjs/go/app/App'

// 状态
const showSettings = ref(false)
const shortcuts = ref<ShortcutConfig[]>([])

// 选中文件夹
const { selectedDir } = useSelectedDir()

// 媒体列表
const { mediaList } = useMediaList()

// 分类查看器
const viewer = useClassifyViewer(mediaList, shortcuts)

// 有效的快捷键
const validShortcuts = computed(() => {
  return shortcuts.value.filter(s => s.key && s.targetDir)
})

// 加载快捷键配置
async function loadShortcuts() {
  try {
    const data = await GetShortcuts()
    shortcuts.value = data || []
  } catch (error) {
    console.error('加载快捷键配置失败:', error)
  }
}

// 开始分类
function startClassify() {
  viewer.open(0)
}

// 选择图片
function handleSelect(index: number) {
  viewer.open(index)
}

// 移动文件
async function handleMove(key: string) {
  await viewer.moveByShortcut(key)
}

// 撤销
async function handleUndo() {
  await viewer.undo()
}

// 保存快捷键后刷新
function handleShortcutsSaved(saved: ShortcutConfig[]) {
  shortcuts.value = saved
}

// 获取状态文本
function getStatusText(): string {
  if (mediaList.value.length > 0) {
    return `${validShortcuts.value.length} 个快捷键已配置`
  }
  return '等待选择文件夹'
}

// 格式化路径显示
function formatPath(path: string): string {
  if (!path) return ''
  if (path.startsWith('.')) return path
  const parts = path.split('/')
  return parts.length > 2 ? `.../${parts.slice(-2).join('/')}` : path
}

onMounted(() => {
  loadShortcuts()
})
</script>

<style scoped>
kbd {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}
</style>

