<script>
import Photo from '../components/Photo.vue'
export default {
	components: {
		Photo
	},
	data: function () {
		return {
			errormsg: null,
			loading: false,
			photos: [],
			token: null
		}
	},
	methods: {
		load() {
			return load
		},
		async refresh() {
			if (!this.token)
				return;
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/stream", {
					headers: {
						Authorization: 'Bearer ' + this.token
					}
				});
				this.photos = response.data.Photos;
			} catch (e) {
				if (e.response && e.response.status == '500') {
					this.errormsg = "An internal error has occured.";
				} else {
					this.errormsg = e.toString();
				}
			}
			this.loading = false;
		}
	},
	mounted() {
		this.token = localStorage.token;
		this.refresh();
	},
}
</script>
<template>
	<div>
		<div v-if="!loading">
			<div
				class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
				<h1 class="h2">Stream</h1>
				<div class="btn-toolbar mb-2 mb-md-0">
					<div class="btn-group me-2">
						<button type="button" class="btn btn-sm btn-outline-primary" @click="refresh">
							Refresh
						</button>
					</div>
				</div>
			</div>
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

			<div v-if="token">
				<div v-if="this.photos && this.photos.length != 0">
					<Photo :photoID=photo v-for="photo in this.photos" :key="photo"></Photo>
				</div>
				<div class=" card" v-else>
					<div class="card-body">
						<p class="card-text">You don't have any photos in your stream.</p>
					</div>
				</div>
			</div>
			<div class="card" v-else>
				<div class="card-body">
					<p class="card-text">
						You can't access profiles without being authenticated!
					</p>
				</div>
			</div>
		</div>
		<LoadingSpinner :loading="loading"></LoadingSpinner>
	</div>
</template>
<style scoped>
.card {
	margin-bottom: 20px;
}
</style>
