<script>
import Layout from '@layouts/main.vue'
import { mapState } from 'vuex'
import axios from 'axios'

export default {
  page: {
    title: 'Import Fuelly',
    meta: [{ name: 'description', content: 'The Import Fuelly page.' }],
  },
  components: { Layout },
  computed: {
    ...mapState('utils', ['isMobile']),
    uploadButtonLabel() {
      if (this.isMobile) {
        if (this.file == null) {
          return this.$t('choosephoto')
        } else {
          return ''
        }
      } else {
        if (this.file == null) {
          return this.$t('choosecsv')
        } else {
          return ''
        }
      }
    },
  },
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  data: function() {
    return {
      file: null,
      tryingToCreate: false,
      errors: [],
    }
  },
  methods: {
    importFuelly() {
      if (this.file == null) {
        return
      }
      this.tryingToCreate = true
      this.errorMessage = ''
      const formData = new FormData()
      formData.append('file', this.file, this.file.name)
      axios
        .post(`/api/import/fuelly`, formData)
        .then((data) => {
          this.$buefy.toast.open({
            message: this.$t('importsuccessfull'),
            type: 'is-success',
            duration: 3000,
          })
          this.file = null
        })
        .catch((ex) => {
          this.$buefy.toast.open({
            duration: 5000,
            message: this.$t('importerror'),
            position: 'is-bottom',
            type: 'is-danger',
          })
          if (ex.response && ex.response.data.errors) {
            this.errors = ex.response.data.errors
          }
        })
        .finally(() => {
          this.tryingToCreate = false
        })
    },
  },
}
</script>

<template>
  <Layout>
    <div class="columns box">
      <div class="column">
        <h1 class="title">{{ $t('importfrom', { 'name': 'Fuelly' }) }}</h1>
      </div>
    </div>
    <br />
    <div class="columns">
      <div class="column">
        <p class="subtitle"> {{ $t('stepstoimport', { 'name': 'Fuelly' }) }}</p>
        <ol>
          <li>{{ $t('importhintcreatecsv', { 'name': 'Fuelly' }) }} <a href="http://docs.fuelly.com/acar-import-export-center" target="_nofollow">{{ $t('here') }}</a>.</li>
          <li>{{ $t('importhintvehiclecreated') }}</li>
          <li>{{ $t('importhintnickname') }}</li>
          <li v-html="$t('importhintcurrdist')"></li>
          <li v-html="$t('importhintunits')"></li>
          <li>{{ $t('checkpointsimportcsv') }}</li>
          <li><b>{{ $t('dontimportagain') }}</b></li>
        </ol>
      </div>
    </div>
    <div class="section box">
      <div class="columns">
        <div class="column is-two-thirds"> <p class="subtitle">{{ $t('choosecsvimport', { 'name': 'Fuelly' }) }}</p></div>
        <div class="column is-one-third is-flex is-align-content-center">
          <form @submit.prevent="importFuelly">
            <div class="columns"
              ><div class="column">
                <b-field class="file is-primary" :class="{ 'has-name': !!file }">
                  <b-upload v-model="file" class="file-label" accept=".csv">
                    <span class="file-cta">
                      <b-icon class="file-icon" icon="upload"></b-icon>
                      <span class="file-label">{{ uploadButtonLabel }}</span>
                    </span>
                    <span v-if="file" class="file-name" :class="isMobile ? 'file-name-mobile' : 'file-name-desktop'">
                      {{ file.name }}
                    </span>
                  </b-upload>
                </b-field>
              </div>
              <div class="column">
                <b-button tag="input" native-type="submit" :disabled="tryingToCreate" type="is-primary" :value="this.$t('uploadfile')" class="control">
                  {{ $t('import') }}
                </b-button>
              </div></div
            >
          </form>
        </div>
      </div>
    </div>
    <b-message v-if="errors.length" type="is-danger">
      <ul>
        <li v-for="error in errors" :key="error">{{ error }}</li>
      </ul>
    </b-message>
  </Layout>
</template>
