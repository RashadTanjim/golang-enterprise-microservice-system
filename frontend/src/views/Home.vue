<template>
  <section class="landing-hero">
    <div class="hero-copy">
      <p class="hero-kicker">Go + Vue Reference Platform</p>
      <h2>
        Build dependable microservices faster with an
        <span class="gradient-text">enterprise-ready foundation</span>
      </h2>
      <p class="hero-lead">
        Enterprise Microservice System is a production-focused sample that demonstrates
        clean service boundaries, secure API patterns, resilient communication, and
        observability-first operations in one cohesive stack.
      </p>
      <div class="hero-dynamic" aria-live="polite">
        <span class="pulse-dot"></span>
        <span>{{ activePhrase }}</span>
      </div>
      <div class="hero-actions">
        <router-link class="primary cta-link" to="/operations">Explore Operations</router-link>
        <router-link class="ghost cta-link" to="/docs">Read API Docs</router-link>
      </div>
    </div>

    <div class="hero-visual" role="presentation">
      <div class="orbit orbit-a"></div>
      <div class="orbit orbit-b"></div>
      <div class="orbit orbit-c"></div>
      <div class="visual-panel">
        <p class="visual-label">Runtime Focus</p>
        <h3>Reliability + Delivery</h3>
        <ul>
          <li>JWT-secured APIs</li>
          <li>Circuit breaker on service calls</li>
          <li>Health, metrics, tracing-ready design</li>
        </ul>
      </div>
    </div>
  </section>

  <section class="kpi-grid reveal-block" :ref="setRevealRef">
    <article v-for="item in animatedStats" :key="item.label" class="kpi-card">
      <p class="kpi-value">{{ item.value }}{{ item.suffix }}</p>
      <p class="kpi-label">{{ item.label }}</p>
      <p class="kpi-note">{{ item.note }}</p>
    </article>
  </section>

  <section class="section landing-section reveal-block" :ref="setRevealRef">
    <div class="section-heading landing-heading">
      <div>
        <p class="eyebrow">Project Purpose</p>
        <h3>Why this frontend exists</h3>
      </div>
    </div>
    <div class="purpose-grid">
      <article class="purpose-card">
        <h4>Unified visibility</h4>
        <p>
          Present service operations, endpoint contracts, and monitoring paths in one UI so
          teams can align product and platform decisions quickly.
        </p>
      </article>
      <article class="purpose-card">
        <h4>Production reference</h4>
        <p>
          Demonstrate practical architecture patterns teams can copy into real projects:
          authentication, migrations, gateway routing, and resilient service calls.
        </p>
      </article>
      <article class="purpose-card">
        <h4>Developer acceleration</h4>
        <p>
          Reduce setup ambiguity with explicit docs, endpoint examples, and a clear path from
          local development to Docker, Kubernetes, and Netlify.
        </p>
      </article>
    </div>
  </section>

  <section class="section landing-section reveal-block" :ref="setRevealRef">
    <div class="section-heading landing-heading">
      <div>
        <p class="eyebrow">System Blueprint</p>
        <h3>How the platform is structured</h3>
      </div>
    </div>
    <div class="blueprint-grid">
      <article class="blueprint-card">
        <p class="blueprint-step">1</p>
        <h4>Domain Services</h4>
        <p>User, order, and audit services with clear ownership and isolated responsibilities.</p>
      </article>
      <article class="blueprint-card">
        <p class="blueprint-step">2</p>
        <h4>Reliability Layer</h4>
        <p>Circuit breaker, timeouts, and defensive defaults for safer inter-service communication.</p>
      </article>
      <article class="blueprint-card">
        <p class="blueprint-step">3</p>
        <h4>Data + Caching</h4>
        <p>PostgreSQL and Redis integration with migration workflows and operational checks.</p>
      </article>
      <article class="blueprint-card">
        <p class="blueprint-step">4</p>
        <h4>Operational Surface</h4>
        <p>Gateway routes, health probes, metrics endpoints, docs, and deployment playbooks.</p>
      </article>
    </div>
  </section>

  <section class="section landing-section reveal-block" :ref="setRevealRef">
    <div class="section-heading landing-heading">
      <div>
        <p class="eyebrow">Quick Start</p>
        <h3>Backend in one command</h3>
      </div>
    </div>
    <div class="command-panel">
      <code>make docker-prod-up</code>
      <p>
        Start production-like backend services with `docker-compose.prod.yml`, then point
        Netlify frontend variables to your API domain.
      </p>
    </div>
  </section>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';

