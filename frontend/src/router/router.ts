import { createRouter, createWebHistory } from 'vue-router'
const Index = () => import('../views/Index.vue');
const Simulation = () => import('../views/Simulation.vue');
const Transcripts = () => import('../views/Transcripts.vue');
const GetTranscript = () => import('../views/GetTranscript.vue');
const Login = () => import('../views/Login.vue');
const UserGuide = () => import('../views/UserGuide.vue');

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
    },
    {
      path: '/get-transcripts',
      name: 'transcripts',
      component: Transcripts,
    },
    {
      path: '/transcripts/:id',
      name: 'transcripts:id',
      component: GetTranscript,
      props: true,
    },
    {
      path: '/login',
      name: 'login',
      component: Login,
    },
    {
      path: '/user-guide',
      name: 'user-guide',
      component: UserGuide,
    }
  ],
})

router.beforeEach((to, from, next) => {
  const publicPages = ['/login']
  const authRequired = !publicPages.includes(to.path)
  const userToken = localStorage.getItem('userToken')

  if (authRequired && !userToken) {
    next('/login')
  } else {
    next()
  }
})

export default router
