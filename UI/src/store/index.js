/*
 * @Author: JerryK
 * @Date: 2021-09-18 21:32:13
 * @LastEditors: JerryK
 * @LastEditTime: 2021-09-22 16:28:16
 * @Description: 
 * @FilePath: /CasaOS-UI/src/store/index.js
 */
import Vue from 'vue'
import Vuex from 'vuex'
//import createPersistedState from "vuex-persistedstate";

Vue.use(Vuex)

export default new Vuex.Store({
  //plugins: [createPersistedState()],
  state: {
    token: "",
    devIp: "192.168.2.217",
    serviceError: false
  },
  mutations: {
    setToken(state, val) {
      state.token = val 
    },
    setServiceError(state, val) {
      state.serviceError = val
    }
  },
  actions: {
  },
  modules: {
  }
})
