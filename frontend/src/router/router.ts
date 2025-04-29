import { createRouter, createWebHistory } from 'vue-router'
const Index = () => import('../views/Index.vue');
const Simulation = () => import('../views/Simulation.vue');

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'index',
      component: Index,
    },
    {
      path: '/simulation',
      name: 'simulation',
      component: Simulation,
      props: true,
    }
  ],
})

export default router
