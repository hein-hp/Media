import {ref} from "vue";
import {EventsOn} from "../../wailsjs/runtime";
import type {SimilarityResult} from "@/types";
import {RemoveSimilarImage} from "../../wailsjs/go/app/App";

// 全局状态，在模块加载时就创建
const similarGroups = ref<SimilarityResult[]>([]);
const isLoading = ref(false);
const isDeleting = ref(false);

// 在模块加载时就注册事件监听，确保不会错过事件
EventsOn("similar-results", (data: SimilarityResult[]) => {
  similarGroups.value = data || [];
  isLoading.value = false;
});

EventsOn("similar-loading", (loading: boolean) => {
  isLoading.value = loading;
});

/**
 * 相似图片管理 composable
 * 监听 similar-results 事件，管理相似图片状态
 */
export function useSimilarImages() {

  /**
   * 删除相似图片
   */
  async function removeImage(groupId: number, imagePath: string) {
    try {
      await RemoveSimilarImage(imagePath);
      // 从本地状态中移除
      const group = similarGroups.value.find((g) => g.groupId === groupId);
      if (group) {
        group.images = group.images.filter((img) => img.path !== imagePath);
        // 如果组中只剩一张图片，移除整个组
        if (group.images.length < 2) {
          similarGroups.value = similarGroups.value.filter(
            (g) => g.groupId !== groupId
          );
        }
      }
    } catch (error) {
      console.error("删除图片失败:", error);
    }
  }

  /**
   * 清空相似图片结果
   */
  function clearResults() {
    similarGroups.value = [];
  }

  /**
   * 获取总组数
   */
  function getTotalGroups() {
    return similarGroups.value.length;
  }

  /**
   * 获取总图片数
   */
  function getTotalImages() {
    return similarGroups.value.reduce(
      (total, group) => total + group.images.length,
      0
    );
  }

  /**
   * 删除每组中较小的文件，只保留最大的
   */
  async function removeSmallerImages() {
    if (isDeleting.value) return;
    
    isDeleting.value = true;
    
    try {
      // 收集所有需要删除的图片
      const toDelete: { groupId: number; path: string }[] = [];
      
      for (const group of similarGroups.value) {
        if (group.images.length < 2) continue;
        
        // 找出最大的图片
        const sortedImages = [...group.images].sort((a, b) => b.size - a.size);
        const largestImage = sortedImages[0];
        
        // 其他图片都要删除
        for (const img of group.images) {
          if (img.path !== largestImage.path) {
            toDelete.push({ groupId: group.groupId, path: img.path });
          }
        }
      }
      
      // 逐个删除
      for (const item of toDelete) {
        try {
          await RemoveSimilarImage(item.path);
          // 从本地状态中移除
          const group = similarGroups.value.find((g) => g.groupId === item.groupId);
          if (group) {
            group.images = group.images.filter((img) => img.path !== item.path);
          }
        } catch (error) {
          console.error("删除图片失败:", item.path, error);
        }
      }
      
      // 清理只剩一张图片的组
      similarGroups.value = similarGroups.value.filter((g) => g.images.length >= 2);
      
    } finally {
      isDeleting.value = false;
    }
  }

  return {
    similarGroups,
    isLoading,
    isDeleting,
    removeImage,
    removeSmallerImages,
    clearResults,
    getTotalGroups,
    getTotalImages,
  };
}

