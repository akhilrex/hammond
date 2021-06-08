<script>
import Layout from '@layouts/main.vue'
import { parseAndFormatDate } from '@utils/format-date'
import { mapState } from 'vuex'
import axios from 'axios'
import currencyFormtter from 'currency-formatter'
import store from '@state/store'
import ShareVehicle from '@components/shareVehicle.vue'
export default {
  page() {
    return {
      title: this.vehicle.nickname,
      meta: [
        {
          name: 'description',
          content: `The vehicle profile for ${this.vehicle.nickname}.`,
        },
      ],
    }
  },
  components: { Layout },
  props: {
    vehicle: {
      type: Object,
      required: true,
    },
  },
  data: function() {
    return {
      fillups: [],
      expenses: [],
      attachments: [],
      showAttachmentForm: false,
      file: null,
      tryingToUpload: false,
      title: '',
      stats: null,
      users: [],
    }
  },
  computed: {
    ...mapState('users', ['me']),
    ...mapState('auth', ['currentUser']),
    ...mapState('utils', ['isMobile']),
    summaryObject() {
      if (this.stats == null) {
        return []
      }
      return this.stats.map((x) => {
        return [
          {
            label: 'Currency',
            value: x.currency,
          },
          {
            label: 'Total Expenditure',
            value: this.formatCurrency(x.expenditureTotal, x.currency),
          },
          {
            label: 'Fillup Costs',
            value: `${this.formatCurrency(x.expenditureFillups, x.currency)} (${x.countFillups})`,
          },
          {
            label: 'Other Expenses',
            value: `${this.formatCurrency(x.expenditureExpenses, x.currency)} (${x.countExpenses})`,
          },
          {
            label: 'Avg Fillup Expense',
            value: `${this.formatCurrency(x.avgFillupCost, x.currency)}`,
          },
          {
            label: 'Avg Fillup Qty',
            value: `${x.avgFuelQty} ${this.vehicle.fuelUnitDetail.short}`,
          },
          {
            label: 'Avg Fuel Cost',
            value: `${this.formatCurrency(x.avgFuelPrice, x.currency)} per ${this.vehicle.fuelUnitDetail.short}`,
          },
        ]
      })
    },
  },
  mounted() {
    this.fetchFillups()
    this.fetchExpenses()

    this.fetchAttachments()
    this.fetchVehicleStats()
    this.fetchVehicleUsers()
  },
  methods: {
    fetchAttachments() {
      store
        .dispatch('vehicles/fetchAttachmentsByVehicleId', { vehicleId: this.vehicle.id })
        .then((data) => {
          this.attachments = data
        })
        .catch((err) => console.log(err))
    },
    fetchFillups() {
      axios
        .get(`/api/vehicles/${this.vehicle.id}/fillups`)
        .then((response) => {
          this.fillups = response.data
        })
        .catch((err) => console.log(err))
    },
    fetchExpenses() {
      axios
        .get(`/api/vehicles/${this.vehicle.id}/expenses`)
        .then((response) => {
          this.expenses = response.data
        })
        .catch((err) => console.log(err))
    },

    fetchVehicleStats() {
      store
        .dispatch('vehicles/fetchStatsByVehicleId', { vehicleId: this.vehicle.id })
        .then((data) => {
          this.stats = data
        })
        .catch((err) => console.log(err))
    },
    fetchVehicleUsers() {
      store
        .dispatch('vehicles/fetchUsersByVehicleId', { vehicleId: this.vehicle.id })
        .then((data) => {
          this.users = data
        })
        .catch((err) => console.log(err))
    },
    deleteFillup(fillupId) {
      var sure = confirm('This will delete this fillup. This step cannot be reversed. Are you sure?')
      if (sure) {
        store
          .dispatch('vehicles/deleteFillupById', { vehicleId: this.vehicle.id, fillupId: fillupId })
          .then((data) => {
            this.fetchVehicleStats()
            this.fetchFillups()
          })
          .catch((err) => console.log(err))
      }
    },
    deleteExpense(expenseId) {
      var sure = confirm('This will delete this expense. This step cannot be reversed. Are you sure?')
      if (sure) {
        store
          .dispatch('vehicles/deleteExpenseById', { vehicleId: this.vehicle.id, expenseId: expenseId })
          .then((data) => {
            this.fetchVehicleStats()
            this.fetchExpenses()
          })
          .catch((err) => console.log(err))
      }
    },
    deleteVehicle() {
      var sure = confirm(
        'This will delete all the expenses and fillups related with this vehicle as well. This step cannot be reversed. Are you sure?'
      )
      if (sure) {
        axios
          .delete(`/api/vehicles/${this.vehicle.id}`)
          .then((data) => {
            this.$buefy.toast.open({
              message: 'Vehicle Deleted Successfully',
              type: 'is-success',
              duration: 3000,
            })
            this.$router.push('/')
          })
          .catch((ex) => {
            this.$buefy.toast.open({
              duration: 5000,
              message: ex.message,
              position: 'is-bottom',
              type: 'is-danger',
            })
          })
          .finally(() => {})
      }
    },
    addAttachment() {
      if (this.file == null) {
        return
      }
      this.tryingToUpload = true
      const formData = new FormData()
      formData.append('file', this.file, this.file.name)
      formData.append('title', this.title)
      axios
        .post(`/api/vehicles/${this.vehicle.id}/attachments`, formData)
        .then((data) => {
          this.$buefy.toast.open({
            message: 'Quick Entry Created Successfully',
            type: 'is-success',
            duration: 3000,
          })
          this.file = null
          this.title = ''
          this.fetchAttachments()
        })
        .catch((ex) => {
          this.$buefy.toast.open({
            duration: 5000,
            message: ex.message,
            position: 'is-bottom',
            type: 'is-danger',
          })
        })
        .finally(() => {
          this.tryingToUpload = false
        })
    },
    formatDate(date) {
      return parseAndFormatDate(date)
    },
    formatCurrency(number, currencyCode) {
      if (!currencyCode) {
        currencyCode = this.me.currency
      }
      return currencyFormtter.format(number, { code: currencyCode })
    },
    columnTdAttrs(row, column) {
      return null
    },
    hiddenDesktop(row, column) {
      return {
        class: 'is-hidden-desktop',
      }
    },
    hiddenMobile(row, column) {
      return {
        class: 'is-hidden-mobile',
      }
    },
    showShareVehicleModal() {
      this.$buefy.modal.open({
        parent: this,
        component: ShareVehicle,
        hasModalCard: false,
        props: { vehicle: this.vehicle },
        onCancel: (x) => {
          this.fetchVehicleUsers()
        },
      })
    },
  },
}
</script>

