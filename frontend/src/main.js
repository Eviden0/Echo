// // main.ts
// import { createApp } from 'vue'
// import ElementPlus from 'element-plus'
// import 'element-plus/dist/index.css'
// import App from './App.vue'
//
// const app = createApp(App)
//
// app.use(ElementPlus)
// app.mount('#app')

import { createApp } from 'vue'
import App from './App.vue'
import './style/reset.scss'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { createPinia } from 'pinia'
const pinia = createPinia()
const app = createApp(App)
app.use(ElementPlus)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
app.use(pinia)

app.mount('#app')
