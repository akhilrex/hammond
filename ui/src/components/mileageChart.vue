<script>
import { Line } from 'vue-chartjs'

import axios from 'axios'
import { mapState } from 'vuex'
import { string } from 'yargs'
export default {
  extends: Line,
  props: {
    vehicle: { type: Object, required: true },
    since: { type: Date, default: '' },
    user: { type: Object, required: true },
    mileageOption: { type: string, default: 'litre_100km' },
  },
  data: function() {
    return {
      chartData: [],
    }
  },
  computed: {
    ...mapState('utils', ['isMobile']),
  },
  watch: {
    since(newOne, old) {
      if (newOne === old) {
        return
      }
      this.fetchMileage()
    },
  },
  mounted() {
    this.fetchMileage()
  },
  methods: {
    showChart() {
      let mileageLabel = ''
      switch (this.mileageOption) {
        case 'litre_100km':
          mileageLabel = 'L/100km'
          break
        case 'km_litre':
          mileageLabel = 'km/L'
          break
        case 'mpg':
          mileageLabel = 'mpg'
          break
      }

      var labels = this.chartData.map((x) => x.date.substr(0, 10))
      var dataset = {
        steppedLine: true,
        label: `Mileage (${mileageLabel})`,
        fill: true,
        data: this.chartData.map((x) => x.mileage),
      }
      this.renderChart({ labels, datasets: [dataset] }, { maintainAspectRatio: false })
    },
    fetchMileage() {
      axios
        .get(`/api/vehicles/${this.vehicle.id}/mileage`, {
          params: {
            since: this.since,
            mileageOption: this.mileageOption,
          },
        })
        .then((response) => {
          this.chartData = response.data
          this.showChart()
        })
        .catch((err) => console.log(err))
    },
  },
}
</script>
