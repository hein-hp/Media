import { ref, onMounted } from 'vue'
import type { ShortcutConfig } from '@/types'
import {
  GetShortcuts,
  SaveShortcuts,
  SelectShortcutTargetDir,
  MoveByShortcut,
  UndoMove,
  GetUndoCount
} from '../../wailsjs/go/app/App'

/**
 * 快捷键管理 composable
 */
export function useShortcuts() {
  const shortcuts = ref<ShortcutConfig[]>([])
  const isLoading = ref(false)
  const undoCount = ref(0)

  // 加载快捷键配置
  async function loadShortcuts() {
    isLoading.value = true
    try {
      const data = await GetShortcuts()
      shortcuts.value = data || []
    } catch (error) {
      console.error('加载快捷键配置失败:', error)
    } finally {
      isLoading.value = false
    }
  }

  // 保存快捷键配置
  async function saveShortcuts(configs: ShortcutConfig[]) {
    try {
      await SaveShortcuts(configs)
      shortcuts.value = configs
      return true
    } catch (error) {
      console.error('保存快捷键配置失败:', error)
      return false
    }
  }

  // 添加快捷键
  function addShortcut(config: ShortcutConfig) {
    shortcuts.value.push(config)
  }

  // 删除快捷键
  function removeShortcut(key: string) {
    const index = shortcuts.value.findIndex(s => s.key === key)
    if (index !== -1) {
      shortcuts.value.splice(index, 1)
    }
  }

  // 更新快捷键
  function updateShortcut(key: string, config: Partial<ShortcutConfig>) {
    const shortcut = shortcuts.value.find(s => s.key === key)
    if (shortcut) {
      Object.assign(shortcut, config)
    }
  }

  // 选择目标文件夹
  async function selectTargetDir(): Promise<string> {
    try {
      return await SelectShortcutTargetDir()
    } catch (error) {
      console.error('选择文件夹失败:', error)
      return ''
    }
  }

  // 通过快捷键移动文件
  async function moveFile(filePath: string, shortcutKey: string): Promise<boolean> {
    try {
      await MoveByShortcut(filePath, shortcutKey)
      await refreshUndoCount()
      return true
    } catch (error) {
      console.error('移动文件失败:', error)
      return false
    }
  }

  // 撤销移动
  async function undo(): Promise<boolean> {
    try {
      await UndoMove()
      await refreshUndoCount()
      return true
    } catch (error) {
      console.error('撤销失败:', error)
      return false
    }
  }

  // 刷新可撤销数量
  async function refreshUndoCount() {
    try {
      undoCount.value = await GetUndoCount()
    } catch (error) {
      console.error('获取撤销数量失败:', error)
    }
  }

  // 根据快捷键获取配置
  function getShortcutByKey(key: string): ShortcutConfig | undefined {
    return shortcuts.value.find(s => s.key.toLowerCase() === key.toLowerCase())
  }

  // 检查快捷键是否有效（已配置目标目录）
  function isValidShortcut(key: string): boolean {
    const shortcut = getShortcutByKey(key)
    return !!shortcut && !!shortcut.targetDir
  }

  onMounted(() => {
    loadShortcuts()
    refreshUndoCount()
  })

  return {
    shortcuts,
    isLoading,
    undoCount,
    loadShortcuts,
    saveShortcuts,
    addShortcut,
    removeShortcut,
    updateShortcut,
    selectTargetDir,
    moveFile,
    undo,
    refreshUndoCount,
    getShortcutByKey,
    isValidShortcut
  }
}

