<template>
  <article class="status-card" :class="variant">
    <div class="status-header">
      <div>
        <p class="status-title">{{ title }}</p>
        <p class="status-subtitle">{{ subtitle }}</p>
      </div>
      <span class="status-pill">{{ status }}</span>
    </div>
    <p class="status-message">{{ message }}</p>
    <div class="status-meta">
      <span>{{ lastChecked }}</span>
      <span>{{ endpoint }}</span>
    </div>
  </article>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  subtitle: {
    type: String,
    required: true
  },
  status: {
    type: String,
    required: true
  },
  message: {
    type: String,
    required: true
  },
  lastChecked: {
    type: String,
    required: true
  },
  endpoint: {
    type: String,
    required: true
  },
  variant: {
    type: String,
    default: 'neutral'
  }
});
</script>
