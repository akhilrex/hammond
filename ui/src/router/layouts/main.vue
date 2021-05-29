<script>
import NavBar from '@components/nav-bar.vue'

export default {
  components: { NavBar },
}
</script>

<template
  ><div>
    <NavBar />
    <div class="section container"> <slot /> </div
  ></div>
</template>
