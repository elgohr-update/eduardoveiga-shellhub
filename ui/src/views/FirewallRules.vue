<template>
  <div>
    <Firewall data-test="firewall-component" />
  </div>
</template>

<script>

import Firewall from '@/components/firewall_rule/FirewallRule';

export default {
  name: 'FirewallRules',

  components: {
    Firewall,
  },
};

</script>
