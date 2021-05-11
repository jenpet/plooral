import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import Board from '../views/Board.vue'
import OrganizationList from "@/views/OrganizationList.vue";
import BoardList from "@/views/BoardList.vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: '/orgs/:orgSlug/boards/:boardSlug',
    name: 'Board',
    component: Board
  },
  {
    path: '/orgs',
    alias: '/',
    name: 'OrganizationList',
    component: OrganizationList
  },
  {
    path: '/orgs/:orgSlug/boards',
    name: 'BoardList',
    component: BoardList
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
