<script>
import Layout from '@layouts/main.vue'
import { authMethods } from '@state/helpers'
import appConfig from '@src/app.config'
import store from '@state/store'

export default {
  page: {
    title: 'Log in',
    meta: [{ name: 'description', content: `Log in to ${appConfig.title}` }],
  },
  components: { Layout },
  data() {
    return {
      username: '',
      password: '',
      authError: null,
      tryingToLogIn: false,
      errorMessage: '',
    }
  },
  computed: {
    placeholders() {
      return process.env.NODE_ENV === 'production'
        ? {}
        : {
            username: this.$t('enterusername'),
            password: this.$t('enterpassword'),
          }
    },
  },
  mounted() {
    // console.log(process.env.API_BASE_URL)
  },
  methods: {
    ...authMethods,
    // Try to log the user in with the username
    // and password they provided.
    tryToLogIn() {
      this.tryingToLogIn = true
      this.errorMessage = ''
      // Reset the authError if it existed.
      this.authError = null
      return this.logIn({
        username: this.username,
        password: this.password,
      })
        .then((token) => {
          this.tryingToLogIn = false
          store.dispatch('users/me').then((data) => {
            this.$router.push(this.$route.query.redirectFrom || { name: 'home' })
          })
          // Redirect to the originally requested page, or to the home page
        })
        .catch((error) => {
          if (error.response.data?.errors?.login) {
            this.errorMessage = error.response.data.errors.login
          }
          this.tryingToLogIn = false
          this.authError = error
        })
    },
  },
}
</script>

<template>
  <Layout>
    <form @submit.prevent="tryToLogIn">
      <b-field :label="$t('email')"> <b-input v-model="username" tag="b-input" name="username" type="email" :placeholder="placeholders.username"/></b-field>
      <b-field :label="$t('password')">
        <b-input v-model="password" tag="b-input" name="password" type="password" :placeholder="placeholders.password" />
      </b-field>
      <b-button tag="input" native-type="submit" :value="$t('login')" :disabled="tryingToLogIn" type="is-primary">
        <BaseIcon v-if="tryingToLogIn" name="sync" spin />
        <span v-else>
          {{ $t('login') }}
        </span>
      </b-button>
      <p v-if="authError"> {{ $t('loginerror', { msg: errorMessage }) }}</p>
    </form>
  </Layout>
</template>
