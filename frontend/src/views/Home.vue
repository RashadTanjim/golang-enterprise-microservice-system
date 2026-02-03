<template>
  <section class="hero">
    <div>
      <p class="eyebrow">System Health</p>
      <h2>Enterprise Microservice Portal</h2>
      <p class="lead">
        A unified command center for user and order services, built for
        resilience, observability, and secure operations.
      </p>
      <div class="hero-actions">
        <button class="primary" type="button">Generate Access Token</button>
        <button class="ghost" type="button">View Runbook</button>
      </div>
    </div>
    <div class="hero-panel">
      <MetricTile
        label="Deployments"
        value="12"
        footnote="Last 30 days"
      />
      <MetricTile
        label="SLO Target"
        value="99.95%"
        footnote="Current window"
      />
      <MetricTile
        label="Avg. Latency"
        value="128ms"
        footnote="P95 across services"
      />
    </div>
  </section>

  <section class="status-grid">
    <StatusCard
      title="User Service"
      subtitle="Identity & profile"
      :status="user.status"
      :message="user.message"
      :last-checked="user.checkedAt"
      endpoint="/api/v1/users"
      :variant="user.variant"
    />
    <StatusCard
      title="Order Service"
      subtitle="Transactions & fulfillment"
      :status="order.status"
      :message="order.message"
      :last-checked="order.checkedAt"
      endpoint="/api/v1/orders"
      :variant="order.variant"
    />
    <StatusCard
      title="Gateway"
      subtitle="Nginx routing layer"
      status="Healthy"
      message="Routing is active for API, metrics, and Swagger endpoints."
      :last-checked="gateway.checkedAt"
      endpoint="/api/v1/*"
      variant="positive"
    />
  </section>
</template>

<script setup>
import { onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { useHealthStore } from '@/stores/health';
import MetricTile from '@/components/MetricTile.vue';
import StatusCard from '@/components/StatusCard.vue';

const healthStore = useHealthStore();
const { user, order, gateway } = storeToRefs(healthStore);

onMounted(() => {
  healthStore.refresh();
});
</script>
