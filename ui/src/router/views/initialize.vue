<script>
import Layout from '@layouts/main.vue'
import axios from 'axios'
import { mapGetters, mapState } from 'vuex'
import store from '@state/store'

export default {
  components: { Layout },
  page: {
    title: 'First Setup',
  },
  data: function() {
    return {
      migrationMode: '',
      url: '',
      testSuccess: false,
      connectionError: '',
      isWorking: false,
      registerModel: {
        name: '',
        email: '',
        password: '',
        distanceUnit: 1,
        currency: '',
      },
    }
  },
  computed: {
    ...mapGetters('auth', ['isInitialized']),
    ...mapState('vehicles', ['currencyMasters', 'distanceUnitMasters']),
    filteredCurrencyMasters() {
      return this.currencyMasters.filter((option) => {
        return (
          option.namePlural
            .toString()
            .toLowerCase()
            .indexOf(this.registerModel.currency.toLowerCase()) >= 0 ||
          option.code
            .toString()
            .toLowerCase()
            .indexOf(this.registerModel.currency.toLowerCase()) >= 0
        )
      })
    },
  },
  mounted() {
    store.dispatch('vehicles/fetchMasters').then((data) => {})
  },
  methods: {
    resetMigrationMode() {
      this.migrationMode = ''
      this.url = ''
      this.registerModel = {
        name: '',
        email: '',
        password: '',
        distanceUnit: 1,
        currency: 'INR',
      }
    },
    showSuccessModal() {
      var message = ''
      if (this.migrationMode === 'clarkson') {
        message =
          this.$t('init.clarkson.success')
      }
      if (this.migrationMode === 'fresh') {
        message =
          this.$t('init.fresh.success')
      }
      this.$buefy.toast.open({
        duration: 10000,
        message: message,
        position: 'is-bottom',
        type: 'is-success',
      })
      setTimeout(() => {
        this.$router
          .push({ name: 'login' })
          .then((succ) => {})
          .catch((err) => console.log('error:', err))
      }, 10000)
    },
    register() {
      this.isWorking = true
      axios
        .post('/api/auth/initialize', this.registerModel)
        .then((response) => {
          const success = response.data.success
          if (success) {
            store.dispatch('auth/systemInitialized').then((data) => {
              this.showSuccessModal()
            })
          }
        })
        .catch((ex) => {
          this.testSuccess = false
          this.$buefy.toast.open({
            duration: 5000,
            message: ex.message,
            position: 'is-bottom',
            type: 'is-danger',
          })
        })
        .finally(() => (this.isWorking = false))
    },
    testConnection() {
      if (!this.url) {
        return
      }
      this.isWorking = true
      axios
        .post('/api/clarkson/check', { url: this.url })
        .then((response) => {
          this.testSuccess = response.data.canMigrate
          if (!this.testSuccess) {
            this.connectionError = response.data.message
          } else {
            this.connectionError = ''
          }
        })
        .catch((ex) => {
          this.testSuccess = false
          this.$buefy.toast.open({
            duration: 5000,
            message: ex.message,
            position: 'is-bottom',
            type: 'is-danger',
          })
        })
        .finally(() => (this.isWorking = false))
    },
    migrate() {
      if (!this.url) {
        return
      }
      this.isWorking = true
      axios
        .post('/api/clarkson/migrate', { url: this.url })
        .then((data) => {
          store.dispatch('auth/systemInitialized').then((data) => {
            this.showSuccessModal()
          })
        })
        .catch((ex) => {
          this.testSuccess = false
          this.$buefy.toast.open({
            duration: 5000,
            message: ex.message,
            position: 'is-bottom',
            type: 'is-danger',
          })
        })
        .finally(() => (this.isWorking = false))
    },
    formatCurrency(option) {
      return `${option.namePlural} (${option.code})`
    },
  },
}
</script>

<template>
  <Layout>
    <div v-if="!migrationMode" class="box">
      <h1 class="title">{{ $t('init.migrateclarkson') }}</h1>
      <p>
        {{ $t('init.migrateclarksondesc') }}
      </p>
      <br />
      <b-field> <b-button type="is-primary" @click="migrationMode = 'clarkson'">{{ $t('init.migrateclarkson') }}</b-button></b-field>
    </div>
    <div v-if="!migrationMode" class="box">
      <h1 class="title">{{ $t('init.freshinstall') }}</h1>
      <p>
        {{ $t('init.freshinstalldesc') }}
      </p>
      <br />
      <b-field>
        <b-button type="is-primary" @click="migrationMode = 'fresh'">{{ $t('init.freshinstall') }}</b-button>
      </b-field>
    </div>
    <div v-if="migrationMode === 'clarkson'" class="box content">
      <h1 class="title">{{ $t('init.migrateclarkson') }}</h1>
      <p v-html="$t('init.clarkson.desc')"></p>
      <b-notification v-if="connectionError" type="is-danger" role="alert" :closable="false">
        {{ connectionError }}
      </b-notification>

      <b-field addons :label="this.$t('mysqlconnstr')">
        <b-input v-model="url" required></b-input>
      </b-field>

      <div class="buttons">
        <b-button v-if="!testSuccess" type="is-primary" :disabled="isWorking" @click="testConnection">{{ $t('testconn') }}</b-button>
        <b-button v-if="testSuccess" type="is-success" :disabled="isWorking" @click="migrate">{{ $t('migrate') }}</b-button>
        <b-button type="is-danger is-light" @click="resetMigrationMode">{{ $t('cancel') }}</b-button>
      </div>
    </div>
    <div v-if="migrationMode === 'fresh'" class="box content">
      <h1 class="title">{{ $t('init.fresh.setupadminuser') }}</h1>
      <form @submit.prevent="register">
        <b-field :label="this.$t('init.fresh.yourname')">
          <b-input v-model="registerModel.name" required></b-input>
        </b-field>
        <b-field :label="this.$t('init.fresh.youremail')">
          <b-input v-model="registerModel.email" type="email" required></b-input>
        </b-field>
        <b-field :label="this.$t('init.fresh.yourpassword')">
          <b-input v-model="registerModel.password" type="password" required minlength="8" password-reveal></b-input>
        </b-field>
        <b-field :label="this.$t('currency')">
          <b-autocomplete
            v-model="registerModel.currency"
            :custom-formatter="formatCurrency"
            :placeholder="this.$t('currency')"
            :data="filteredCurrencyMasters"
            :keep-first="true"
            :open-on-focus="true"
            required
            @select="(option) => (selected = option)"
          ></b-autocomplete>
        </b-field>
        <b-field :label="this.$t('distanceunit')">
          <b-select v-model.number="registerModel.distanceUnit" :placeholder="this.$t('distanceunit')" required expanded>
            <option v-for="(option, key) in distanceUnitMasters" :key="key" :value="key">
              {{ `${$t('unit.long.' + option.key)} (${$t('unit.short.' + option.key)})` }}
            </option>
          </b-select>
        </b-field>
        <br />
        <div class="buttons">
          <b-button type="is-primary" native-type="submit" tag="input" :value="this.$t('save')"></b-button>

          <b-button type="is-danger is-light" @click="resetMigrationMode">{{ $t('cancel') }}</b-button>
        </div>
      </form>
    </div>
  </Layout>
</template>
