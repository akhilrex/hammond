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
          title: 'Home',
        },
      ],
      loggedInNavRoutes: [
        {
          name: 'quickEntries',
          title: () => 'Quick Entries',
          badge: () => this.unprocessedQuickEntries.length,
        },
        {
          name: 'import',
          title: () => 'Import',
        },
        {
          name: 'settings',
          title: 'Settings',
        },
        {
          name: 'logout',
          title: 'Log out',
        },
      ],
      loggedOutNavRoutes: [
        {
          name: 'login',
          title: 'Log in',
        },
      ],
      adminNavRoutes: [
        {
          name: 'site-settings',
          title: 'Site Settings',
        },
        {
          name: 'users',
          title: 'Users',
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
        <b-navbar-dropdown v-if="loggedIn && isAdmin" label="Admin">
          <NavBarRoutes :routes="adminNavRoutes" />
        </b-navbar-dropdown>
      </template>
    </b-navbar>
  </div>
</template>