<template>
  <Layout>
    <div class="columns box">
      <div class="column is-two-thirds" :class="isMobile ? 'has-text-centered' : ''">
        <p class="title">{{ vehicle.nickname }} - {{ vehicle.registration }}</p>
        <p class="subtitle">
          {{ [vehicle.make, vehicle.model, vehicle.fuelTypeDetail.long].join(' | ') }}

          <template v-if="users.length > 1">
            | Shared with :
            {{
              users
                .map((x) => {
                  if (x.userId === me.id) {
                    return 'You'
                  } else {
                    return x.name
                  }
                })
                .join(', ')
            }}
          </template>
        </p>
      </div>
      <div class="column is-one-third buttons has-text-centered">
        <b-button type="is-primary" tag="router-link" :to="`/vehicles/${vehicle.id}/fillup`">Add Fillup</b-button>
        <b-button type="is-primary" tag="router-link" :to="`/vehicles/${vehicle.id}/expense`">Add Expense</b-button>
        <b-button
          v-if="vehicle.isOwner"
          tag="router-link"
          title="Edit Vehicle"
          :to="{
            name: 'vehicle-edit',
            props: { vehicle: vehicle },
            params: { id: vehicle.id },
          }"
        >
          <b-icon pack="fas" icon="edit" type="is-info"> </b-icon
        ></b-button>
        <b-button v-if="vehicle.isOwner" title="Share vehicle" @click="showShareVehicleModal">
          <b-icon pack="fas" icon="share" type="is-info"> </b-icon
        ></b-button>
        <b-button v-if="vehicle.isOwner" title="Delete Vehicle" @click="deleteVehicle">
          <b-icon pack="fas" icon="trash" type="is-danger"> </b-icon
        ></b-button>
      </div>
    </div>
    <div v-for="(currencyLevel, index) in summaryObject" :key="index" class="level box">
      <div v-for="item in currencyLevel" :key="item.label" class="level-item has-text-centered">
        <div>
          <p class="heading">{{ item.label }}</p>
          <p class="title is-4">{{ item.value }}</p>
        </div>
      </div>
    </div>
    <div class="box">
      <h1 class="title is-4">Past Fillups</h1>

      <b-table :data="fillups" hoverable mobile-cards :detailed="isMobile" detail-key="id" paginated per-page="10">
        <b-table-column v-slot="props" field="date" label="Date" :td-attrs="columnTdAttrs" sortable date>
          {{ formatDate(props.row.date) }}
        </b-table-column>
        <b-table-column v-slot="props" field="fuelQuantity" label="Qty." :td-attrs="hiddenMobile" numeric>
          {{ `${props.row.fuelQuantity} ${props.row.fuelUnitDetail.short}` }}
        </b-table-column>
        <b-table-column
          v-slot="props"
          field="perUnitPrice"
          :label="'Price per ' + vehicle.fuelUnitDetail.short"
          :td-attrs="hiddenMobile"
          numeric
          sortable
        >
          {{ `${formatCurrency(props.row.perUnitPrice, props.row.currency)}` }}
        </b-table-column>
        <b-table-column v-if="isMobile" v-slot="props" field="totalAmount" label="Total" :td-attrs="hiddenDesktop" sortable numeric>
          {{ `${me.currency} ${props.row.totalAmount}` }} ({{ `${props.row.fuelQuantity} ${props.row.fuelUnitDetail.short}` }} @
          {{ `${me.currency} ${props.row.perUnitPrice}` }})
        </b-table-column>
        <b-table-column v-if="!isMobile" v-slot="props" field="totalAmount" label="Total" :td-attrs="hiddenMobile" sortable numeric>
          {{ `${formatCurrency(props.row.totalAmount, props.row.currency)}` }}
        </b-table-column>
        <b-table-column v-slot="props" width="20" field="isTankFull" label="Tank Full" :td-attrs="hiddenMobile">
          <b-icon pack="fas" :icon="props.row.isTankFull ? 'check' : 'times'" type="is-info"> </b-icon>
        </b-table-column>
        <b-table-column v-slot="props" field="odoReading" label="Odometer Reading" :td-attrs="hiddenMobile" numeric>
          {{ `${props.row.odoReading} ${me.distanceUnitDetail.short}` }}
        </b-table-column>
        <b-table-column v-slot="props" field="fillingStation" label="Fillup Station" :td-attrs="hiddenMobile">
          {{ `${props.row.fillingStation}` }}
        </b-table-column>
        <b-table-column v-slot="props" field="userId" label="By" :td-attrs="hiddenMobile">
          {{ `${props.row.user.name}` }}
        </b-table-column>
        <b-table-column v-slot="props">
          <b-button
            type="is-ghost"
            tag="router-link"
            :to="{
              name: 'vehicle-edit-fillup',
              props: { fillup: props.row, vehicle: vehicle },
              params: { fillupId: props.row.id, id: vehicle.id },
            }"
          >
            <b-icon pack="fas" icon="edit" type="is-info"> </b-icon
          ></b-button>
          <b-button type="is-ghost" title="Delete this fillup" @click="deleteFillup(props.row.id)">
            <b-icon pack="fas" icon="trash" type="is-danger"> </b-icon
          ></b-button>
        </b-table-column>
        <template v-slot:empty> No Fillups so far</template>
        <template v-slot:detail="props">
          <p>{{ props.row.id }}</p>
        </template>
      </b-table>
    </div>
    <br />
    <div class="box">
      <h1 class="title is-4">Past Expenses</h1>

      <b-table :data="expenses" hoverable mobile-cards paginated per-page="10">
        <b-table-column v-slot="props" field="date" label="Date" :td-attrs="columnTdAttrs" date>
          {{ formatDate(props.row.date) }}
        </b-table-column>
        <b-table-column v-slot="props" field="expenseType" label="Expense Type">
          {{ `${props.row.expenseType}` }}
        </b-table-column>

        <b-table-column v-slot="props" field="amount" label="Total" :td-attrs="hiddenMobile" sortable numeric>
          {{ `${formatCurrency(props.row.amount, props.row.currency)}` }}
        </b-table-column>

        <b-table-column v-slot="props" field="odoReading" label="Odometer Reading" :td-attrs="columnTdAttrs" numeric>
          {{ `${props.row.odoReading} ${me.distanceUnitDetail.short}` }}
        </b-table-column>

        <b-table-column v-slot="props" field="userId" label="By" :td-attrs="columnTdAttrs">
          {{ `${props.row.user.name}` }}
        </b-table-column>
        <b-table-column v-slot="props">
          <b-button
            type="is-ghost"
            tag="router-link"
            :to="{
              name: 'vehicle-edit-expense',
              props: { expense: props.row, vehicle: vehicle },
              params: { expenseId: props.row.id, id: vehicle.id },
            }"
          >
            <b-icon pack="fas" icon="edit" type="is-info"> </b-icon
          ></b-button>
          <b-button type="is-ghost" title="Delete this expense" @click="deleteExpense(props.row.id)">
            <b-icon pack="fas" icon="trash" type="is-danger"> </b-icon
          ></b-button>
        </b-table-column>
        <template v-slot:empty> No Expenses so far</template>
      </b-table>
    </div>
    <br />
    <div class="box">
      <div class="columns">
        <div class="column is-three-quarters"> <h1 class="title is-4">Attachments</h1></div>
        <div class="column buttons">
          <b-button type="is-primary" @click="showAttachmentForm = true">
            Add Attachment
          </b-button>
        </div>
      </div>
      <div v-if="showAttachmentForm" class="box">
        <div class="columns">
          <div class="column"></div>
          <div class="column is-two-thirds">
            <form @submit.prevent="addAttachment">
              <b-field :grouped="!isMobile">
                <b-field class="file is-primary" :class="{ 'has-name': !!file }">
                  <b-upload v-model="file" class="file-label" required>
                    <span class="file-cta">
                      <b-icon class="file-icon" icon="upload"></b-icon>
                      <span class="file-label">Choose File</span>
                    </span>
                    <span v-if="file" class="file-name" :class="isMobile ? 'file-name-mobile' : 'file-name-desktop'">
                      {{ file.name }}
                    </span>
                  </b-upload>
                </b-field>
                <b-field>
                  <b-input v-model="title" required placeholder="Label for this file"></b-input>
                </b-field>

                <b-field class="buttons">
                  <b-button tag="input" native-type="submit" :disabled="tryingToUpload" type="is-primary" label="Upload File" value="Upload File">
                  </b-button>
                  <b-button
                    tag="input"
                    native-type="submit"
                    :disabled="tryingToUpload"
                    type="is-danger"
                    label="Upload File"
                    value="Cancel"
                    @click="showAttachmentForm = false"
                  >
                  </b-button>
                </b-field>
              </b-field>
            </form>
          </div>
        </div>
      </div>

      <b-table :data="attachments" hoverable mobile-cards>
        <b-table-column v-slot="props" field="title" label="Title" :td-attrs="columnTdAttrs">
          {{ `${props.row.title}` }}
        </b-table-column>

        <b-table-column v-slot="props" field="originalName" label="Name" :td-attrs="columnTdAttrs">
          {{ `${props.row.originalName}` }}
        </b-table-column>
        <b-table-column v-slot="props" field="id" label="Download" :td-attrs="columnTdAttrs">
          <b-button tag="a" :href="`/api/attachments/${props.row.id}/file?access_token=${currentUser.token}`" :download="props.row.originalName">
            <b-icon type="is-primary" icon="download"></b-icon>
          </b-button>
        </b-table-column>
        <template v-slot:empty> No Attachments so far</template>
      </b-table>
    </div>
  </Layout>
</template>
