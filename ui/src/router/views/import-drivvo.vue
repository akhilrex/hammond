<script>
import Layout from '@layouts/main.vue'
import { mapState } from 'vuex'
import axios from 'axios'

export default {
  page: {
    title: 'Import Drivvo',
    meta: [{ name: 'description', content: 'The Import Drivvo page.' }],
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
      myVehicles: [],
      file: null,
      selectedVehicle: null,
      tryingToCreate: false,
      errors: [],
      importLocation: true,
    }
  },
  computed: {
    ...mapState('utils', ['isMobile']),
    ...mapState('vehicles', ['vehicles']),
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
  mounted() {
    this.myVehicles = this.vehicles
  },
  methods: {
    importDrivvo() {
      console.log('Import from drivvo')
      if (this.file == null) {
        return
      }
      this.tryingToCreate = true
      this.errorMessage = ''
      const formData = new FormData()
      formData.append('vehicleID', this.selectedVehicle)
      formData.append('importLocation', this.importLocation)
      formData.append('file', this.file, this.file.name)
      axios
        .post(`/api/import/drivvo`, formData)
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
        <h1 class="title">Import from Drivvo</h1>
      </div>
    </div>
    <br />
    <div class="columns">
      <div class="column">
        <p class="subtitle"> Steps to import data from Drivvo</p>
        <ol>
          <li>Export your data from Drivvo in the CSV format.</li>
          <li>Select the vehicle the exported data is for. You may need to create the vehicle in Hammond first if you haven't already done so</li>
          <li
            >Make sure that the <u>Currency</u> and <u>Distance Unit</u> are set correctly in Hammond. Drivvo does not include this information in
            their export, instead Hammond will use the values set for the user.</li
          >
          <li>Similiarly, make sure that the <u>Fuel Unit</u> and <u>Fuel Type</u> are correctly set in the Vehicle.</li>
          <li>Once you have checked all these points, select the vehicle and import the CSV below.</li>
          <li><b>Make sure that you do not import the file again as that will create repeat entries.</b></li>
        </ol>
      </div>
    </div>
    <p
      ><b>PS:</b> If you have <em>'income'</em> and <em>'trips'</em> in your export, they will not be imported to Hammond. The fields
      <em>'Second fuel'</em> and <em>'Third fuel'</em> are are are also ignored as the use case for these is not understood by us. If you have a use
      case for this, please open a issue on
      <a href="https://github.com/akhilrex/hammond/issues">issue tracker</a>
    </p>
    <div class="section box">
      <div class="columns is-multiline">
        <div class="column is-full"> <p class="subtitle">Choose the vehicle, then select the Drivvo CSV and press the import button.</p></div>
        <div class="column is-full is-flex is-align-content-center">
          <form @submit.prevent="importDrivvo">
            <div class="columns">
              <div class="column">
                <b-field label="Vehicle" label-position="on-border">
                  <b-select v-model="selectedVehicle" placeholder="Select Vehicle" required>
                    <option v-for="vehicle in myVehicles" :key="vehicle.id" :value="vehicle.id">{{ vehicle.nickname }}</option>
                  </b-select>
                </b-field>
              </div>
              <div class="column">
                <b-field>
                  <b-tooltip label="Whether to import the location for fillups and services or not." multilined>
                    <b-checkbox v-model="importLocation">Import Location?</b-checkbox>
                  </b-tooltip>
                </b-field>
              </div>

              <div class="column">
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
