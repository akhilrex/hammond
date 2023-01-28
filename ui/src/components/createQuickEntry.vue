<script>
import { mapState } from 'vuex'
import store from '@state/store'
import axios from 'axios'
export default {
  data: function() {
    return {
      file: null,
      tryingToCreate: false,
    }
  },
  computed: {
    ...mapState('utils', ['isMobile']),
    uploadButtonLabel() {
      if (this.isMobile) {
        if (this.file == null) {
          return 'Upload/Click Photo'
        } else {
          return ''
        }
      } else {
        if (this.file == null) {
          return this.$t('uploadphoto')
        } else {
          return ''
        }
      }
    },
  },
  methods: {
    createQuickEntry() {
      if (this.file == null) {
        return
      }
      this.tryingToCreate = true
      const formData = new FormData()
      formData.append('file', this.file, this.file.name)
      axios
        .post(`/api/quickEntries`, formData)
        .then((data) => {
          this.$buefy.toast.open({
            message: this.$t('quickentrycreatedsuccessfully'),
            type: 'is-success',
            duration: 3000,
          })
          this.file = null
          store.dispatch('vehicles/fetchQuickEntries', { force: true }).then((data) => {
            this.quickEntries = data
          })
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
          this.tryingToCreate = false
        })
    },
  },
}
</script>

<template>
  <div class="section box">
    <div class="columns">
      <div class="column is-two-thirds">
        <p class="title">{{ $tc('quickentry',1) }}</p>
        <p class="subtitle"
          >{{ $t('quickentrydesc') }}</p
        ></div
      >
      <div class="column is-one-third is-flex is-align-content-center">
        <form @submit.prevent="createQuickEntry">
          <div class="columns"
            ><div class="column">
              <b-field class="file is-primary" :class="{ 'has-name': !!file }">
                <b-upload v-model="file" class="file-label" accept="image/*">
                  <span class="file-cta">
                    <b-icon class="file-icon" icon="upload"></b-icon>
                    <span class="file-label">{{ uploadButtonLabel }}</span>
                  </span>
                  <span
                    v-if="file"
                    class="file-name"
                    :class="isMobile ? 'file-name-mobile' : 'file-name-desktop'"
                  >
                    {{ file.name }}
                  </span>
                </b-upload>
              </b-field>
            </div>
            <div class="column">
              <b-button
                tag="input"
                native-type="submit"
                :disabled="tryingToCreate"
                type="is-primary"
                :value="this.$t('uploadfile')"
                class="control"
              >
                {{ $t('uploadfile') }}
              </b-button>
            </div></div
          >
        </form>
      </div>
    </div>
  </div>
</template>

<style>
.file-name-desktop {
  max-width: 9em;
}
.file-name-mobile {
  max-width: 12em;
}
</style>
