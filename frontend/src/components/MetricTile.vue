<template>
  <div class="metric-tile">
    <p class="metric-label">{{ label }}</p>
    <p class="metric-value">{{ value }}</p>
    <p class="metric-footnote">{{ footnote }}</p>
  </div>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    required: true
  },
  value: {
    type: String,
    required: true
  },
  footnote: {
    type: String,
    required: true
  }
});
</script>
