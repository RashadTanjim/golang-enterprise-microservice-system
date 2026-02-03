<template>
  <div class="app-shell">
    <Sidebar />
    <div class="app-main">
      <header class="topbar">
        <div>
          <p class="eyebrow">Enterprise Control Plane</p>
          <h1>Microservice Operations Portal</h1>
        </div>
        <div class="topbar-actions">
          <button class="ghost" type="button">Launch Incident Runbook</button>
          <button class="primary" type="button">Create Deployment</button>
        </div>
      </header>
      <main class="content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import Sidebar from '@/components/Sidebar.vue';
</script>
