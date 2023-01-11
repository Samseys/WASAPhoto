import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ChangeNameView from '../views/ChangeNameView.vue'
import PhotoUploadView from '../views/PhotoUploadView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: '/', component: HomeView },
		{ path: '/login', component: LoginView },
		{ path: '/profile/changename', component: ChangeNameView },
		{ path: '/photos/upload', component: PhotoUploadView },
	]
})

export default router
