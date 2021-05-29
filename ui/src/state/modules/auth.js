import axios from 'axios'

export const state = {
  currentUser: getSavedState('auth.currentUser'),
  initialized: getSavedState('system.initialized'),
}

export const mutations = {
  SET_CURRENT_USER(state, newValue) {
    state.currentUser = newValue
    saveState('auth.currentUser', newValue)
    setDefaultAuthHeaders(state)
  },
  SET_INITIALIZATION_STATUS(state, newValue) {
    state.initialized = newValue
    saveState('system.initialized', newValue)
  },
}

export const getters = {
  // Whether the user is currently logged in.
  loggedIn(state) {
    return !!state.currentUser
  },
  isInitialized(state) {
    return state.initialized == null || state.initialized.initialized
  },
}

export const actions = {
  // This is automatically run in `src/state/store.js` when the app
  // starts, along with any other actions named `init` in other modules.
  init({ state, dispatch }) {
    dispatch('systemInitialized')
    setDefaultAuthHeaders(state)
    dispatch('validate')
  },

  logIn({ commit, dispatch, getters }, { username, password } = {}) {
    if (getters.loggedIn) return dispatch('validate')

    return axios.post('/api/login', { email: username, password }).then((response) => {
      const user = response.data
      commit('SET_CURRENT_USER', user)
      dispatch('vehicles/fetchMasters', null, { root: true })
      return user
    })
  },

  // Logs out the current user.
  logOut({ commit }) {
    commit('SET_CURRENT_USER', null)
  },

  // Validates the current user's token and refreshes it
  // with new data from the API.
  validate({ commit, state }) {
    if (!state.currentUser) return Promise.resolve(null)

    return axios
      .post('/api/refresh', { refreshToken: state.currentUser.refreshToken })
      .then((response) => {
        const user = response.data
        commit('SET_CURRENT_USER', user)
        return user
      })
      .catch((ex) => {
        commit('SET_CURRENT_USER', null)
      })
  },

  systemInitialized({ commit, state }) {
    return axios.get('/api/system/status').then((response) => {
      const data = response.data
      commit('SET_INITIALIZATION_STATUS', data)
      return data
    })
  },
}

// ===
// Private helpers
// ===

function getSavedState(key) {
  return JSON.parse(window.localStorage.getItem(key))
}

function saveState(key, state) {
  window.localStorage.setItem(key, JSON.stringify(state))
}

function setDefaultAuthHeaders(state) {
  axios.defaults.headers.common.Authorization = state.currentUser ? state.currentUser.token : ''
}
