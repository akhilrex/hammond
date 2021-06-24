<script>
import Layout from '@layouts/main.vue'
import store from '@state/store'
import { mapState } from 'vuex'
import axios from 'axios'

import { parseAndFormatDate } from '@utils/format-date'

export default {
  components: { Layout },
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
      users: [],
      showUserForm: false,
      isWorking: false,
      registerModel: {
        name: '',
        email: '',
        password: '',
        distanceUnit: this.settings.distanceUnit,
        currency: this.settings.currency,
        role: 1,
      },
    }
  },
  page() {
    return {
      title: 'User Management ',
    }
  },
  computed: {
    ...mapState('vehicles', ['currencyMasters', 'distanceUnitMasters', 'roleMasters']),
  },
  mounted() {
    this.getUsers()
  },
  methods: {
    getUsers() {
      store.dispatch('users/users').then((data) => {
        this.users = data
      })
    },
    formatDate(date) {
      return parseAndFormatDate(date)
    },
    changeDisabledStatus(userId,status){
        this.$buefy.dialog.confirm({
        title: status?'Disable User':"Enable User",
        message: 'Are you sure you want to do this?',
        cancelText: 'Cancel',
        confirmText: 'Go Ahead',
        onConfirm: () => {

          var url = `/api/users/${userId}/${status?"disable":"enable"}`
          axios
            .post(url, {})
            .then((data) => {
              this.$buefy.toast.open({
                message: status?"User disabled successfully":'User enabled successfully',
                type: 'is-success',
                duration: 3000,
              })
             this.getUsers();
            })
            .catch((ex) => {
              this.$buefy.toast.open({
                duration: 5000,
                message: ex.message,
                position: 'is-bottom',
                type: 'is-danger',
              })
            })
        },
      })
    },
    resetUserForm() {
      this.registerModel = {
        name: '',
        email: '',
        password: '',
        distanceUnit: this.settings.distanceUnit,
        currency: this.settings.currency,
        role: 1,
      }
      this.showUserForm = false
    },
    register() {
      this.isWorking = true
      axios
        .post('/api/register', this.registerModel)
        .then((response) => {
          const success = response.data.success
          if (success) {
            this.$buefy.toast.open({
              duration: 10000,
              message: 'User Created Successfully',
              position: 'is-bottom',
              type: 'is-success',
            })
            this.getUsers()
            this.resetUserForm()
          }
        })
        .catch((ex) => {
          this.$buefy.toast.open({
            duration: 5000,
            message: ex.message,
            position: 'is-bottom',
            type: 'is-danger',
          })
        })
        .finally(() => (this.isWorking = false))
    },
  },
}
</script>

<template>
  <Layout>
    <div class="box">
      <div class="columns">
        <div class="column is-three-quarters"> <h1 class="title is-4">Users</h1> </div>
        <div class="column is-one-quarter">
          <b-button type="is-primary" @click="showUserForm = true">Add User</b-button>
        </div>
      </div>

      <div v-if="showUserForm" class="box content">
        <h1 class="title">Create New User</h1>
        <form @submit.prevent="register">
          <b-field label="Name">
            <b-input v-model="registerModel.name" required></b-input>
          </b-field>
          <b-field label="Email">
            <b-input v-model="registerModel.email" type="email" required></b-input>
          </b-field>
          <b-field label="Password">
            <b-input
              v-model="registerModel.password"
              type="password"
              required
              minlength="6"
              password-reveal
            ></b-input>
          </b-field>
          <b-field label="Role">
            <b-select v-model.number="registerModel.role" placeholder="Role" required expanded>
              <option v-for="(option, key) in roleMasters" :key="key" :value="key">
                {{ `${option.long}` }}
              </option>
            </b-select>
          </b-field>
          <b-field label="Currency">
            <b-select v-model="registerModel.currency" placeholder="Currency" required expanded>
              <option v-for="option in currencyMasters" :key="option.code" :value="option.code">
                {{ `${option.namePlural} (${option.code})` }}
              </option>
            </b-select>
          </b-field>
          <b-field label="Distance Unit">
            <b-select
              v-model.number="registerModel.distanceUnit"
              placeholder="Distance Unit"
              required
              expanded
            >
              <option v-for="(option, key) in distanceUnitMasters" :key="key" :value="key">
                {{ `${option.long} (${option.short})` }}
              </option>
            </b-select>
          </b-field>
          <br />
          <div class="buttons">
            <b-button type="is-primary" native-type="submit" tag="input"></b-button>

            <b-button type="is-danger is-light" @click="resetUserForm">Cancel</b-button>
          </div>
        </form>
      </div>
      <b-table :data="users" hoverable mobile-cards detail-key="id" paginated per-page="10" :row-class="(row, index) => row.isDisabled && 'is-disabled'">
        <b-table-column v-slot="props" field="name" label="Name">
          {{ `${props.row.name}` }} <template v-if="props.row.id === user.id">(You)</template>
        </b-table-column>
        <b-table-column v-slot="props" field="email" label="Email">
          {{ `${props.row.email}` }}
        </b-table-column>
        <b-table-column v-slot="props" field="role" label="Role">
          {{ `${props.row.roleDetail.short}` }}
        </b-table-column>
        <b-table-column v-slot="props" field="createdAt" label="Created" sortable date>
          {{ formatDate(props.row.createdAt) }}
        </b-table-column>
         <b-table-column v-slot="props">
           <b-button type="is-success" v-if="props.row.isDisabled && props.row.roleDetail.long === 'USER'" @click="changeDisabledStatus(props.row.id, false)">Enable</b-button>
            <b-button type="is-danger" v-if="!props.row.isDisabled && props.row.roleDetail.long === 'USER'" @click="changeDisabledStatus(props.row.id, true)">Disable</b-button>
         </b-table-column>
      </b-table>
    </div>
  </Layout>
</template>
