<template>
  <fragment>
    <v-chip
      v-if="isEnterprise"
      class="float-right"
      data-test="devices-chip"
      :to="{name: 'namespaceSettings'}"
    >
      {{ namespace.devices_count }}

      <font-awesome-icon
        class="ml-3 fa-lg"
        icon="microchip"
      />
    </v-chip>
  </fragment>
</template>

<script>

export default {
  name: 'MaxDevice',

  computed: {
    namespace() {
      return this.$store.getters['namespaces/get'];
    },

    isEnterprise() {
      return this.$env.isEnterprise;
    },
  },
};

</script>
