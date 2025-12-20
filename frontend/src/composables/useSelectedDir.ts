import {onMounted, onUnmounted, ref} from "vue";
import {EventsOff, EventsOn} from "../../wailsjs/runtime";

/**
 * 选中文件路径 composable 监听 selected-dir 事件，获取选中文件路径
 */
export function useSelectedDir() {
  const selectedDir = ref<string>("");

  onMounted(() => {
    EventsOn("selected-dir", (data: string) => {
      selectedDir.value = data;
    });
  });

  onUnmounted(() => {
    EventsOff("selected-dir");
  });

  return {
    selectedDir: selectedDir,
  };
}
