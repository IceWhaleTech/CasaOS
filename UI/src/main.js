import Vue from 'vue'
import App from '@/App.vue'
import router from '@/router'
import store from '@/store'
import api from '@/service/api.js'
import Buefy from 'buefy'
import '@/assets/scss/app.scss'

Vue.use(Buefy)

Vue.config.productionTip = false
Vue.prototype.$api = api;
new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
