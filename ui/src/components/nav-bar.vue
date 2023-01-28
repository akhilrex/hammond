<script>
import { authComputed } from '@state/helpers'
import { mapGetters } from 'vuex'
import NavBarRoutes from './nav-bar-routes.vue'

export default {
  components: { NavBarRoutes },
  data() {
    return {
      persistentNavRoutes: [
        {
          name: 'home',
          title: this.$t('menu.home'),
        },
      ],
      loggedInNavRoutes: [
        {
          name: 'quickEntries',
          title: () => this.$t('menu.quickentries'),
          badge: () => this.unprocessedQuickEntries.length,
        },
        {
          name: 'import',
          title: () => this.$t('menu.import'),
        },
        {
          name: 'settings',
          title: this.$t('menu.settings'),
        },
        {
          name: 'logout',
          title: this.$t('menu.logout'),
        },
      ],
      loggedOutNavRoutes: [
        {
          name: 'login',
          title: this.$t('menu.login'),
        },
      ],
      adminNavRoutes: [
        {
          name: 'site-settings',
          title: this.$t('menu.sitesettings'),
        },
        {
          name: 'users',
          title: this.$t('menu.users'),
        },
      ],
    }
  },
  computed: {
    ...authComputed,
    ...mapGetters('vehicles', ['unprocessedQuickEntries']),
    isAdmin() {
      return this.loggedIn && this.currentUser.role === 'ADMIN'
    },
  },
}
</script>

<template>
  <div class="container">
    <b-navbar class="" spaced>
      <template v-slot:brand>
        <b-navbar-item tag="router-link" :to="{ path: '/' }">
          <h1 class="title" style="font-family:Arial Black">Hammond</h1>
        </b-navbar-item>
      </template>
      <template v-slot:end>
        <NavBarRoutes :routes="persistentNavRoutes" />
        <NavBarRoutes v-if="loggedIn" :routes="loggedInNavRoutes" />
        <NavBarRoutes v-else :routes="loggedOutNavRoutes" />
        <b-navbar-dropdown v-if="loggedIn && isAdmin" :label="$t('menu.admin')">
          <NavBarRoutes :routes="adminNavRoutes" />
        </b-navbar-dropdown>
      </template>
    </b-navbar>
  </div>
</template>
