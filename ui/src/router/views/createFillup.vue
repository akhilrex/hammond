<script>
import Layout from '@layouts/main.vue'
import QuickEntryDisplay from '@components/quickEntryDisplay.vue'
import store from '@state/store'

import { mapState } from 'vuex'
import axios from 'axios'
import { round } from 'lodash'

export default {
  page: {
    title: 'Create Fillup',
  },
  components: { Layout, QuickEntryDisplay },
  props: {
    vehicle: {
      type: Object,
      required: true,
    },
    fillup: {
      type: Object,
      required: false,
      default: function() {
        return {}
      },
    },
  },
  data() {
    return {
      authError: null,
      tryingToCreate: false,
      showMore: false,
      quickEntry: null,
      myVehicles: [],
      selectedVehicle: this.vehicle,
      fillupModel: this.fillup,
      processQuickEntry: false,
    }
  },

  computed: {
    user() {
      return store.state.auth.currentUser
    },
    ...mapState('users', ['me']),
    ...mapState('vehicles', ['fuelUnitMasters', 'fuelTypeMasters', 'vehicles']),
  },
  watch: {
    'fillupModel.fuelQuantity': function(old, newOne) {
      this.fillupModel.totalAmount = round(
        this.fillupModel.fuelQuantity * this.fillupModel.perUnitPrice,
        2
      )
    },
    'fillupModel.perUnitPrice': function(old, newOne) {
      this.fillupModel.totalAmount = round(
        this.fillupModel.fuelQuantity * this.fillupModel.perUnitPrice,
        2
      )
    },
    quickEntry: function(newOne, old) {
      if (old == null && newOne !== null) {
        this.processQuickEntry = true
      }
    },
  },
  mounted() {
    this.myVehicles = this.vehicles
    this.selectedVehicle = this.vehicle
    if (!this.fillup.id) {
      this.fillupModel = this.getEmptyFillup()
    }
  },
  methods: {
    getEmptyFillup() {
      return {
        vehicleId: this.selectedVehicle.id,
        fuelUnit: this.selectedVehicle.fuelUnit,
        perUnitPrice: null,
        fuelQuantity: null,
        totalAmount: null,
        odoReading: '',
        isTankFull: true,
        hasMissedFillup: false,
        date: new Date(),
        fillingStation: '',
        comments: '',
      }
    },
    async createFillup() {
      this.tryingToCreate = true
      this.fillupModel.vehicleId = this.selectedVehicle.id
      this.fillupModel.userId = this.me.id
      if (this.fillup.id) {
        axios
          .put(
            `/api/vehicles/${this.selectedVehicle.id}/fillups/${this.fillup.id}`,
            this.fillupModel
          )
          .then((data) => {
            this.$buefy.toast.open({
              message: 'Fillup Updated Successfully',
              type: 'is-success',
              duration: 3000,
            })
            this.fillupModel = this.getEmptyFillup()
            if (this.processQuickEntry) {
              store
                .dispatch('vehicles/setQuickEntryAsProcessed', { id: this.quickEntry.id })
                .then((data) => {
                  this.quickEntry = null
                })
            }
          })
          .catch((ex) => {
            this.$buefy.toast.open({
              duration: 5000,
              message: ex,
              position: 'is-bottom',
              type: 'is-danger',
            })
          })
          .finally(() => {
            this.tryingToCreate = false
          })
      } else {
        axios
          .post(`/api/vehicles/${this.selectedVehicle.id}/fillups`, this.fillupModel)
          .then((data) => {
            this.$buefy.toast.open({
              message: 'Fillup Created Successfully',
              type: 'is-success',
              duration: 3000,
            })
            this.fillupModel = this.getEmptyFillup()
            if (this.processQuickEntry) {
              store
                .dispatch('vehicles/setQuickEntryAsProcessed', { id: this.quickEntry.id })
                .then((data) => {})
            }
          })
          .catch((ex) => {
            this.$buefy.toast.open({
              duration: 5000,
              message: ex,
              position: 'is-bottom',
              type: 'is-danger',
            })
          })
          .finally(() => {
            this.tryingToCreate = false
          })
      }
    },
  },
}
</script>

