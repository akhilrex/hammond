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
        <h1 class="title">{{ $t('importdata') }}</h1>
        <p class="subtitle">{{ $t('importdatadesc') }}</p>
      </div></div
    >
    <br />
    <div class="columns">
      <div class="box column is-one-third" to="/import-fuelly">
        <h1 class="title">Fuelly</h1>
        <p>{{ $t('importcsv', { 'name': 'Fuelly' }) }}</p>
        <br />
        <b-button type="is-primary" tag="router-link" to="/import/fuelly">{{ $t('import') }}</b-button>
      </div>
    </div>
  </Layout>
</template>
