import {ref} from "vue";
import {EventsOn} from "../../wailsjs/runtime";
import type {MediaInfo} from "@/types";

// 全局状态，在模块加载时就创建
const mediaList = ref<MediaInfo[]>([]);
const mediaCount = ref(0);

// 在模块加载时就注册事件监听，确保不会错过事件
EventsOn("media-list", (data: MediaInfo[]) => {
  mediaList.value = data;
});

EventsOn("media-count", (data: number) => {
  mediaCount.value = data;
});

/**
 * 媒体列表管理 composable 监听 media-list 事件，管理媒体列表状态
 */
export function useMediaList() {
  return {
    mediaList,
    mediaCount
  };
}
