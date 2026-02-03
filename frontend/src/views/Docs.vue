<template>
  <section class="section">
    <div class="section-heading">
      <div>
        <p class="eyebrow">Docs & Access</p>
        <h2>Developer Enablement</h2>
      </div>
      <button class="primary" type="button">Open Swagger</button>
    </div>

    <div class="card-grid">
      <article class="card">
        <h3>Swagger UI</h3>
        <p>Live API references for each service.</p>
        <ul>
          <li>/swagger/user/</li>
          <li>/swagger/order/</li>
        </ul>
      </article>
      <article class="card">
        <h3>Postman Collections</h3>
        <p>Automated collection with assertions for key flows.</p>
        <ul>
          <li>docs/collections/enterprise-microservice-system.postman_collection.json</li>
          <li>docs/collections/enterprise-microservice-system.postman_environment.json</li>
        </ul>
      </article>
      <article class="card">
        <h3>Security</h3>
        <p>JWT tokens with role-based authorization.</p>
        <ul>
          <li>POST /api/v1/auth/token</li>
          <li>Authorization: Bearer &lt;token&gt;</li>
        </ul>
      </article>
    </div>
  </section>
</template>

<script setup>
</script>
