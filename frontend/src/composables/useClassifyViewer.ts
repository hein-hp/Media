import { computed, onMounted, onUnmounted, ref, type Ref } from 'vue'
import type { MediaInfo, ShortcutConfig, ClassifyStats } from '@/types'
import { MoveByShortcut, UndoMove, GetUndoCount } from '../../wailsjs/go/app/App'

export interface ClassifyViewerOptions {
  step?: number
  maxScale?: number
  minScale?: number
}

/**
 * 分类查看器 composable
 * 支持快捷键分类、撤销操作
 */
export function useClassifyViewer(
  mediaList: Ref<MediaInfo[]>,
  shortcuts: Ref<ShortcutConfig[]>,
  options: ClassifyViewerOptions = {}
) {
  const {
    step = 0.1,
    maxScale = 3,
    minScale = 0.5
  } = options

  // 核心状态
  const isOpen = ref(false)
  const currentIndex = ref(0)
  const scale = ref(1)
  const offsetX = ref(0)
  const offsetY = ref(0)
  const isDragging = ref(false)
  const isProcessing = ref(false)
  const lastAction = ref<string>('')
  const undoCount = ref(0)

  // 分类统计
  const stats = ref<ClassifyStats>({
    processed: 0,
    total: 0,
    categories: {}
  })

  // 当前媒体
  const currentMedia = computed(() => {
    const list = mediaList.value
    if (list.length === 0 || currentIndex.value < 0 || currentIndex.value >= list.length) {
      return null
    }
    return list[currentIndex.value]
  })

  const hasNext = computed(() => currentIndex.value < mediaList.value.length - 1)
  const hasPrev = computed(() => currentIndex.value > 0)
  const progress = computed(() => {
    if (mediaList.value.length === 0) return 0
    return Math.round((stats.value.processed / mediaList.value.length) * 100)
  })

  // 重置变换
  function resetTransform() {
    scale.value = 1
    offsetX.value = 0
    offsetY.value = 0
  }

  // 打开查看器
  function open(index = 0) {
    const listLength = mediaList.value.length
    if (listLength === 0) return

    currentIndex.value = Math.max(0, Math.min(index, listLength - 1))
    stats.value.total = listLength
    resetTransform()
    isOpen.value = true
  }

  // 关闭查看器
  function close() {
    isOpen.value = false
    resetTransform()
  }

  // 下一张
  function next() {
    if (hasNext.value) {
      currentIndex.value++
      resetTransform()
    }
  }

  // 上一张
  function prev() {
    if (hasPrev.value) {
      currentIndex.value--
      resetTransform()
    }
  }

  // 跳过（空格键）
  function skip() {
    next()
    lastAction.value = '跳过'
  }

  // 缩放
  function zoomIn() {
    scale.value = Math.min(scale.value + step, maxScale)
  }

  function zoomOut() {
    scale.value = Math.max(scale.value - step, minScale)
  }

  function resetZoom() {
    resetTransform()
  }

  // 拖拽
  function startDrag() {
    if (scale.value > 1) {
      isDragging.value = true
    }
  }

  function endDrag() {
    isDragging.value = false
  }

  function drag(deltaX: number, deltaY: number) {
    if (isDragging.value && scale.value > 1) {
      offsetX.value += deltaX
      offsetY.value += deltaY
    }
  }

  // 通过快捷键移动文件
  async function moveByShortcut(key: string): Promise<boolean> {
    const media = currentMedia.value
    if (!media || isProcessing.value) return false

    // 查找快捷键配置
    const shortcut = shortcuts.value.find(s => s.key.toLowerCase() === key.toLowerCase())
    if (!shortcut || !shortcut.targetDir) {
      lastAction.value = `快捷键 ${key} 未配置`
      return false
    }

    isProcessing.value = true
    try {
      await MoveByShortcut(media.path, key)

      // 更新统计
      stats.value.processed++
      stats.value.categories[shortcut.label] = (stats.value.categories[shortcut.label] || 0) + 1

      // 从列表中移除
      const curIdx = currentIndex.value
      mediaList.value.splice(curIdx, 1)

      // 更新索引
      if (mediaList.value.length === 0) {
        close()
        return true
      }

      if (curIdx >= mediaList.value.length) {
        currentIndex.value = mediaList.value.length - 1
      }

      resetTransform()
      lastAction.value = `已移动到 ${shortcut.label}`
      await refreshUndoCount()
      return true
    } catch (error) {
      console.error('移动文件失败:', error)
      lastAction.value = '移动失败'
      return false
    } finally {
      isProcessing.value = false
    }
  }

  // 撤销
  async function undo(): Promise<boolean> {
    if (undoCount.value === 0 || isProcessing.value) return false

    isProcessing.value = true
    try {
      await UndoMove()
      stats.value.processed = Math.max(0, stats.value.processed - 1)
      lastAction.value = '已撤销'
      await refreshUndoCount()
      return true
    } catch (error) {
      console.error('撤销失败:', error)
      lastAction.value = '撤销失败'
      return false
    } finally {
      isProcessing.value = false
    }
  }

  // 刷新撤销数量
  async function refreshUndoCount() {
    try {
      undoCount.value = await GetUndoCount()
    } catch (error) {
      console.error('获取撤销数量失败:', error)
    }
  }

  // 键盘事件处理
  async function handleKeydown(e: KeyboardEvent) {
    if (!isOpen.value) return

    // 忽略输入框中的按键
    if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) {
      return
    }

    const key = e.key.toLowerCase()

    // 系统快捷键
    switch (e.key) {
      case 'Escape':
        close()
        return
      case 'ArrowRight':
        next()
        return
      case 'ArrowLeft':
        prev()
        return
      case ' ': // 空格跳过
        e.preventDefault()
        skip()
        return
      case '+':
      case '=':
        zoomIn()
        return
      case '-':
        zoomOut()
        return
    }

    // Cmd/Ctrl + Z 撤销
    if ((e.metaKey || e.ctrlKey) && key === 'z') {
      e.preventDefault()
      await undo()
      return
    }

    // 检查是否是配置的快捷键
    const shortcut = shortcuts.value.find(s => s.key.toLowerCase() === key)
    if (shortcut && shortcut.targetDir) {
      e.preventDefault()
      await moveByShortcut(key)
    }
  }

  onMounted(() => {
    document.addEventListener('keydown', handleKeydown)
    refreshUndoCount()
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeydown)
  })

  return {
    // 状态
    isOpen,
    currentIndex,
    currentMedia,
    scale,
    offsetX,
    offsetY,
    isDragging,
    isProcessing,
    lastAction,
    undoCount,
    stats,
    hasNext,
    hasPrev,
    progress,
    // 方法
    open,
    close,
    next,
    prev,
    skip,
    zoomIn,
    zoomOut,
    resetZoom,
    startDrag,
    endDrag,
    drag,
    moveByShortcut,
    undo,
    refreshUndoCount
  }
}

