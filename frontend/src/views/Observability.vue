<template>
  <section class="section">
    <div class="section-heading">
      <div>
        <p class="eyebrow">Observability</p>
        <h2>Metrics, Logs, and Health</h2>
      </div>
      <button class="ghost" type="button">Export Metrics</button>
    </div>

    <div class="card-grid">
      <article class="card">
        <h3>Health Endpoints</h3>
        <p>Check readiness and liveness via the gateway.</p>
        <ul>
          <li>/health/user</li>
          <li>/health/order</li>
        </ul>
      </article>
      <article class="card">
        <h3>Metrics Endpoints</h3>
        <p>Prometheus-compatible metrics via gateway routes.</p>
        <ul>
          <li>/metrics/user</li>
          <li>/metrics/order</li>
        </ul>
      </article>
      <article class="card">
        <h3>Logging</h3>
        <p>Structured JSON logs with request IDs for traceability.</p>
        <ul>
          <li>Zap logging</li>
          <li>Request ID middleware</li>
        </ul>
      </article>
    </div>
  </section>
</template>

<script setup>
</script>
