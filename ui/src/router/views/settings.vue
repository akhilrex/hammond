<script>
import Layout from '@layouts/main.vue'
import { mapState } from 'vuex'
import store from '@state/store'
import axios from 'axios'

export default {
  page: {
    title: 'Settings',
  },
  components: { Layout },
  props: {
    user: {
      type: Object,
      required: true,
    },
    me: {
      type: Object,
      required: true,
    },
  },
  data: function() {
    return {
      settingsModel: {
        currency: this.me.currency,
        distanceUnit: this.me.distanceUnit,
        dateFormat: this.me.dateFormat,
      },
      tryingToSave: false,
      changePassModel: {
        old: '',
        new: '',
        renew: '',
      },
      dateFormatMasters: ['dd/MM/yyyy', 'MM/dd/yyyy', 'yyyy/MM/dd'],
    }
  },
  computed: {
    ...mapState('vehicles', ['currencyMasters', 'distanceUnitMasters']),
    passwordValid() {
      if (this.changePassModel.new === '' || this.changePassModel.renew === '') {
        return true
      }

      return this.changePassModel.new === this.changePassModel.renew
    },
    filteredCurrencyMasters() {
      return this.currencyMasters.filter((option) => {
        return (
          option.namePlural
            .toString()
            .toLowerCase()
            .indexOf(this.settingsModel.currency.toLowerCase()) >= 0 ||
          option.code
            .toString()
            .toLowerCase()
            .indexOf(this.settingsModel.currency.toLowerCase()) >= 0
        )
      })
    },
  },
  methods: {
    changePassword() {
      if (!this.passwordValid) {
        return
      }
      this.tryingToSavePass = true

      axios
        .post('/api/changePassword', {
          oldPassword: this.changePassModel.old,
          newPassword: this.changePassModel.new,
        })
        .then((data) => {
          this.$buefy.toast.open({
            message: 'Password changed successfully. You will be logged out now.',
            type: 'is-success',
            duration: 3000,
          })
          this.$router
            .push({ name: 'logout' })
            .then((succ) => {
              console.log(succ)
            })
            .catch((err) => console.log('error:', err))
        })
        .catch((ex) => {
          let errorMessage = ex.message
          if (ex.response && ex.response.data?.errors?.changePassword) {
            errorMessage = ex.response.data?.errors?.changePassword
          }
          this.$buefy.toast.open({
            duration: 5000,
            message: errorMessage,
            position: 'is-bottom',
            type: 'is-danger',
          })
        })
        .finally(() => {
          this.tryingToSavePass = false
        })
    },
    saveSettings() {
      this.tryingToSave = true
      store
        .dispatch(`utils/saveUserSettings`, { settings: this.settingsModel })
        .then((data) => {
          this.$buefy.toast.open({
            message: this.$t('settingssaved'),
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
    formatCurrency(option) {
      return `${option.namePlural} (${option.code})`
    },
  },
}
</script>

<template>
  <Layout>
    <h1 class="title">{{ $t('yoursettings') }}</h1>
    <div class="columns"
      ><div class="column">
        <form class="box " @submit.prevent="saveSettings">
          <h1 class="subtitle">
            {{ $t('settingdesc') }}
          </h1>
          <b-field :label="$t('currency')">
            <b-autocomplete
              v-model="settingsModel.currency"
              :custom-formatter="formatCurrency"
              :placeholder="$t('currency')"
              :data="filteredCurrencyMasters"
              :keep-first="true"
              :open-on-focus="true"
              required
              @select="(option) => (selected = option)"
            ></b-autocomplete>
          </b-field>
          <b-field :label="$t('distanceunit')">
            <b-select v-model.number="settingsModel.distanceUnit" placeholder="Distance Unit" required expanded>
              <option v-for="(option, key) in distanceUnitMasters" :key="key" :value="key">
                {{ `${$t('unit.long.' + option.key)} (${$t('unit.short.' + option.key)})` }}
              </option>
            </b-select>
          </b-field>
          <b-field :label="$t('dateformat')">
            <b-select v-model.number="settingsModel.dateFormat" placeholder="Date Format" required expanded>
              <option v-for="option in dateFormatMasters" :key="option" :value="option">
                {{ `${option}` }}
              </option>
            </b-select>
          </b-field>
          <br />
          <b-field>
            <b-button tag="input" native-type="submit" :disabled="tryingToSave" type="is-primary" :value="$t('save')" expanded> </b-button>
          </b-field>
        </form>
      </div>
      <div class="column">
        <form class="box" @submit.prevent="changePassword">
          <h1 class="subtitle">{{ $t('changepassword') }}</h1>
          <b-field :label="$t('oldpassword')">
            <b-input v-model="changePassModel.old" required minlength="6" password-reveal type="password"></b-input>
          </b-field>
          <b-field :label="$t('newpassword')">
            <b-input v-model="changePassModel.new" required minlength="6" password-reveal type="password"></b-input>
          </b-field>
          <b-field :label="$t('repeatnewpassword')">
            <b-input v-model="changePassModel.renew" required minlength="6" password-reveal type="password"></b-input>
          </b-field>
          <p v-if="!passwordValid" class="help is-danger">{{ $t('passworddontmatch') }}</p>
          <b-field>
            <b-button tag="input" native-type="submit" :disabled="!passwordValid" type="is-primary" :value="$t('changepassword')" expanded>
            </b-button>
          </b-field>
        </form>
      </div>
    </div>
    <hr />
    <div class="columns">
      <div class="twelve">
        <h3 class="title">{{ $t('moreinfo') }}</h3>
        <table class="table is-hoverable">
          <tr>
            <td>{{ $t('currentversion') }}</td>
            <td>2022.07.06</td>
          </tr>
          <tr>
            <td>Website</td>
            <td><a href="https://github.com/alfhou/hammond" target="_blank">https://github.com/alfhou/hammond</a></td>
          </tr>
          <tr>
            <td>{{ $t('foundabug') }}</td>
            <td
              ><a href="https://github.com/alfhou/hammond/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc" target="_blank" rel="noopener noreferrer"
                >Report here</a
              ></td
            >
          </tr>
          <tr>
            <td>{{ $t('featurerequest') }}</td>
            <td
              ><a href="https://github.com/alfhou/hammond/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc" target="_blank" rel="noopener noreferrer"
                >Request here</a
              ></td
            >
          </tr>
        </table>
      </div>
    </div>
  </Layout>
</template>