const phrases = [
  'Security-first APIs with JWT authorization',
  'Observability built in from day one',
  'Deployment-ready from laptop to cluster'
];

const activePhraseIndex = ref(0);
const statValues = ref([0, 0, 0, 0]);
const revealElements = [];

let phraseIntervalId = null;
let revealObserver = null;

const statTargets = [
  { label: 'Core Services', target: 4, suffix: '', note: 'User, order, audit, migration' },
  { label: 'Deployment Modes', target: 3, suffix: '', note: 'Docker, Kubernetes, Netlify frontend' },
  { label: 'SLO Mindset', target: 99.95, suffix: '%', note: 'High-availability operational target' },
  { label: 'Request Tracing', target: 100, suffix: '%', note: 'Structured logs with request IDs' }
];

const activePhrase = computed(() => phrases[activePhraseIndex.value]);

const animatedStats = computed(() =>
  statTargets.map((item, index) => ({
    ...item,
    value: item.target % 1 === 0 ? Math.round(statValues.value[index]) : statValues.value[index].toFixed(2)
  }))
);

const animateStats = () => {
  const durationMs = 1400;
  const startTime = performance.now();

  const tick = (now) => {
    const progress = Math.min((now - startTime) / durationMs, 1);
    const eased = 1 - Math.pow(1 - progress, 3);

    statValues.value = statTargets.map((item) => item.target * eased);

    if (progress < 1) {
      requestAnimationFrame(tick);
    }
  };

  requestAnimationFrame(tick);
};

const setRevealRef = (el) => {
  if (el) {
    revealElements.push(el);
  }
};

onMounted(() => {
  animateStats();

  phraseIntervalId = window.setInterval(() => {
    activePhraseIndex.value = (activePhraseIndex.value + 1) % phrases.length;
  }, 2400);

  if (typeof IntersectionObserver === 'undefined') {
    revealElements.forEach((element) => element.classList.add('visible'));
    return;
  }

  revealObserver = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible');
          revealObserver.unobserve(entry.target);
        }
      });
    },
    { threshold: 0.2 }
  );

  revealElements.forEach((element) => revealObserver.observe(element));
});

onBeforeUnmount(() => {
  if (phraseIntervalId) {
    window.clearInterval(phraseIntervalId);
  }

  if (revealObserver) {
    revealObserver.disconnect();
  }
});
</script>

