import { createRouter, createWebHistory } from 'vue-router';
import Home from '@/views/Home.vue';
import Operations from '@/views/Operations.vue';
import Observability from '@/views/Observability.vue';
import Docs from '@/views/Docs.vue';

const routes = [
  { path: '/', name: 'home', component: Home },
  { path: '/operations', name: 'operations', component: Operations },
  { path: '/observability', name: 'observability', component: Observability },
  { path: '/docs', name: 'docs', component: Docs }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

export default router;
