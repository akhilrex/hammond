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
          return 'Choose Photo'
        } else {
          return ''
        }
      } else {
        if (this.file == null) {
          return 'Choose CSV'
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
            message: 'Data Imported Successfully',
            type: 'is-success',
            duration: 3000,
          })
          this.file = null
        })
        .catch((ex) => {
          this.$buefy.toast.open({
            duration: 5000,
            message: 'There was some issue with importing the file. Please check the error message',
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
        <h1 class="title">Import from Fuelly</h1>
      </div>
    </div>
    <br />
    <div class="columns">
      <div class="column">
        <p class="subtitle"> Steps to import data from Fuelly</p>
        <ol>
          <li
            >Export your data from Fuelly in the CSV format. Steps to do that can be found
            <a href="http://docs.fuelly.com/acar-import-export-center" target="_nofollow">here</a>.</li
          >
          <li>Make sure that you have already created the vehicles in Hammond platform.</li>
          <li>Make sure that the Vehicle nickname in Hammond is exactly the same as the name on Fuelly CSV or the import will not work.</li>
          <li
            >Make sure that the <u>Currency</u> and <u>Distance Unit</u> are set correctly in Hammond. Import will not autodetect Currency from the
            CSV but use the one set for the user.</li
          >
          <li>Similiarly, make sure that the <u>Fuel Unit</u> and <u>Fuel Type</u> are correctly set in the Vehicle.</li>
          <li>Once you have checked all these points,just import the CSV below.</li>
          <li><b>Make sure that you do not import the file again and that will create repeat entries.</b></li>
        </ol>
      </div>
    </div>
    <div class="section box">
      <div class="columns">
        <div class="column is-two-thirds"> <p class="subtitle">Choose the Fuelly CSV and press the import button.</p></div>
        <div class="column is-one-third is-flex is-align-content-center">
          <form @submit.prevent="importFuelly">
            <div class="columns"
              ><div class="column">
                <b-field class="file is-primary" :class="{ 'has-name': !!file }">
                  <b-upload v-model="file" class="file-label" accept=".csv" required>
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
                <b-button tag="input" native-type="submit" :disabled="tryingToCreate" type="is-primary" value="Upload File" class="control">
                  Import
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
