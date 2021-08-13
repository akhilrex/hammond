import axios from 'axios'
import { filter } from 'lodash'

import parseISO from 'date-fns/parseISO'
export const state = {
  vehicles: [],
  roleMasters: [],
  fuelUnitMasters: [],
  distanceUnitMasters: [],
  currencyMasters: [],
  fuelTypeMasters: [],
  quickEntries: [],
  vehicleStats: new Map(),
}

export const getters = {
  unprocessedQuickEntries: (state) => {
    return filter(state.quickEntries, (o) => o.processDate == null)
  },
  processedQuickEntries: (state) => {
    return filter(state.quickEntries, (o) => o.processDate != null)
  },
}

export const mutations = {
  CACHE_VEHICLE(state, newVehicles) {
    state.vehicles = newVehicles
  },
  CACHE_VEHICLE_STATS(state, stats) {
    state.vehicleStats.set(stats.vehicleId, stats)
  },
  CACHE_FUEL_UNIT_MASTERS(state, masters) {
    state.fuelUnitMasters = masters
  },
  CACHE_DISTANCE_UNIT_MASTERS(state, masters) {
    state.distanceUnitMasters = masters
  },
  CACHE_FUEL_TYPE_MASTERS(state, masters) {
    state.fuelTypeMasters = masters
  },
  CACHE_CURRENCY_MASTERS(state, masters) {
    state.currencyMasters = masters
  },
  CACHE_QUICK_ENTRIES(state, entries) {
    state.quickEntries = entries
  },
  CACHE_ROLE_MASTERS(state, roles) {
    state.roleMasters = roles
  },
}

export const actions = {
  init({ dispatch, rootState }) {
    const { currentUser } = rootState.auth
    if (currentUser) {
      dispatch('fetchVehicles')
      dispatch('fetchMasters')
      dispatch('fetchQuickEntries', { force: true })
    }
  },
  fetchMasters({ commit, state, rootState }) {
    return axios.get('/api/masters').then((response) => {
      const fuelUnitMasters = response.data.fuelUnits
      const fuelTypeMasters = response.data.fuelTypes
      commit('CACHE_FUEL_UNIT_MASTERS', fuelUnitMasters)
      commit('CACHE_FUEL_TYPE_MASTERS', fuelTypeMasters)
      commit('CACHE_CURRENCY_MASTERS', response.data.currencies)
      commit('CACHE_DISTANCE_UNIT_MASTERS', response.data.distanceUnits)
      commit('CACHE_ROLE_MASTERS', response.data.roles)
      return response.data
    })
  },
  fetchVehicles({ commit, state, rootState }) {
    return axios.get('/api/me/vehicles').then((response) => {
      const data = response.data
      commit('CACHE_VEHICLE', data)
      return data
    })
  },
  fetchQuickEntries({ commit, state, rootState }, { force }) {
    if (state.quickEntries && !force) {
      return Promise.resolve(state.quickEntries)
    }
    return axios.get('/api/me/quickEntries').then((response) => {
      const data = response.data
      commit('CACHE_QUICK_ENTRIES', data)
      return data
    })
  },
  fetchVehicleById({ commit, state, rootState }, { vehicleId }) {
    const matchedVehicle = state.vehicles.find((vehicle) => vehicle.id === vehicleId)
    if (matchedVehicle) {
      return Promise.resolve(matchedVehicle)
    }
    return axios.get('/api/vehicles/' + vehicleId).then((response) => {
      const data = response.data
      // commit('CACHE_VEHICLE', data)
      return data
    })
  },
  fetchFillupById({ commit, state, rootState }, { vehicleId, fillupId }) {
    return axios.get(`/api/vehicles/${vehicleId}/fillups/${fillupId}`).then((response) => {
      const data = response.data
      data.date = parseISO(data.date)
      return data
    })
  },
  deleteFillupById({ commit, state, rootState }, { vehicleId, fillupId }) {
    return axios.delete(`/api/vehicles/${vehicleId}/fillups/${fillupId}`).then((response) => {
      const data = response.data
      return data
    })
  },
  fetchExpenseById({ commit, state, rootState }, { vehicleId, expenseId }) {
    return axios.get(`/api/vehicles/${vehicleId}/expenses/${expenseId}`).then((response) => {
      const data = response.data
      data.date = parseISO(data.date)
      return data
    })
  },
  deleteExpenseById({ commit, state, rootState }, { vehicleId, expenseId }) {
    return axios.delete(`/api/vehicles/${vehicleId}/expenses/${expenseId}`).then((response) => {
      const data = response.data
      return data
    })
  },
  fetchAttachmentsByVehicleId({ commit, state, rootState }, { vehicleId }) {
    return axios.get(`/api/vehicles/${vehicleId}/attachments`).then((response) => {
      const data = response.data

      return data
    })
  },
  fetchUsersByVehicleId({ commit, state, rootState }, { vehicleId, force }) {
    return axios.get(`/api/vehicles/${vehicleId}/users`).then((response) => {
      const data = response.data
      // data.vehicleId = vehicleId
      // commit('CACHE_VEHICLE_STATS', data)

      return data
    })
  },
  fetchFuelSubtypesByVehicleId({ commit, state, rootState }, { vehicleId, force }) {
    return axios.get(`/api/vehicles/${vehicleId}/fuelSubTypes`).then((response) => {
      const data = response.data
      return data
    })
  },
  fetchStatsByVehicleId({ commit, state, rootState }, { vehicleId, force }) {
    if (state.vehicleStats.has(vehicleId) && !force) {
      return Promise.resolve(state.vehicleStats.get(vehicleId))
    }
    return axios.get(`/api/vehicles/${vehicleId}/stats`).then((response) => {
      const data = response.data
      data.vehicleId = vehicleId
      commit('CACHE_VEHICLE_STATS', data)

      return data
    })
  },
  setQuickEntryAsProcessed({ commit, state, rootState, dispatch }, { id }) {
    return axios.post(`/api/quickEntries/${id}/process`, {}).then((response) => {
      const data = response.data
      dispatch('fetchQuickEntries', { force: true })
      return data
    })
  },
  deleteQuickEntry({ commit, state, rootState, dispatch }, { id }) {
    return axios.delete(`/api/quickEntries/${id}`).then((response) => {
      const data = response.data
      dispatch('fetchQuickEntries', { force: true })
      return data
    })
  },
}
