<script>
import Layout from '@layouts/main.vue'

export default {
  page: {
    title: 'Import',
    meta: [{ name: 'Import from other apps', content: 'The Import page.' }],
  },
  components: { Layout },
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
}
</script>

<template>
  <Layout>
    <div class="columns box"
      ><div class="column">
        <h1 class="title">Import data into Hammond</h1>
        <p class="subtitle">Choose from the following options to import data into Fuelly</p>
      </div></div
    >
    <br />
    <div class="columns">
      <div class="box column is-one-third" to="/import-fuelly">
        <h1 class="title">Fuelly</h1>
        <p>If you have been using Fuelly to store your vehicle data, export the CSV file from Fuelly and click here to import.</p>
        <br />
        <b-button type="is-primary" tag="router-link" to="/import/fuelly">Import</b-button>
      </div>
    </div>
  </Layout>
</template>
