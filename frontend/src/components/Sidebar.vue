<template>
  <aside class="sidebar">
    <div class="brand">
      <div class="brand-mark">EMS</div>
      <div>
        <p class="brand-name">Enterprise Portal</p>
        <p class="brand-subtitle">Microservice System</p>
      </div>
    </div>
    <nav class="nav">
      <router-link to="/" class="nav-link">Overview</router-link>
      <router-link to="/operations" class="nav-link">Operations</router-link>
      <router-link to="/observability" class="nav-link">Observability</router-link>
      <router-link to="/docs" class="nav-link">Docs & Access</router-link>
    </nav>
    <div class="sidebar-footer">
      <p class="small">Environment</p>
      <div class="chip">Local / Docker</div>
      <div class="status-row">
        <span class="status-dot"></span>
        <span>All services monitored</span>
      </div>
    </div>
  </aside>
</template>

<script setup>
</script>
