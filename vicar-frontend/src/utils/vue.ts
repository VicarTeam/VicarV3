import {h} from "vue";

export function renderIcon(icon: string, type: string = "fas") {
  return () => {
    return h("i", {class: `${type} fa-${icon}`});
  }
}