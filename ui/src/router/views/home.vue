<script>
import currencyFormtter from 'currency-formatter'

import appConfig from '@src/app.config'
import Layout from '@layouts/main.vue'
import CreateQuickEntry from '@components/createQuickEntry.vue'
import StatsWidget from '@components/statsWidget.vue'
import { chunk } from 'lodash'

import { parseAndFormatDate } from '@utils/format-date'
// import store from '@state/store'
import { mapGetters, mapState } from 'vuex'

export default {
  page: {
    title: 'Home',
    meta: [{ name: 'description', content: appConfig.description }],
  },
  components: { Layout, CreateQuickEntry, StatsWidget },
  props: {
    user: {
      type: Object,
      required: true,
    },
  },

  data: function() {
    return {
      myVehicles: [],

      tryingToCreate: false,
    }
  },
  computed: {
    ...mapState('users', ['me']),
    ...mapState('vehicles', ['vehicles']),
    ...mapState('utils', ['isMobile']),
    ...mapGetters('vehicles', ['unprocessedQuickEntries']),
    chunkedVehicles() {
      return chunk(this.myVehicles, 3)
    },
  },
  watch: {
    vehicles(old, newOne) {
      this.myVehicles = newOne
    },
  },
  mounted() {
    this.myVehicles = this.vehicles
  },
  methods: {
    formatDate(date) {
      return parseAndFormatDate(date)
    },
    formatCurrency(number) {
      return currencyFormtter.format(number, { code: this.me.currency })
    },
  },
}
</script>

<template>
  <Layout>
    <b-notification v-if="myVehicles.length === 0" type="is-warning is-light" :closable="false">
      <div class="columns">
        <div class="column">
          It seems you have not yet created a vehicle in the system. Start by creating an entry for
          one of the vehicles you want to track.
        </div>
        <div class="column" :class="!isMobile ? 'has-text-right' : ''">
          <b-button type="is-warning" class="" tag="router-link" :to="`/vehicles/create`"
            >Create Now</b-button
          ></div
        >
      </div>
    </b-notification>
    <b-notification
      v-if="unprocessedQuickEntries.length"
      type="is-warning is-light"
      :closable="false"
    >
      <div class="columns">
        <div class="column">
          {{
            `You have ${unprocessedQuickEntries.length} quick ${
              unprocessedQuickEntries.length === 1 ? 'entry' : 'entries'
            } pending to be processed.`
          }}
        </div>
        <div class="column" :class="!isMobile ? 'has-text-right' : ''">
          <b-button type="is-warning" class="is-small" tag="router-link" :to="`/quickEntries`"
            >Process Now</b-button
          ></div
        >
      </div>
    </b-notification>
    <CreateQuickEntry />
    <StatsWidget :user="me" />
    <br />
    <section>
      <div class="columns" :class="isMobile ? 'has-text-centered' : ''"
        ><div class="column is-three-quarters"> <h1 class="title">Your Vehicles</h1></div>
        <div class="column is-one-quarter buttons" :class="!isMobile ? 'has-text-right' : ''">
          <b-button type="is-primary" tag="router-link" :to="`/vehicles/create`"
            >Add Vehicle</b-button
          >
        </div></div
      >

      <div v-if="myVehicles.length">
  <div v-for="chunk,index in chunkedVehicles" :key="index"  class="columns">
        <div v-for="vehicle in chunk" :key="vehicle.id" class="column is-4">
          <b-collapse animation="slide" aria-id="contentIdForA11y3" class="card" :open="!isMobile">
            <template v-slot:trigger="props">
              <div class="card-header" role="button" aria-controls="contentIdForA11y3">
                <div class="card-header-title">
                  <div>{{ `${vehicle.nickname} - ${vehicle.registration}` }} </div>
                </div>
                <a class="card-header-icon">
                  <b-icon :icon="props.open ? 'angle-down' : 'angle-up'"> </b-icon>
                </a>
              </div>
            </template>
            <div v-if="vehicle.fillups.length" class="card-content">
              <div class="content">
                <table class="table">
                  <div class="columns">
                    <div class="column is-one-third">Last Fillup</div>
                    <div class="column"
                      >{{ formatDate(vehicle.fillups[0].date) }} <br />
                      {{ `${formatCurrency(vehicle.fillups[0].totalAmount)}` }} ({{
                        `${vehicle.fillups[0].fuelQuantity} ${vehicle.fillups[0].fuelUnitDetail.short}`
                      }}
                      @ {{ `${formatCurrency(vehicle.fillups[0].perUnitPrice)}` }})</div
                    >
                  </div>

                  <div class="columns">
                    <div class="column is-one-third">Odometer</div>
                    <div class="column">
                      <template v-if="vehicle.fillups.length">
                        {{ vehicle.fillups[0].odoReading }}&nbsp;{{
                          me.distanceUnitDetail.short
                        }}</template
                      >
                    </div>
                  </div>
                </table>
              </div>
            </div>
            <footer class="card-footer">
              <router-link class="card-footer-item" :to="'/vehicles/' + vehicle.id">
                Details
              </router-link>
              <router-link class="card-footer-item" :to="`/vehicles/${vehicle.id}/fillup`">
                Add Fillup </router-link
              ><router-link class="card-footer-item" :to="`/vehicles/${vehicle.id}/expense`">
                Add Expense
              </router-link>
            </footer>
          </b-collapse>
        </div>
        </div>
      </div>
    </section>
  </Layout>
</template>
