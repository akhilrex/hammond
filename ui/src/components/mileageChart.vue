<script>
import { Line } from 'vue-chartjs'

import axios from 'axios'
import { mapState } from 'vuex'
export default {
  extends: Line,
  props: { vehicle: { type: Object, required: true }, since: { type: Date, default: '' }, user: { type: Object, required: true } },
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
  data: function() {
    return {
      chartData: [],
    }
  },
  mounted() {
    this.fetchMileage()
  },
  methods: {
    showChart() {
      var labels = this.chartData.map((x) => x.date.substr(0, 10))
      var dataset = {
        steppedLine: true,
        label: `Mileage (${this.user.distanceUnitDetail.short}/${this.vehicle.fuelUnitDetail.short})`,
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
