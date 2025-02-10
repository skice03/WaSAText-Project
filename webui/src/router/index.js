import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import ChatsView from '../views/ChatsView.vue'


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/profile', component: ProfileView},
		{path: '/home', component: HomeView},
	//	{path: '/some/:id/link', component: HomeView},
	]
})

export default router
