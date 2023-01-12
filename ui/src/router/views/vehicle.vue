<script>
import Layout from '@layouts/main.vue'
import { parseAndFormatDate } from '@utils/format-date'
import { mapState } from 'vuex'
import { addDays, addMonths } from 'date-fns'
import axios from 'axios'
import currencyFormatter from 'currency-formatter'
import store from '@state/store'
import ShareVehicle from '@components/shareVehicle.vue'
import MileageChart from '@components/mileageChart.vue'

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
  components: { Layout, MileageChart },
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
      dateRangeOptions: [
        { label: this.$t('thisweek'), value: 'this_week' },
        { label: this.$t('thismonth'), value: 'this_month' },
        { label: this.$tc('pastxdays', 30), value: 'past_30_days' },
        { label: this.$tc('pastxmonths', 3), value: 'past_3_months' },
        { label: this.$t('thisyear'), value: 'this_year' },
        { label: this.$t('alltime'), value: 'all_time' },
      ],
      dateRangeOption: 'past_30_days',
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
            label: this.$t('currency'),
            value: x.currency,
          },
          {
            label: this.$t('totalexpenses'),
            value: this.formatCurrency(x.expenditureTotal, x.currency),
          },
          {
            label: this.$t('fillupcost'),
            value: `${this.formatCurrency(x.expenditureFillups, x.currency)} (${x.countFillups})`,
          },
          {
            label: this.$t('otherexpenses'),
            value: `${this.formatCurrency(x.expenditureExpenses, x.currency)} (${x.countExpenses})`,
          },
          {
            label: this.$t('avgfillupexpense'),
            value: `${this.formatCurrency(x.avgFillupCost, x.currency)}`,
          },
          {
            label: this.$t('avgfillupqty'),
            value: `${x.avgFuelQty} ${this.vehicle.fuelUnitDetail.short}`,
          },
          {
            label: this.$t('avgfuelcost'),
            value: this.$t('per', {'0': this.formatCurrency(x.avgFuelPrice, x.currency), '1': this.vehicle.fuelUnitDetail.short}),
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
      // const config = { headers: { 'Content-Type': 'multipart/form-data; boundary=' + formData._boundary } }
      fetch(`/api/vehicles/${this.vehicle.id}/attachments`, {
        method: 'POST',
        body: formData,
        headers: {
          Authorization: this.currentUser.token,
        },
      })
        .then((data) => {
          this.$buefy.toast.open({
            message: 'File uploaded Successfully',
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
      return currencyFormatter.format(number, { code: currencyCode })
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
    getStartDate() {
      const toDate = new Date()
      switch (this.dateRangeOption) {
        case 'this_week':
          var currentDayOfWeek = toDate.getDay()
          var toSubtract = 0
          if (currentDayOfWeek === 0) {
            toSubtract = -6
          }
          if (currentDayOfWeek > 1) {
            toSubtract = -1 * (currentDayOfWeek - 1)
          }
          return addDays(toDate, toSubtract)
        case 'this_month':
          return new Date(toDate.getFullYear(), toDate.getMonth(), 1)
        case 'past_30_days':
          return addDays(toDate, -30)
        case 'past_3_months':
          return addMonths(toDate, -3)
        case 'this_year':
          return new Date(toDate.getFullYear(), 0, 1)
        case 'all_time':
          return new Date(1969, 4, 20)
        default:
          return new Date(1969, 4, 20)
      }
    },
  },
}
</script>

<template>
  <Layout>
    <div class="columns box">
      <div class="column is-one-half" :class="isMobile ? 'has-text-centered' : ''">
        <p class="title">{{ vehicle.nickname }} - {{ vehicle.registration }}</p>
        <p class="subtitle">
          {{ [vehicle.make, vehicle.model, this.$t('fuel.' + vehicle.fuelTypeDetail.short)].join(' | ') }}

          <template v-if="users.length > 1">
            | {{ $t("sharedwith") }} :
            {{
              users
                .map((x) => {
                  if (x.userId === me.id) {
                    return this.$t('you')
                  } else {
                    return x.name
                  }
                })
                .join(', ')
            }}
          </template>
        </p>
      </div>
      <div :class="(!isMobile ? 'has-text-right ' : '') + 'column is-one-half buttons'">
        <b-button type="is-primary" tag="router-link" :to="`/vehicles/${vehicle.id}/fillup`">{{ this.$t('addfillup') }}</b-button>
        <b-button type="is-primary" tag="router-link" :to="`/vehicles/${vehicle.id}/expense`">{{ this.$t('addexpense') }}</b-button>
        <b-button
          v-if="vehicle.isOwner"
          tag="router-link"
          :title="$t('editvehicle')"
          :to="{
            name: 'vehicle-edit',
            props: { vehicle: vehicle },
            params: { id: vehicle.id },
          }"
        >
          <b-icon pack="fas" icon="edit" type="is-info"> </b-icon
        ></b-button>
        <b-button v-if="vehicle.isOwner" :title="$t('sharevehicle')" @click="showShareVehicleModal">
          <b-icon pack="fas" icon="user-friends" type="is-info"> </b-icon
        ></b-button>
        <b-button v-if="vehicle.isOwner" :title="$t('deletevehicle')" @click="deleteVehicle">
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
      <h1 class="title is-4">{{ $t('pastfillups') }}</h1>

      <b-table :data="fillups" hoverable mobile-cards :detailed="isMobile" detail-key="id" paginated per-page="10">
        <b-table-column v-slot="props" field="date" :label="this.$t('date')" :td-attrs="columnTdAttrs" sortable date>
          {{ formatDate(props.row.date) }}
        </b-table-column>
        <b-table-column v-slot="props" field="fuelSubType" :label="this.$t('fuelsubtype')" :td-attrs="columnTdAttrs">
          {{ props.row.fuelSubType }}
        </b-table-column>
        <b-table-column v-slot="props" field="fuelQuantity" :label="this.$t('quantity')" :td-attrs="hiddenMobile" numeric>
          {{ `${props.row.fuelQuantity} ${props.row.fuelUnitDetail.short}` }}
        </b-table-column>
        <b-table-column
          v-slot="props"
          field="perUnitPrice"
          :label="this.$t('per', { '0': this.$t('price'), '1': vehicle.fuelUnitDetail.short })"
          :td-attrs="hiddenMobile"
          numeric
          sortable
        >
          {{ `${formatCurrency(props.row.perUnitPrice, props.row.currency)}` }}
        </b-table-column>
        <b-table-column v-if="isMobile" v-slot="props" field="totalAmount" :label="this.$t('total')" :td-attrs="hiddenDesktop" sortable numeric>
          {{ `${me.currency} ${props.row.totalAmount}` }} ({{ `${props.row.fuelQuantity} ${props.row.fuelUnitDetail.short}` }} @
          {{ `${me.currency} ${props.row.perUnitPrice}` }})
        </b-table-column>
        <b-table-column v-if="!isMobile" v-slot="props" field="totalAmount" :label="this.$t('total')" :td-attrs="hiddenMobile" sortable numeric>
          {{ `${formatCurrency(props.row.totalAmount, props.row.currency)}` }}
        </b-table-column>
        <b-table-column v-slot="props" width="20" field="isTankFull" :label="this.$t('fulltank')" :td-attrs="hiddenMobile">
          <b-icon pack="fas" :icon="props.row.isTankFull ? 'check' : 'times'" type="is-info"> </b-icon>
        </b-table-column>
        <b-table-column v-slot="props" field="odoReading" :label="this.$t('odometer')" :td-attrs="hiddenMobile" numeric>
          {{ `${props.row.odoReading} ${me.distanceUnitDetail.short}` }}
        </b-table-column>
        <b-table-column v-slot="props" field="fillingStation" :label="this.$t('gasstation')" :td-attrs="hiddenMobile">
          {{ `${props.row.fillingStation}` }}
        </b-table-column>
        <b-table-column v-slot="props" field="userId" :label="this.$t('by')" :td-attrs="hiddenMobile">
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
          <b-button type="is-ghost" :title="$t('deletefillup')" @click="deleteFillup(props.row.id)">
            <b-icon pack="fas" icon="trash" type="is-danger"> </b-icon
          ></b-button>
        </b-table-column>
        <template v-slot:empty> {{ $t('nofillups') }}</template>
        <template v-slot:detail="props">
          <p>{{ props.row.id }}</p>
        </template>
      </b-table>
    </div>
    <br />
    <div class="box">
      <h1 class="title is-4">{{ $t('expenses') }}</h1>

      <b-table :data="expenses" hoverable mobile-cards paginated per-page="10">
        <b-table-column v-slot="props" field="date" :label="this.$t('date')" :td-attrs="columnTdAttrs" date>
          {{ formatDate(props.row.date) }}
        </b-table-column>
        <b-table-column v-slot="props" field="expenseType" :label="this.$t('expensetype')">
          {{ `${props.row.expenseType}` }}
        </b-table-column>

        <b-table-column v-slot="props" field="amount" :label="this.$t('total')" :td-attrs="hiddenMobile" sortable numeric>
          {{ `${formatCurrency(props.row.amount, props.row.currency)}` }}
        </b-table-column>

        <b-table-column v-slot="props" field="odoReading" :label="this.$t('odometer')" :td-attrs="columnTdAttrs" numeric>
          {{ `${props.row.odoReading} ${me.distanceUnitDetail.short}` }}
        </b-table-column>

        <b-table-column v-slot="props" field="userId" :label="this.$t('by')" :td-attrs="columnTdAttrs">
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
          <b-button type="is-ghost" :title="$t('deleteexpense')" @click="deleteExpense(props.row.id)">
            <b-icon pack="fas" icon="trash" type="is-danger"> </b-icon
          ></b-button>
        </b-table-column>
        <template v-slot:empty> {{ $t('noexpenses') }}</template>
      </b-table>
    </div>
    <br />
    <div class="box">
      <div class="columns">
        <div class="column is-three-quarters"> <h1 class="title is-4">{{ $t('attachments') }}</h1></div>
        <div class="column buttons">
          <b-button type="is-primary" @click="showAttachmentForm = true">
            {{ $t('addattachment') }}
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
                      <span class="file-label">{{ $t('choosefile') }}</span>
                    </span>
                    <span v-if="file" class="file-name" :class="isMobile ? 'file-name-mobile' : 'file-name-desktop'">
                      {{ file.name }}
                    </span>
                  </b-upload>
                </b-field>
                <b-field>
                  <b-input v-model="title" required :placeholder="this.$t('labelforfile')"></b-input>
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
        <b-table-column v-slot="props" field="title" :label="this.$t('title')" :td-attrs="columnTdAttrs">
          {{ `${props.row.title}` }}
        </b-table-column>

        <b-table-column v-slot="props" field="originalName" :label="this.$t('name')" :td-attrs="columnTdAttrs">
          {{ `${props.row.originalName}` }}
        </b-table-column>
        <b-table-column v-slot="props" field="id" :label="this.$t('download')" :td-attrs="columnTdAttrs">
          <b-button tag="a" :href="`/api/attachments/${props.row.id}/file?access_token=${currentUser.token}`" :download="props.row.originalName">
            <b-icon type="is-primary" icon="download"></b-icon>
          </b-button>
        </b-table-column>
        <template v-slot:empty> {{ $t('noattachments') }}</template>
      </b-table>
    </div>
    <div class="box">
      <div class="columns">
        <div class="column" :class="isMobile ? 'has-text-centered' : ''"> <h1 class="title">{{ $t('statistics') }}</h1></div>
        <div class="column">
          <b-select v-model="dateRangeOption" class="is-pulled-right is-medium">
            <option v-for="option in dateRangeOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </b-select></div
        >
      </div>
      <MileageChart :vehicle="vehicle" :since="getStartDate()" :user="me" :height="300" />
    </div>
  </Layout>
</template>
