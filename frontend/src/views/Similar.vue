<template>
  <div class="min-h-screen bg-white">
    <Header>相似图片分析 {{ similarGroups.length > 0 ? `- 共 ${similarGroups.length} 组` : '' }}</Header>

    <!-- 主内容区域 -->
    <main class="pb-12">
      <SimilarGroups
        :groups="similarGroups"
        :loading="isLoading"
        :is-deleting="isDeleting"
        @remove="handleRemove"
        @remove-smaller="handleRemoveSmaller"
      />
    </main>

    <!-- 底部栏 -->
    <Footer
      :status="getStatusText()"
      :status-type="similarGroups.length > 0 ? 'success' : 'info'"
    >
      {{ selectedDir }}
    </Footer>
  </div>
</template>

<script lang="ts" setup>
import {SimilarGroups} from '@/components'
import {useSelectedDir, useSimilarImages} from '@/composables'
import {Footer, Header} from '@/layout'

// 选中文件夹
const {selectedDir} = useSelectedDir()

// 相似图片状态
const {similarGroups, isLoading, isDeleting, removeImage, removeSmallerImages} = useSimilarImages()

// 处理删除单张图片
async function handleRemove(groupId: number, imagePath: string) {
  await removeImage(groupId, imagePath)
}

// 处理删除较小文件
async function handleRemoveSmaller() {
  await removeSmallerImages()
}

// 获取状态文本
function getStatusText(): string {
  if (isDeleting.value) {
    return '正在删除...'
  }
  if (isLoading.value) {
    return '正在分析...'
  }
  if (similarGroups.value.length > 0) {
    const totalImages = similarGroups.value.reduce((total, group) => total + group.images.length, 0)
    return `已找到 ${similarGroups.value.length} 组相似图片，共 ${totalImages} 张`
  }
  return '等待分析'
}
</script>

