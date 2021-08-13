<script>
import Layout from '@layouts/main.vue'
import { parseAndFormatDateTime } from '@utils/format-date'
import store from '@state/store'
import { chunk, filter } from 'lodash'
import { mapState, mapGetters } from 'vuex'
// import axios from 'axios'

export default {
  page: {
    title: 'Quick Entries',
  },
  components: { Layout },
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  data: function() {
    return {
      showUnprocessedOnly: false,
    }
  },
  computed: {
    chunkedQuickEntries() {
      var source = this.showUnprocessedOnly ? filter(this.quickEntries, (x) => x.processDate === null) : this.quickEntries
      return chunk(source, 3)
    },
    ...mapState('vehicles', ['vehicles', 'quickEntries']),
    ...mapGetters('vehicles', ['unprocessedQuickEntries']),
  },
  created() {
    store.dispatch('vehicles/fetchQuickEntries', { force: true }).then((data) => {})
  },
  methods: {
    parseAndFormatDateTime(date) {
      return parseAndFormatDateTime(date)
    },
    markProcessed(entry) {
      store.dispatch('vehicles/setQuickEntryAsProcessed', { id: entry.id }).then((data) => {})
    },
    deleteQuickEntry(entry) {
      var sure = confirm('This will delete this Quick Entry. This step cannot be reversed. Are you sure?')
      if (sure) {
        store.dispatch('vehicles/deleteQuickEntry', { id: entry.id }).then((data) => {})
      }
    },
    imageModal(url) {
      const h = this.$createElement
      const vnode = h('p', { class: 'image' }, [h('img', { attrs: { src: url } })])
      this.$buefy.modal.open({
        content: [vnode],
      })
    },
  },
}
</script>

<template>
  <Layout>
    <h1 class="title">Quick Entries</h1>
    <b-field>
      <b-switch v-if="unprocessedQuickEntries.length" v-model="showUnprocessedOnly">Show unprocessed only</b-switch>
    </b-field>
    <div v-for="(chunk, index) in chunkedQuickEntries" :key="index" class="tile is-ancestor">
      <div v-for="entry in chunk" :key="entry.id" class="tile is-parent" :class="{ 'is-4': quickEntries.length <= 3 }">
        <div class="tile is-child">
          <div class="card">
            <div class="card-header">
              <div class="card-header-title">
                {{ parseAndFormatDateTime(entry.createdAt) }}
              </div>
              <b-tag v-if="entry.processDate === null" class="is-align-content-center" type="is-primary">unprocessed</b-tag>
            </div>
            <div class="card-image">
              <!-- prettier-ignore -->
              <img
                class="is-clickable"
                :src="`/api/attachments/${entry.attachmentId}/file?access_token=${user.token}`"
                alt="Placeholder image"
                @click="imageModal(`/api/attachments/${entry.attachmentId}/file?access_token=${user.token}`)"
            />
            </div>
            <div class="card-content is-flex"
              ><p>{{ entry.comments }}</p></div
            >
            <footer class="card-footer">
              <router-link v-if="entry.processDate === null && vehicles.length" :to="`/vehicles/${vehicles[0].id}/fillup`" class="card-footer-item"
                >Create Fillup</router-link
              >
              <router-link v-if="entry.processDate === null && vehicles.length" :to="`/vehicles/${vehicles[0].id}/expense`" class="card-footer-item"
                >Create Expense</router-link
              >

              <a v-if="entry.processDate === null" class="card-footer-item" @click="markProcessed(entry)">Mark Processed</a>
              <p v-else class="card-footer-item">Processed on {{ parseAndFormatDateTime(entry.processDate) }}</p>
              <a class="card-footer-item" type="is-danger" @click="deleteQuickEntry(entry)"> Delete</a>
            </footer>
          </div>
        </div>
      </div>
    </div>
    <div v-if="!quickEntries.length" class="box">
      <p>No Quick Entries right now.</p>
    </div>
  </Layout>
</template>

<style module>
.card-equal-height {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.card-equal-height .card-footer {
  margin-top: auto;
}
</style>
