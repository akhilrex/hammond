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
          'We have successfully migrated the data from Clarkson. You will be redirected to the login screen shortly where you can login using your existing email and password : hammond'
      }
      if (this.migrationMode === 'fresh') {
        message =
          'You have been registered successfully. You will be redirected to the login screen shortly where you can login and start using the system.'
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
      <h1 class="title">Migrate from Clarkson</h1>
      <p>
        If you have an existing Clarkson deployment and you want to migrate your data from that, press the following button.
      </p>
      <br />
      <b-field> <b-button type="is-primary" @click="migrationMode = 'clarkson'">Migrate from Clarkson</b-button></b-field>
    </div>
    <div v-if="!migrationMode" class="box">
      <h1 class="title">Fresh Install</h1>
      <p>
        If you want a fresh install of Hammond, press the following button.
      </p>
      <br />
      <b-field>
        <b-button type="is-primary" @click="migrationMode = 'fresh'">Fresh Install</b-button>
      </b-field>
    </div>
    <div v-if="migrationMode === 'clarkson'" class="box content">
      <h1 class="title">Migrate from Clarkson</h1>
      <p>You need to make sure that this deployment of Hammond can access the MySQL database used by Clarkson.</p>
      <p>If that is not directly possible, you can make a copy of that database somewhere accessible from this instance.</p>
      <p>Once that is done, enter the connection string to the MySQL instance in the following format.</p>
      <p
        >All the users imported from Clarkson will have their username as their email in Clarkson database and pasword set to
        <span class="" style="font-weight:bold">hammond</span></p
      >
      <code>
        user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
      </code>
      <br />
      <br />
      <b-notification v-if="connectionError" type="is-danger" role="alert" :closable="false">
        {{ connectionError }}
      </b-notification>

      <b-field addons label="Mysql Connection String">
        <b-input v-model="url" required></b-input>
      </b-field>

      <div class="buttons">
        <b-button v-if="!testSuccess" type="is-primary" :disabled="isWorking" @click="testConnection">Test Connection</b-button
        ><b-button v-if="testSuccess" type="is-success" :disabled="isWorking" @click="migrate">Migrate</b-button>
        <b-button type="is-danger is-light" @click="resetMigrationMode">Cancel</b-button>
      </div>
    </div>
    <div v-if="migrationMode === 'fresh'" class="box content">
      <h1 class="title">Setup Admin Users</h1>
      <form @submit.prevent="register">
        <b-field label="Your Name">
          <b-input v-model="registerModel.name" required></b-input>
        </b-field>
        <b-field label="Your Email">
          <b-input v-model="registerModel.email" type="email" required></b-input>
        </b-field>
        <b-field label="Your Password">
          <b-input v-model="registerModel.password" type="password" required minlength="8" password-reveal></b-input>
        </b-field>
        <b-field label="Currency">
          <b-autocomplete
            v-model="registerModel.currency"
            :custom-formatter="formatCurrency"
            placeholder="Currency"
            :data="filteredCurrencyMasters"
            :keep-first="true"
            :open-on-focus="true"
            required
            @select="(option) => (selected = option)"
          ></b-autocomplete>
        </b-field>
        <b-field label="Distance Unit">
          <b-select v-model.number="registerModel.distanceUnit" placeholder="Distance Unit" required expanded>
            <option v-for="(option, key) in distanceUnitMasters" :key="key" :value="key">
              {{ `${option.long} (${option.short})` }}
            </option>
          </b-select>
        </b-field>
        <br />
        <div class="buttons">
          <b-button type="is-primary" native-type="submit" tag="input"></b-button>

          <b-button type="is-danger is-light" @click="resetMigrationMode">Cancel</b-button>
        </div>
      </form>
    </div>
  </Layout>
</template>
