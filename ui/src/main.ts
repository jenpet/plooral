import Vue, { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

// Bootstrap specific imports
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

// import './app.scss'
// TODO: Not working
// $b-custom-control-indicator-size-lg: $custom-control-indicator-size * 1.25 !default;

createApp(App)
    .use(store)
    .use(router)
    .mount('#app')
