/// <reference types="vite/client" />
// wangeditor 类型声明
declare module '@wangeditor/editor-for-vue' {
  import { DefineComponent } from 'vue'
  export const Editor: DefineComponent<any, any, any>
  export const Toolbar: DefineComponent<any, any, any>
}
