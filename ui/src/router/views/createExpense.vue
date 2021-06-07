<script>
import Layout from '@layouts/main.vue'
import QuickEntryDisplay from '@components/quickEntryDisplay.vue'
import { mapState } from 'vuex'
import axios from 'axios'
import store from '@state/store'

export default {
  page: {
    title: 'Create Expense',
  },
  components: { Layout, QuickEntryDisplay },
  props: {
    vehicle: {
      type: Object,
      required: true,
    },
    expense: {
      type: Object,
      required: false,
      default: function() {
        return {}
      },
    },
  },
  data() {
    return {
      tryingToCreate: false,
      showMore: false,
      quickEntry: null,
      myVehicles: [],
      users: [],
      selectedVehicle: this.vehicle,
      expenseModel: this.expense,
      processQuickEntry: false,
    }
  },

  computed: {
    user() {
      return store.state.auth.currentUser
    },
    ...mapState('utils', ['isMobile']),
    ...mapState('users', ['me']),
    ...mapState('vehicles', ['fuelUnitMasters', 'fuelTypeMasters', 'vehicles']),
  },
  watch: {
    quickEntry: function(newOne, old) {
      if (old == null && newOne !== null) {
        this.processQuickEntry = true
      }
    },
  },
  mounted() {
    this.myVehicles = this.vehicles
    this.selectedVehicle = this.vehicle
    this.fetchVehicleUsers()
    if (!this.expense.id) {
      this.expenseModel = this.getEmptyExpense()
      this.expenseModel.userId = this.me.id
    }
  },
  methods: {
    getEmptyExpense() {
      return {
        vehicleId: this.selectedVehicle.id,
        amount: null,
        expenseType: '',
        odoReading: '',
        date: new Date(),
        comments: '',
        userId: '',
      }
    },
    fetchVehicleUsers() {
      store
        .dispatch('vehicles/fetchUsersByVehicleId', { vehicleId: this.selectedVehicle.id })
        .then((data) => {
          this.users = data
        })
        .catch((err) => console.log(err))
    },
    createExpense() {
      this.tryingToCreate = true
      this.expenseModel.vehicleId = this.selectedVehicle.id
      //   this.expenseModel.userId = this.me.id

      if (this.expense.id) {
        axios
          .put(`/api/vehicles/${this.selectedVehicle.id}/expenses/${this.expense.id}`, this.expenseModel)
          .then((data) => {
            this.$buefy.toast.open({
              message: 'Expense Updated Successfully',
              type: 'is-success',
              duration: 3000,
            })
            this.expenseModel = this.getEmptyExpense()
            if (this.processQuickEntry) {
              store.dispatch('vehicles/setQuickEntryAsProcessed', { id: this.quickEntry.id }).then((data) => {})
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
          .post(`/api/vehicles/${this.selectedVehicle.id}/expenses`, this.expenseModel)
          .then((data) => {
            this.$buefy.toast.open({
              message: 'Expense Created Successfully',
              type: 'is-success',
              duration: 3000,
            })
            this.expenseModel = this.getEmptyExpense()
            if (this.processQuickEntry) {
              store.dispatch('vehicles/setQuickEntryAsProcessed', { id: this.quickEntry.id }).then((data) => {
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
      }
    },
  },
}
</script>

<template>
  <Layout>
    <div class="columns">
      <div class="column is-two-thirds">
        <h1 class="title">Create Expense</h1>
        <h1 class="subtitle">
          {{ [selectedVehicle.nickname, selectedVehicle.registration, selectedVehicle.make, selectedVehicle.model].join(' | ') }}
        </h1>
      </div>
      <div class="column is-one-thirds">
        <QuickEntryDisplay v-model="quickEntry" :user="user" />
      </div>
    </div>
    <form @submit.prevent="createExpense">
      <b-field label="Select a vehicle">
        <b-select v-model="selectedVehicle" placeholder="Vehicle" required expanded :disabled="expense.id">
          <option v-for="option in myVehicles" :key="option.id" :value="option">
            {{ option.nickname }}
          </option>
        </b-select>
      </b-field>
      <b-field label="Expense by">
        <b-select v-model="expenseModel.userId" placeholder="User" required expanded :disabled="expense.id">
          <option v-for="option in users" :key="option.userId" :value="option.userId">
            {{ option.name }}
          </option>
        </b-select>
      </b-field>
      <b-field label="Expense Date">
        <b-datepicker v-model="expenseModel.date" placeholder="Click to select..." icon="calendar" :max-date="new Date()"> </b-datepicker>
      </b-field>
      <b-field label="Expense Type*">
        <b-input v-model="expenseModel.expenseType" expanded required></b-input>
      </b-field>

      <b-field label="Total Amount Paid">
        <p class="control">
          <span class="button is-static">{{ me.currency }}</span>
        </p>
        <b-input v-model.number="expenseModel.amount" type="number" min="0" expanded step=".001" required></b-input>
      </b-field>
      <b-field label="Odometer Reading">
        <p class="control">
          <span class="button is-static">{{ me.distanceUnitDetail.short }}</span>
        </p>
        <b-input v-model.number="expenseModel.odoReading" type="number" min="0" expanded required></b-input>
      </b-field>

      <b-field>
        <b-switch v-model="showMore">Fill more details</b-switch>
      </b-field>
      <fieldset v-if="showMore">
        <b-field label="Comments">
          <b-input v-model="expenseModel.comments" type="textarea" expanded></b-input>
        </b-field>
      </fieldset>
      <b-field>
        <b-switch v-if="quickEntry" v-model="processQuickEntry">Mark selected Quick Entry as processed</b-switch>
      </b-field>
      <br />
      <b-field>
        <b-button tag="input" native-type="submit" :disabled="tryingToCreate" type="is-primary" label="Create Expense" expanded> </b-button>
      </b-field>
    </form>
  </Layout>
</template>
