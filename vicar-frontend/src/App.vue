<script setup lang="ts">
import VConfirmHost from "@/components/VConfirmHost.vue";
import VToastHost from "@/components/VToastHost.vue";
</script>

<template>
  <div class="smoke"/>
  <VToastHost>
    <VConfirmHost>
      <RouterView />
    </VConfirmHost>
  </VToastHost>
</template>

<style scoped lang="scss">

</style>
