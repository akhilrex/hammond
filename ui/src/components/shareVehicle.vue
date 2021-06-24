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
    transferVehicle(model) {
      if (!model.isShared) {
        return
      }
      this.$buefy.dialog.confirm({
        title: 'Transfer Vehicle',
        message: 'Are you sure you want to do this? You will lose ownership and all editing rights if you confirm.',
        cancelText: 'Cancel',
        confirmText: 'Go Ahead',
        onConfirm: () => {
          var url = `/api/vehicles/${this.vehicle.id}/users/${model.id}/transfer`
          axios
            .post(url, {})
            .then((data) => {
              this.$buefy.toast.open({
                message: 'Vehicle Transferred Successfully',
                type: 'is-success',
                duration: 3000,
              })
              setTimeout(() => {
              this.$router.go()
              }, 3000);
            })
            .catch((ex) => {
              this.$buefy.toast.open({
                duration: 5000,
                message: ex.message,
                position: 'is-bottom',
                type: 'is-danger',
              })
            })
        },
      })
    },
  },
}
</script>

<template>
  <div class="box" style="max-width:600px">
    <h1 class="subtitle">Share {{ vehicle.nickname }}</h1>
    <section>
      <div class="columns is-mobile" v-for="model in models" :key="model.id">
        <div class="column is-one-third">
          <b-field>
            <b-switch v-model="model.isShared" :disabled="model.isOwner" @input="changeShareStatus(model)">
              {{ model.name }}
            </b-switch>
          </b-field> </div
        ><div class="column is-three-quarters">
          <b-field>
            <b-button v-if="model.isShared && !model.isOwner" type="is-primary is-small" @click="transferVehicle(model)">Make Owner</b-button>
          </b-field></div
        ></div
      >
    </section>
  </div>
</template>
