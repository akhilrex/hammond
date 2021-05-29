import axios from 'axios'

export const state = {
  cached: [],
  me: null,
}

export const getters = {}

export const mutations = {
  CACHE_USER(state, newUser) {
    state.cached.push(newUser)
  },
  CACHE_MY_USER(state, newUser) {
    state.me = newUser
  },
}

export const actions = {
  init({ dispatch, rootState }) {
    const { currentUser } = rootState.auth
    if (currentUser != null) {
      dispatch('me')
    }
  },
  forceMe({ commit, state }) {
    return axios
      .get('/api/me')
      .then((response) => {
        commit('CACHE_MY_USER', response.data)
        return response.data
      })
      .catch((error) => {
        if (error.response && error.response.status === 401) {
          commit('CACHE_MY_USER', null)
        } else {
          console.warn(error)
        }
        return null
      })
  },
  users() {
    return axios
      .get('/api/users')
      .then((response) => {
        return response.data
      })
      .catch((error) => {
        if (error.response && error.response.status === 401) {
        } else {
          console.warn(error)
        }
        return null
      })
  },
  me({ commit, state }) {
    if (state.me) {
      return Promise.resolve(state.me)
    }
    return axios
      .get('/api/me')
      .then((response) => {
        commit('CACHE_MY_USER', response.data)
        return response.data
      })
      .catch((error) => {
        if (error.response && error.response.status === 401) {
          commit('CACHE_MY_USER', null)
        } else {
          console.warn(error)
        }
        return null
      })
  },
  fetchUser({ commit, state, rootState }, { username }) {
    // 1. Check if we already have the user as a current user.
    const { currentUser } = rootState.auth
    if (currentUser && currentUser.username === username) {
      return Promise.resolve(currentUser)
    }

    // 2. Check if we've already fetched and cached the user.
    const matchedUser = state.cached.find((user) => user.username === username)
    if (matchedUser) {
      return Promise.resolve(matchedUser)
    }
  },
}
