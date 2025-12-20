import {onMounted, onUnmounted, ref} from "vue";
import {EventsOff, EventsOn} from "../../wailsjs/runtime";

export function useRoute() {
  const path = ref<string>("");

  onMounted(() => {
    EventsOn("router", (data: string) => {
      path.value = data;
      console.log("路由变化，当前路径：", data);
    });
  });

  onUnmounted(() => {
    EventsOff("router");
  });

  return {
    path,
  };
}
