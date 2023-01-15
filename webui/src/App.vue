<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
	data: function () {
		return {
			loaded: false,
			token: null
		}
	},
	methods: {
		async checkLogged() {
			try {
				let username = localStorage.username;
				if (username != null) {
					let response = await this.$axios.get("/userids/" + username);
					let responseID = response.data.UserID;
					if (responseID != parseInt(localStorage.getItem("token"))) {
						localStorage.removeItem("token");
						localStorage.removeItem("username");
					}
				}
			} catch (e) {
				localStorage.removeItem("token");
				localStorage.removeItem("username");
			}
		},

		onLoginLogout() {
			this.$nextTick(() => {
				this.token = (localStorage.token ?? null);
			});
		}
	},
	async mounted() {
		await this.checkLogged();
		this.token = localStorage.token;
		this.loaded = true;
	}
}
</script>

<template>
	<div v-if="loaded">
		<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
			<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASA Photos</a>
			<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse"
				data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false"
				aria-label="Toggle navigation">
				<span class="navbar-toggler-icon"></span>
			</button>
			<div class="btn-group me-2">
				<RouterLink to="/login" class="nav-link">
					<button type="button" class="btn btn-sm btn-outline-primary">
						<span v-if="token">
							Logout
						</span>
						<span v-else>
							Login
						</span>
					</button>
				</RouterLink>
			</div>
		</header>

		<div class="container-fluid">
			<div class="row">
				<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
					<div class="position-sticky pt-3 sidebar-sticky">
						<h6
							class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
							<span>General</span>
						</h6>
						<ul class="nav flex-column">
							<li class="nav-item">
								<RouterLink to="/" class="nav-link">
									<svg class="feather">
										<use href="/feather-sprite-v4.29.0.svg#home" />
									</svg>
									Stream
								</RouterLink>
							</li>
							<li class="nav-item">
								<RouterLink to="/photos/upload" class="nav-link">
									<svg class="feather">
										<use href="/feather-sprite-v4.29.0.svg#upload" />
									</svg>
									Upload Photo
								</RouterLink>
							</li>
							<li class="nav-item">
								<RouterLink to="/search" class="nav-link">
									<svg class="feather">
										<use href="/feather-sprite-v4.29.0.svg#search" />
									</svg>
									Search by Username
								</RouterLink>
							</li>
						</ul>

						<h6
							class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
							<span>Profile</span>
						</h6>
						<ul class="nav flex-column">
							<li class="nav-item">
								<RouterLink to="/profile/me" class="nav-link">
									<svg class="feather">
										<use href="/feather-sprite-v4.29.0.svg#user" />
									</svg>
									My Profile
								</RouterLink>
							</li>
							<li class="nav-item">
								<RouterLink to="/profile/changename" class="nav-link">
									<svg class="feather">
										<use href="/feather-sprite-v4.29.0.svg#settings" />
									</svg>
									Change Name
								</RouterLink>
							</li>
						</ul>
					</div>
				</nav>
				<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
					<RouterView @on-login-logout="onLoginLogout()" />
				</main>
			</div>
		</div>
	</div>
</template>

<style>

</style>