<template>
  <Layout>
    <div class="has-text-centered">
      <div class="columns">
        <div class="column is-two-thirds">
          <h1 class="title">Create Fillup</h1>
          <h1 class="subtitle">
            {{
              [
                selectedVehicle.nickname,
                selectedVehicle.registration,
                selectedVehicle.make,
                selectedVehicle.model,
              ].join(' | ')
            }}
          </h1>
        </div>
        <div class="column is-one-thirds">
          <QuickEntryDisplay v-model="quickEntry" :user="user" />
        </div>
      </div>
    </div>
    <form class="" @submit.prevent="createFillup">
      <b-field label="Select a vehicle">
        <b-select
          v-model="selectedVehicle"
          placeholder="Vehicle"
          required
          expanded
          :disabled="fillup.id"
        >
          <option v-for="option in myVehicles" :key="option.id" :value="option">
            {{ option.nickname }}
          </option>
        </b-select>
      </b-field>
      <b-field label="Fillup Date">
        <b-datepicker
          v-model="fillupModel.date"
          placeholder="Click to select..."
          icon="calendar"
          trap-focus
          :max-date="new Date()"
        >
        </b-datepicker>
      </b-field>
      <b-field label="Quantity*" addons>
        <b-input
          v-model.number="fillupModel.fuelQuantity"
          type="number"
          step=".001"
          min="0"
          expanded
          required
        ></b-input>
        <b-select v-model="fillupModel.fuelUnit" placeholder="Fuel Unit" required>
          <option v-for="(option, key) in fuelUnitMasters" :key="key" :value="key">
            {{ option.long }}
          </option>
        </b-select>
      </b-field>
      <b-field :label="'Price per ' + vehicle.fuelUnitDetail.short + '*'"
        ><p class="control">
          <span class="button is-static">{{ me.currency }}</span>
        </p>
        <b-input
          v-model.number="fillupModel.perUnitPrice"
          type="number"
          min="0"
          step=".001"
          expanded
          required
        ></b-input>
      </b-field>
      <b-field label="Total Amount Paid">
        <p class="control">
          <span class="button is-static">{{ me.currency }}</span>
        </p>
        <b-input
          v-model.number="fillupModel.totalAmount"
          type="number"
          min="0"
          step=".001"
          expanded
          required
        ></b-input>
      </b-field>
      <b-field label="Odometer Reading">
        <p class="control">
          <span class="button is-static">{{ me.distanceUnitDetail.short }}</span>
        </p>
        <b-input
          v-model.number="fillupModel.odoReading"
          type="number"
          min="0"
          expanded
          required
        ></b-input>
      </b-field>
      <b-field>
        <b-checkbox v-model="fillupModel.isTankFull">Did you get a full tank?</b-checkbox>
      </b-field>
      <b-field>
        <b-checkbox v-model="fillupModel.hasMissedFillup"
          >Did you miss the fillup entry before this one?</b-checkbox
        >
      </b-field>
      <b-field>
        <b-switch v-model="showMore">Fill more details</b-switch>
      </b-field>
      <fieldset v-if="showMore">
        <b-field label="Filling Station Name">
          <b-input v-model="fillupModel.fillingStation" type="text" expanded></b-input>
        </b-field>
        <b-field label="Comments">
          <b-input v-model="fillupModel.comments" type="textarea" expanded></b-input>
        </b-field>
      </fieldset>
      <b-field>
        <b-switch v-if="quickEntry" v-model="processQuickEntry"
          >Mark selected Quick Entry as processed</b-switch
        >
      </b-field>
      <br />
      <b-field>
        <b-button
          tag="input"
          native-type="submit"
          :disabled="tryingToCreate"
          type="is-primary"
          label="Create Fillup"
          expanded
        >
        </b-button>
        <p v-if="authError">
          There was an error logging in to your account.
        </p>
      </b-field>
    </form>
  </Layout>
</template>
