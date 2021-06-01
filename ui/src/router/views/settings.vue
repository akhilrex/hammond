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
      },
      tryingToSave: false,
      changePassModel: {
        old: '',
        new: '',
        renew: '',
      },
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
          let errorMessage= ex.message;
          if(ex.response && ex.response.data?.errors?.changePassword){
            errorMessage=ex.response.data?.errors?.changePassword
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
    <h1 class="title">Your Settings</h1>
    <div class="columns"
      ><div class="column">
        <form class="box " @submit.prevent="saveSettings">
          <h1 class="subtitle">
            These will be used as default values whenever you create a new fillup or expense.
          </h1>
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
      </div>
      <div class="column">
        <form class="box" @submit.prevent="changePassword">
          <h1 class="subtitle">Change password</h1>
          <b-field label="Old Password">
            <b-input v-model="changePassModel.old" required minlength="6" password-reveal type="password"></b-input>
          </b-field>
          <b-field label="New Password">
            <b-input v-model="changePassModel.new" required minlength="6" password-reveal type="password"></b-input>
          </b-field>
          <b-field label="Repeat New Password">
            <b-input v-model="changePassModel.renew" required minlength="6" password-reveal type="password"></b-input>
          </b-field>
          <p v-if="!passwordValid" class="help is-danger">Password values don't match</p>
          <b-field>
            <b-button tag="input" native-type="submit" :disabled="!passwordValid" type="is-primary" value="Change Password" expanded> </b-button>
          </b-field>
        </form>
      </div>
    </div>
    <hr />
    <div class="columns">
      <div class="twelve">
        <h3 class="title">More Info</h3>
        <p style="font-style: italic;">
          This project is under active development which means I release new updates very frequently. I will eventually build the version
          management/update checking mechanism. Until then it is recommended that you use something like watchtower which will automatically update
          your containers whenever I release a new version or periodically rebuild the container with the latest image manually.
        </p>
        <br />
        <table class="table is-hoverable">
          <tr>
            <td>Current Version</td>
            <td>2021.06.01</td>
          </tr>
          <tr>
            <td>Website</td>
            <td><a href="https://github.com/akhilrex/hammond" target="_blank">https://github.com/akhilrex/hammond</a></td>
          </tr>
          <tr>
            <td>Found a bug</td>
            <td
              ><a
                href="https://github.com/akhilrex/hammond/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc"
                target="_blank"
                rel="noopener noreferrer"
                >Report here</a
              ></td
            >
          </tr>
          <tr>
            <td>Feature Request</td>
            <td
              ><a
                href="https://github.com/akhilrex/hammond/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc"
                target="_blank"
                rel="noopener noreferrer"
                >Request here</a
              ></td
            >
          </tr>
          <tr>
            <td>Support the developer</td>
            <td><a href="https://www.buymeacoffee.com/akhilrex" target="_blank" rel="noopener noreferrer">Support here</a></td>
          </tr>
        </table>
      </div>
    </div>
  </Layout>
</template>
