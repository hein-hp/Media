<template>
  <div class="similar-groups">
    <!-- 统计信息和操作栏 -->
    <div v-if="groups.length > 0" class="mb-4 px-6 mt-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-4 text-sm text-gray-500">
          <span class="flex items-center gap-1.5">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
            </svg>
            共 {{ groups.length }} 组
          </span>
          <span class="flex items-center gap-1.5">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
            </svg>
            {{ totalImages }} 张图片
          </span>
          <span v-if="toDeleteCount > 0" class="flex items-center gap-1.5 text-red-500">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
            可清理 {{ toDeleteCount }} 张
          </span>
        </div>
        
        <!-- 批量删除按钮 -->
        <button
          @click="$emit('removeSmaller')"
          :disabled="isDeleting"
          class="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-red-500 to-red-600 
                 hover:from-red-600 hover:to-red-700 text-white text-sm font-medium rounded-lg
                 shadow-md hover:shadow-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg v-if="!isDeleting" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
          </svg>
          <svg v-else class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ isDeleting ? '删除中...' : '删除较小文件' }}
        </button>
      </div>
    </div>

    <!-- 分组列表 -->
    <div class="space-y-8 px-6 pb-6">
      <div
        v-for="group in groups"
        :key="group.groupId"
        class="group-card bg-gradient-to-br from-gray-50 to-white rounded-2xl border border-gray-100 
               shadow-sm hover:shadow-md transition-shadow duration-300 overflow-hidden"
      >
        <!-- 组标题 -->
        <div class="flex items-center justify-between px-5 py-3 bg-white/60 border-b border-gray-100">
          <div class="flex items-center gap-2">
            <div class="w-2 h-2 rounded-full bg-amber-400"></div>
            <span class="text-sm font-medium text-gray-700">第 {{ group.groupId }} 组</span>
            <span class="text-xs text-gray-400">{{ group.images.length }} 张相同图片</span>
          </div>
        </div>

        <!-- 图片网格 -->
        <div class="p-4">
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-3">
            <div
              v-for="image in group.images"
              :key="image.path"
              class="relative group/item aspect-square rounded-xl overflow-hidden bg-gray-100
                     border border-gray-200 hover:border-amber-300 transition-all duration-200"
            >
              <!-- 图片 -->
              <img
                :src="image.url"
                :alt="image.name"
                class="w-full h-full object-cover transition-transform duration-300 group-hover/item:scale-105"
                loading="lazy"
              />

              <!-- 悬浮操作层 -->
              <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent
                          opacity-0 group-hover/item:opacity-100 transition-opacity duration-200">
                <!-- 文件信息 -->
                <div class="absolute bottom-0 left-0 right-0 p-3">
                  <p class="text-white text-xs font-medium truncate mb-1">{{ image.name }}</p>
                  <p class="text-white/70 text-xs">{{ formatSize(image.size) }}</p>
                </div>

                <!-- 删除按钮 -->
                <button
                  @click.stop="$emit('remove', group.groupId, image.path)"
                  class="absolute top-2 right-2 w-8 h-8 flex items-center justify-center
                         bg-red-500/90 hover:bg-red-600 rounded-full text-white
                         shadow-lg transform scale-90 hover:scale-100 transition-all duration-200"
                  title="删除此图片"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="groups.length === 0 && !loading" class="flex flex-col items-center justify-center py-20 text-gray-400">
      <svg class="w-16 h-16 mb-4 text-gray-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
      </svg>
      <p class="text-sm">暂无相似图片分析结果</p>
      <p class="text-xs text-gray-300 mt-1">使用菜单中的"查找相似"功能开始分析</p>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-20">
      <div class="w-10 h-10 border-3 border-amber-200 border-t-amber-500 rounded-full animate-spin mb-4"></div>
      <p class="text-sm text-gray-500">正在分析相似图片...</p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {computed} from 'vue'
import type {SimilarityResult} from '@/types'

const props = defineProps<{
  groups: SimilarityResult[]
  loading?: boolean
  isDeleting?: boolean
}>()

defineEmits<{
  remove: [groupId: number, imagePath: string]
  removeSmaller: []
}>()

/**
 * 计算总图片数
 */
const totalImages = computed(() => {
  return props.groups.reduce((total, group) => total + group.images.length, 0)
})

/**
 * 计算可删除的图片数量（每组中除了最大的都可删除）
 */
const toDeleteCount = computed(() => {
  return props.groups.reduce((total, group) => {
    // 每组中除了最大的一张，其他都要删除
    return total + Math.max(0, group.images.length - 1)
  }, 0)
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

<style scoped>
.group-card {
  animation: fadeInUp 0.4s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.border-3 {
  border-width: 3px;
}
</style>

