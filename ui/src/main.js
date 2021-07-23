import Vue from 'vue'
import Buefy from 'buefy'
import router from '@router'
import store from '@state/store'
import { library } from '@fortawesome/fontawesome-svg-core'
import {
  faCheck,
  faTimes,
  faArrowUp,
  faAngleLeft,
  faAngleRight,
  faCalendar,
  faEdit,
  faAngleDown,
  faAngleUp,
  faUpload,
  faExclamationCircle,
  faDownload,
  faEye,
  faEyeSlash,
  faTrash,
  faShare,
  faUserFriends,
  faTimesCircle,
} from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

import App from './app.vue'

// Globally register all `_base`-prefixed components
import '@components/_globals'

import 'buefy/dist/buefy.css'
import 'nprogress/nprogress.css'

Vue.component('vue-fontawesome', FontAwesomeIcon)
library.add(
  faCheck,
  faTimes,
  faArrowUp,
  faAngleLeft,
  faAngleRight,
  faCalendar,
  faEdit,
  faAngleDown,
  faAngleUp,
  faUpload,
  faExclamationCircle,
  faDownload,
  faEye,
  faEyeSlash,
  faTrash,
  faShare,
  faUserFriends,
  faTimesCircle
)
Vue.use(Buefy, {
  defaultIconComponent: 'vue-fontawesome',
  defaultIconPack: 'fas',
})

// Don't warn about using the dev version of Vue in development.
Vue.config.productionTip = process.env.NODE_ENV === 'production'

// If running inside Cypress...
if (process.env.VUE_APP_TEST === 'e2e') {
  // Ensure tests fail when Vue emits an error.
  Vue.config.errorHandler = window.Cypress.cy.onUncaughtException
}

const app = new Vue({
  router,
  store,

  render: (h) => h(App),
}).$mount('#app')

// If running e2e tests...
if (process.env.VUE_APP_TEST === 'e2e') {
  // Attach the app to the window, which can be useful
  // for manually setting state in Cypress commands
  // such as `cy.logIn()`.
  window.__app__ = app
}
