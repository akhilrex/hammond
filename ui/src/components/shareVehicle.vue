<script>
import store from '@state/store'
import { sortBy } from 'lodash'
import axios from 'axios'

export default {
  props: {
    vehicle: {
      type: Object,
      required: true,
    },
  },
  data: function() {
    return {
      models: [],
    }
  },
  mounted() {
    store
      .dispatch('users/users')
      .then((allUsers) => {
        store.dispatch('vehicles/fetchUsersByVehicleId', { vehicleId: this.vehicle.id }).then((data) => {
          const arr = []
          for (const user of allUsers) {
            const toAdd = {
              id: user.id,
              name: user.name,
              isShared: false,
              isOwner: false,
            }
            for (const mappedUser of data) {
              if (mappedUser.userId === user.id) {
                toAdd.isShared = true
                toAdd.isOwner = mappedUser.isOwner
              }
            }
            arr.push(toAdd)
          }
          this.models = sortBy(arr, ['isOwner', 'name'])
        })
      })
      .catch((err) => console.log(err))
  },
  methods: {
    changeShareStatus(model) {
      var url = `/api/vehicles/${this.vehicle.id}/users/${model.id}`
      if (model.isShared) {
        axios.post(url, {}).then((data) => {})
      } else {
        axios.delete(url).then((data) => {})
      }
    },
  },
}
</script>

<template>
  <div class="box">
    <h1 class="subtitle">Share {{ vehicle.nickname }}</h1>
    <section>
      <b-field v-for="model in models" :key="model.id">
        <b-switch v-model="model.isShared" :disabled="model.isOwner" @input="changeShareStatus(model)">
          {{ model.name }}
        </b-switch>
      </b-field>
    </section>
  </div>
</template>
