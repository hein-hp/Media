import {ref} from "vue";
import {EventsOn} from "../../wailsjs/runtime";

// 全局状态，在模块加载时就创建
const selectedDir = ref<string>("");

// 在模块加载时就注册事件监听，确保不会错过事件
EventsOn("selected-dir", (data: string) => {
  selectedDir.value = data;
});

/**
 * 选中文件路径 composable 监听 selected-dir 事件，获取选中文件路径
 */
export function useSelectedDir() {
  return {
    selectedDir,
  };
}
