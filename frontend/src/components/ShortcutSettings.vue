<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="isOpen"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
        @click.self="$emit('close')">
        <div class="bg-white rounded-2xl shadow-2xl w-[600px] max-h-[80vh] overflow-hidden">
          <!-- 头部 -->
          <div class="flex items-center justify-between px-6 py-4 border-b border-gray-100">
            <h2 class="text-lg font-semibold text-gray-800">快捷键设置</h2>
            <button
              class="p-2 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors"
              @click="$emit('close')">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- 内容区域 -->
          <div class="px-6 py-4 max-h-[60vh] overflow-y-auto">
            <!-- 快捷键列表 -->
            <div class="space-y-3">
              <div
                v-for="(shortcut, index) in localShortcuts"
                :key="index"
                class="flex items-center gap-3 p-3 rounded-xl bg-gray-50 hover:bg-gray-100 transition-colors group">
                <!-- 快捷键输入 -->
                <div class="flex-shrink-0">
                  <input
                    v-model="shortcut.key"
                    type="text"
                    maxlength="1"
                    class="w-12 h-12 text-center text-lg font-bold rounded-xl border-2 border-gray-200
                           focus:border-blue-500 focus:ring-2 focus:ring-blue-200 outline-none
                           uppercase bg-white transition-all"
                    placeholder="键"
                    @input="onKeyInput($event, index)"
                  />
                </div>

                <!-- 标签输入 -->
                <input
                  v-model="shortcut.label"
                  type="text"
                  class="flex-shrink-0 w-24 px-3 py-2 rounded-lg border border-gray-200
                         focus:border-blue-500 focus:ring-2 focus:ring-blue-200 outline-none
                         bg-white transition-all text-sm"
                  placeholder="标签名"
                />

                <!-- 目录显示/选择 -->
                <div class="flex-1 flex items-center gap-2">
                  <input
                    v-model="shortcut.targetDir"
                    type="text"
                    class="flex-1 px-3 py-2 rounded-lg border border-gray-200
                           focus:border-blue-500 focus:ring-2 focus:ring-blue-200 outline-none
                           bg-white transition-all text-sm text-gray-600"
                    placeholder="目标文件夹路径（可相对路径如 .delete）"
                  />
                  <button
                    class="flex-shrink-0 px-3 py-2 rounded-lg bg-blue-500 text-white text-sm
                           hover:bg-blue-600 transition-colors"
                    @click="selectDir(index)">
                    选择
                  </button>
                </div>

                <!-- 删除按钮 -->
                <button
                  class="flex-shrink-0 p-2 rounded-lg text-gray-400 hover:text-red-500
                         hover:bg-red-50 transition-colors opacity-0 group-hover:opacity-100"
                  @click="removeShortcut(index)">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>

            <!-- 添加按钮 -->
            <button
              class="mt-4 w-full py-3 rounded-xl border-2 border-dashed border-gray-300
                     text-gray-500 hover:border-blue-400 hover:text-blue-500
                     transition-colors flex items-center justify-center gap-2"
              @click="addShortcut">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              添加快捷键
            </button>

            <!-- 使用提示 -->
            <div class="mt-6 p-4 rounded-xl bg-blue-50 text-sm text-blue-700">
              <p class="font-medium mb-2">💡 使用提示</p>
              <ul class="space-y-1 text-blue-600">
                <li>• 在预览图片时按对应快捷键，图片将移动到目标文件夹</li>
                <li>• 支持相对路径：如 <code class="px-1 py-0.5 bg-blue-100 rounded">.delete</code> 会在当前目录创建</li>
                <li>• 按 <kbd class="px-1.5 py-0.5 bg-white rounded border">Cmd+Z</kbd> 可撤销操作</li>
                <li>• 按 <kbd class="px-1.5 py-0.5 bg-white rounded border">空格</kbd> 跳过当前图片</li>
              </ul>
            </div>
          </div>

          <!-- 底部操作 -->
          <div class="flex items-center justify-end gap-3 px-6 py-4 border-t border-gray-100 bg-gray-50">
            <button
              class="px-4 py-2 rounded-lg text-gray-600 hover:bg-gray-200 transition-colors"
              @click="$emit('close')">
              取消
            </button>
            <button
              class="px-6 py-2 rounded-lg bg-blue-500 text-white hover:bg-blue-600 transition-colors
                     disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="isSaving"
              @click="save">
              {{ isSaving ? '保存中...' : '保存' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import type { ShortcutConfig } from '@/types'
import { SaveShortcuts, SelectShortcutTargetDir } from '../../wailsjs/go/app/App'

const props = defineProps<{
  isOpen: boolean
  shortcuts: ShortcutConfig[]
}>()

const emit = defineEmits<{
  'close': []
  'saved': [shortcuts: ShortcutConfig[]]
}>()

const localShortcuts = ref<ShortcutConfig[]>([])
const isSaving = ref(false)

// 同步外部数据
watch(() => props.isOpen, (open) => {
  if (open) {
    localShortcuts.value = JSON.parse(JSON.stringify(props.shortcuts))
  }
})

// 处理按键输入（只允许字母和数字）
function onKeyInput(event: Event, index: number) {
  const input = event.target as HTMLInputElement
  const value = input.value.replace(/[^a-zA-Z0-9]/g, '').slice(0, 1)
  localShortcuts.value[index].key = value
}

// 添加快捷键
function addShortcut() {
  localShortcuts.value.push({
    key: '',
    targetDir: '',
    label: ''
  })
}

// 删除快捷键
function removeShortcut(index: number) {
  localShortcuts.value.splice(index, 1)
}

// 选择目录
async function selectDir(index: number) {
  try {
    const dir = await SelectShortcutTargetDir()
    if (dir) {
      localShortcuts.value[index].targetDir = dir
    }
  } catch (error) {
    console.error('选择目录失败:', error)
  }
}

// 保存
async function save() {
  // 过滤掉空的配置
  const validShortcuts = localShortcuts.value.filter(s => s.key && s.label)

  isSaving.value = true
  try {
    await SaveShortcuts(validShortcuts)
    emit('saved', validShortcuts)
    emit('close')
  } catch (error) {
    console.error('保存失败:', error)
  } finally {
    isSaving.value = false
  }
}
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

kbd {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.75rem;
}

code {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}
</style>

