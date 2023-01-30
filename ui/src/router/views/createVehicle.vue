<script>
import Layout from '@layouts/main.vue'
import { mapState } from 'vuex'
import axios from 'axios'

export default {
  page: {
    title: 'Create Vehicle',
  },
  components: { Layout },
  props: {
    vehicle: {
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
      myVehicles: [],
      vehicleModel: {},
    }
  },

  computed: {
    ...mapState('users', ['me']),
    ...mapState('vehicles', ['fuelUnitMasters', 'fuelTypeMasters', 'vehicles']),
  },
  watch: {},
  mounted() {
    if (!this.vehicle) {
      this.vehicleModel = this.getEmptyVehicle()
    } else {
      this.vehicleModel = this.getEmptyVehicle(this.vehicle)
    }
    this.myVehicles = this.vehicles
  },
  methods: {
    getEmptyVehicle(veh) {
      if (!veh || !veh.id) {
        return {
          fuelUnit: null,
          fuelType: null,
          registration: '',
          vin: '',
          nickname: '',
          engineSize: null,
          make: '',
          model: '',
          yearOfManufacture: null,
        }
      } else {
        return {
          fuelUnit: veh.fuelUnit,
          fuelType: veh.fuelType,
          registration: veh.registration,
          vin: veh.vin,
          nickname: veh.nickname,
          engineSize: veh.engineSize,
          make: veh.make,
          model: veh.model,
          yearOfManufacture: veh.yearOfManufacture,
        }
      }
    },
    createVehicle() {
      this.tryingToCreate = true
      this.vehicleModel.userId = this.me.id
      if (this.vehicle.id) {
        axios
          .put(`/api/vehicles/${this.vehicle.id}`, this.vehicleModel)
          .then((data) => {
            this.$buefy.toast.open({
              message: this.$t('vehiclesavedsuccessfully'),
              type: 'is-success',
              duration: 3000,
            })
            // this.vehicleModel = this.getEmptyVehicle()
          })
          .catch((ex) => {
            this.$buefy.toast.open({
              duration: 5000,
              message: ex.message,
              position: 'is-bottom',
              type: 'is-danger',
            })
          })
          .finally(() => {
            this.tryingToCreate = false
          })
      } else {
        axios
          .post(`/api/vehicles`, this.vehicleModel)
          .then((data) => {
            this.$buefy.toast.open({
              message: this.$t('vehiclesavedsuccessfully'),
              type: 'is-success',
              duration: 3000,
            })
            this.vehicleModel = this.getEmptyVehicle()
          })
          .catch((ex) => {
            this.$buefy.toast.open({
              duration: 5000,
              message: ex.message,
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
      <div class="column is-three-quarters">
        <h1 class="title">{{ $t('createvehicle') }}</h1>
      </div>
      <div class="column is-one-quarter">
        <router-link tag="b-button" type="is-primary" to="/">
          {{ $t('back') }}
        </router-link>
      </div>
    </div>
    <form @submit.prevent="createVehicle">
      <b-field :label="this.$t('nickname') + `*`">
        <b-input v-model="vehicleModel.nickname" type="text" expanded required></b-input>
      </b-field>
      <b-field :label="this.$t('registration') + `*`">
        <b-input v-model="vehicleModel.registration" type="text" expanded required></b-input>
      </b-field>
      <b-field label="VIN">
        <b-input v-model="vehicleModel.vin" type="text" expanded></b-input>
      </b-field>
      <b-field :label="this.$t('fueltype') + `*`">
        <b-select v-model.number="vehicleModel.fuelType" :placeholder="this.$t('fueltype')" required expanded>
          <option v-for="(option, key) in fuelTypeMasters" :key="key" :value="key">
            {{ $t('fuel.' + option.key) }}
          </option>
        </b-select>
      </b-field>

      <b-field :label="this.$t('fuelunit') + `*`">
        <b-select v-model.number="vehicleModel.fuelUnit" :placeholder="this.$t('fuelunit')" required expanded>
          <option v-for="(option, key) in fuelUnitMasters" :key="key" :value="key">
            {{ $t('unit.long.' + option.key) }}
          </option>
        </b-select>
      </b-field>

      <b-field :label="this.$t('make') + `*`">
        <b-input v-model="vehicleModel.make" type="text" required expanded></b-input>
      </b-field>
      <b-field :label="this.$t('model') + `*`">
        <b-input v-model="vehicleModel.model" type="text" required expanded></b-input>
      </b-field>
      <b-field :label="this.$t('yearmanufacture') + `*`">
        <b-input v-model.number="vehicleModel.yearOfManufacture" type="number" expanded number></b-input>
      </b-field>
      <b-field :label="this.$t('enginesize')">
        <b-input v-model.number="vehicleModel.engineSize" type="number" expanded number></b-input>
      </b-field>

      <br />
      <b-field>
        <b-button
          tag="input"
          native-type="submit"
          :disabled="tryingToCreate"
          type="is-primary"
          :value="this.$t('save')"
          :label="this.$t('createvehicle')"
          expanded
        >
          <BaseIcon v-if="tryingToCreate" name="sync" spin />
        </b-button>
        <p v-if="authError">
          {{ $t('loginerror') }}
        </p>
      </b-field>
    </form>
  </Layout>
</template>
