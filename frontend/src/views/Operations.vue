<template>
  <section class="section">
    <div class="section-heading">
      <div>
        <p class="eyebrow">Operations</p>
        <h2>Service Actions</h2>
      </div>
      <button class="primary" type="button">Open Change Request</button>
    </div>

    <div class="card-grid">
      <article class="card">
        <h3>User Service</h3>
        <p>Manage user profiles, active status, and audit-safe updates.</p>
        <ul>
          <li>POST /api/v1/users</li>
          <li>GET /api/v1/users/{id}</li>
          <li>PUT /api/v1/users/{id}</li>
        </ul>
      </article>
      <article class="card">
        <h3>Order Service</h3>
        <p>Coordinate order lifecycles and validate users via service token.</p>
        <ul>
          <li>POST /api/v1/orders</li>
          <li>GET /api/v1/orders/{id}</li>
          <li>PUT /api/v1/orders/{id}</li>
        </ul>
      </article>
      <article class="card">
        <h3>Migrations</h3>
        <p>Automated migrations run on startup, with manual triggers ready.</p>
        <ul>
          <li>make migrate-user</li>
          <li>make migrate-order</li>
        </ul>
      </article>
    </div>
  </section>
</template>

<script setup>
</script>