<style scoped>
.landing-hero {
  display: grid;
  grid-template-columns: 1.1fr 0.9fr;
  gap: 28px;
  padding: 34px;
  border-radius: 26px;
  background:
    radial-gradient(circle at 12% 10%, rgba(34, 139, 230, 0.16), transparent 34%),
    radial-gradient(circle at 90% 90%, rgba(241, 143, 1, 0.2), transparent 38%),
    linear-gradient(135deg, #ffffff 0%, #f4f9ff 56%, #fff5e6 100%);
  border: 1px solid rgba(15, 92, 122, 0.16);
  box-shadow: 0 26px 52px rgba(8, 37, 52, 0.12);
  overflow: hidden;
}

.hero-copy {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.hero-kicker {
  font-size: 12px;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: #4b7289;
  font-weight: 600;
}

.hero-copy h2 {
  font-family: 'Sora', sans-serif;
  font-size: clamp(30px, 4vw, 45px);
  line-height: 1.15;
  max-width: 20ch;
}

.gradient-text {
  background: linear-gradient(95deg, #0f5c7a 0%, #1f8a5f 60%, #f18f01 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.hero-lead {
  color: #334a59;
  line-height: 1.75;
  max-width: 62ch;
}

.hero-dynamic {
  margin-top: 8px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 999px;
  background: rgba(15, 92, 122, 0.09);
  color: #0f5c7a;
  font-weight: 600;
  width: fit-content;
}

.pulse-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: #f18f01;
  box-shadow: 0 0 0 rgba(241, 143, 1, 0.45);
  animation: pulse 1.6s ease-in-out infinite;
}

.hero-actions {
  margin-top: 10px;
  display: flex;
  gap: 12px;
}

.cta-link {
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.hero-visual {
  position: relative;
  min-height: 280px;
  display: grid;
  place-items: center;
}

.orbit {
  position: absolute;
  border-radius: 50%;
  border: 1px dashed rgba(15, 92, 122, 0.24);
  animation: spin 18s linear infinite;
}

.orbit-a {
  width: 250px;
  height: 250px;
}

.orbit-b {
  width: 190px;
  height: 190px;
  animation-duration: 14s;
  animation-direction: reverse;
}

.orbit-c {
  width: 132px;
  height: 132px;
  animation-duration: 9s;
}

.visual-panel {
  z-index: 1;
  width: min(100%, 300px);
  background: rgba(255, 255, 255, 0.86);
  border: 1px solid rgba(11, 31, 42, 0.1);
  backdrop-filter: blur(8px);
  border-radius: 18px;
  padding: 18px;
  box-shadow: 0 18px 40px rgba(11, 31, 42, 0.14);
}

.visual-label {
  font-size: 11px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: #607b8b;
}

.visual-panel h3 {
  margin-top: 7px;
  font-family: 'Sora', sans-serif;
  font-size: 20px;
}

.visual-panel ul {
  margin-top: 12px;
  list-style: none;
  display: grid;
  gap: 9px;
  color: #355062;
}

.visual-panel li {
  padding-left: 14px;
  position: relative;
}

.visual-panel li::before {
  content: '';
  position: absolute;
  left: 0;
  top: 8px;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #1f8a5f;
}

.kpi-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.kpi-card {
  background: #ffffff;
  border: 1px solid rgba(11, 31, 42, 0.1);
  border-radius: 18px;
  padding: 20px;
  box-shadow: 0 16px 30px rgba(11, 31, 42, 0.08);
}

.kpi-value {
  font-family: 'Sora', sans-serif;
  font-size: 33px;
  color: #0f5c7a;
}

.kpi-label {
  margin-top: 7px;
  font-weight: 600;
}

.kpi-note {
  margin-top: 6px;
  font-size: 13px;
  color: #5c6b78;
}

.landing-section {
  border: 1px solid rgba(11, 31, 42, 0.08);
}

.landing-heading h3 {
  font-family: 'Sora', sans-serif;
  font-size: 30px;
  margin-top: 7px;
}

.purpose-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.purpose-card {
  border-radius: 16px;
  border: 1px solid #d7e4ee;
  background: linear-gradient(165deg, #ffffff, #f6fbff);
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.purpose-card h4,
.blueprint-card h4 {
  font-family: 'Sora', sans-serif;
  font-size: 19px;
}

.purpose-card p,
.blueprint-card p {
  color: #4e6170;
  line-height: 1.65;
}

.blueprint-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.blueprint-card {
  background: linear-gradient(140deg, #f5faf6 0%, #ffffff 60%);
  border: 1px solid #d5e7dd;
  border-radius: 16px;
  padding: 20px;
}

.blueprint-step {
  font-family: 'Sora', sans-serif;
  color: #1f8a5f;
  font-size: 26px;
  margin-bottom: 8px;
}

.command-panel {
  border-radius: 18px;
  background: linear-gradient(125deg, #0b1f2a 0%, #113a4d 60%, #0f5c7a 100%);
  color: #e8f2fa;
  border: 1px solid rgba(255, 255, 255, 0.14);
  padding: 24px;
}

.command-panel code {
  display: inline-block;
  font-size: 24px;
  font-weight: 700;
  letter-spacing: 0.03em;
  color: #ffe2b7;
}

.command-panel p {
  margin-top: 10px;
  line-height: 1.7;
  max-width: 72ch;
}

.reveal-block {
  opacity: 0;
  transform: translateY(14px);
  transition: opacity 0.6s ease, transform 0.6s ease;
}

.reveal-block.visible {
  opacity: 1;
  transform: translateY(0);
}

@keyframes pulse {
  0%,
  100% {
    box-shadow: 0 0 0 0 rgba(241, 143, 1, 0.45);
  }
  50% {
    box-shadow: 0 0 0 8px rgba(241, 143, 1, 0);
  }
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 1200px) {
  .kpi-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .purpose-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 900px) {
  .landing-hero {
    grid-template-columns: 1fr;
    padding: 24px;
  }

  .hero-copy h2 {
    font-size: 32px;
  }

  .blueprint-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 680px) {
  .hero-actions {
    flex-direction: column;
  }

  .cta-link {
    width: 100%;
  }

  .kpi-grid,
  .purpose-grid {
    grid-template-columns: 1fr;
  }

  .kpi-value {
    font-size: 30px;
  }

  .landing-heading h3 {
    font-size: 25px;
  }

  .command-panel code {
    font-size: 20px;
  }
}
</style>
