<script>
import Layout from '@layouts/main.vue'
import { mapState } from 'vuex'
import store from '@state/store'

export default {
  components: { Layout },
  page: {
    title: 'Site Settings',
  },
  props: {
    user: {
      type: Object,
      required: true,
    },
    settings: {
      type: Object,
      required: true,
    },
  },
  data: function() {
    return {
      settingsModel: {
        currency: this.settings.currency,
        distanceUnit: this.settings.distanceUnit,
      },
      tryingToSave: false,
    }
  },
  computed: {
    ...mapState('vehicles', ['currencyMasters', 'distanceUnitMasters']),
  },
  methods: {
    saveSettings() {
      this.tryingToSave = true
      store
        .dispatch(`utils/saveSettings`, { settings: this.settingsModel })
        .then((data) => {
          this.$buefy.toast.open({
            message: 'Settings saved successfully',
            type: 'is-success',
            duration: 3000,
          })
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
          this.tryingToSave = false
        })
    },
  },
}
</script>

<template>
  <Layout>
    <div class="">
      <div class="columns">
        <div class="column">
          <h1 class="title">Site Settings</h1>
          <h1 class="subtitle">
            Update site level settings. These will be used as default values for new users.
          </h1>
        </div>
      </div>
    </div>
    <br />
    <form class="" @submit.prevent="saveSettings">
      <b-field label="Currency">
        <b-select v-model="settingsModel.currency" placeholder="Currency" required expanded>
          <option v-for="option in currencyMasters" :key="option.code" :value="option.code">
            {{ `${option.namePlural} (${option.code})` }}
          </option>
        </b-select>
      </b-field>
      <b-field label="Distance Unit">
        <b-select v-model.number="settingsModel.distanceUnit" placeholder="Distance Unit" required expanded>
          <option v-for="(option, key) in distanceUnitMasters" :key="key" :value="key">
            {{ `${option.long} (${option.short})` }}
          </option>
        </b-select>
      </b-field>
      <br />
      <b-field>
        <b-button tag="input" native-type="submit" :disabled="tryingToSave" type="is-primary" value="Save" expanded> </b-button>
      </b-field>
    </form>
  </Layout>
</template>
