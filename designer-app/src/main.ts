import { createSSRApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'

/*** UniApp entry file - creates Vue 3 SSR app with Pinia state management ***/
export function createApp() {
  const app = createSSRApp(App)
  const pinia = createPinia()
  
  app.use(pinia)
  
  return {
    app,
    pinia
  }
}
