<script>
import { mapGetters, mapState } from 'vuex'
import { parseAndFormatDateTime } from '@utils/format-date'

export default {
  model: {
    prop: 'quickEntry',
    event: 'change',
  },
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  data: function() {
    return {
      quickEntry: null,
    }
  },
  computed: {
    ...mapState('utils', ['isMobile']),
    ...mapGetters('vehicles', ['unprocessedQuickEntries', 'processedQuickEntries']),
  },
  methods: {
    parseAndFormatDateTime(date) {
      return parseAndFormatDateTime(date)
    },
    showQuickEntry(entry) {
      const h = this.$createElement
      const vnode = h('p', { class: 'image' }, [
        h('img', {
          attrs: {
            src: `/api/attachments/${entry.attachmentId}/file?access_token=${this.user.token}`,
          },
        }),
      ])
      this.$buefy.modal.open({
        content: [vnode],
      })
      this.$emit('change', entry)
    },
  },
}
</script>

<template>
  <div class="level">
    <b-field class="level-right">
      <b-select
        v-if="unprocessedQuickEntries.length"
        v-model="quickEntry"
        :placeholder="this.$t('referquickentry')"
        expanded
        @input="showQuickEntry($event)"
      >
        <option v-for="option in unprocessedQuickEntries" :key="option.id" :value="option">
          {{ $t('created') }}: {{ parseAndFormatDateTime(option.createdAt) }}
        </option>
      </b-select>
      <p class="control">
        <b-button v-if="quickEntry" type="is-primary" @click="showQuickEntry(quickEntry)"
          >Show</b-button
        ></p
      >
    </b-field>
  </div></template
>
