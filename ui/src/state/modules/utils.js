import axios from 'axios'

export const state = {
  isMobile: false,
  settings: null,
}
export const mutations = {
  CACHE_ISMOBILE(state, isMobile) {
    state.isMobile = isMobile
  },
  CACHE_SETTINGS(state, settings) {
    state.settings = settings
  },
}
export const getters = {}
export const actions = {
  init({ dispatch, rootState }) {
    dispatch('checkSize')
    const { currentUser } = rootState.auth
    if (currentUser) {
      dispatch('getSettings')
    }
  },
  checkSize({ commit }) {
    commit('CACHE_ISMOBILE', window.innerWidth < 600)
    return window.innerWidth < 600
  },
  getSettings({ commit }) {
    return axios.get(`/api/settings`).then((response) => {
      const data = response.data
      commit('CACHE_SETTINGS', data)
      return data
    })
  },
  saveSettings({ commit, dispatch }, { settings }) {
    return axios.post(`/api/settings`, { ...settings }).then((response) => {
      const data = response.data
      dispatch('getSettings')
      return data
    })
  },
  saveUserSettings({ commit, dispatch }, { settings }) {
    return axios.post(`/api/me/settings`, { ...settings }).then((response) => {
      const data = response.data
      dispatch('users/forceMe', {}, { root: true }).then((data) => {})
      return data
    })
  },
}
