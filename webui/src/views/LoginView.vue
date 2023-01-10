<script>
export default {
    components: {},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			loginInfo: {
                Username: ""
            },
            token: null
		}
	},
	methods: {
		async login() {
            this.loading = true;
            this.errormsg = null;
            try {
                let response = await this.$axios.post("/session", this.loginInfo);
                this.token = response.data.UserID;
            } catch (e) {
                if (e.response && e.response.status == '400') {
                    this.errormsg = "The username is empty";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured";
                } else {
                    this.errormsg = e.toString();
                }
            }
            this.loading = false;
        },
        async logout() {
            this.token = null
            localStorage.removeItem("token")
        }
    },
    mounted() {
        this.token = localStorage.token;
    },
    watch: {
        token(token) {
            if (token) {
                localStorage.token = token;
            }
        }
    }
}
</script>

<template>
	<div v-show="!token">
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Login</h1>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

        <LoadingSpinner :loading="loading"></LoadingSpinner>

        <div v-if="!loading">
            <div class="mb-3">
				<label for="username" class="form-label">Username</label>
				<input type="text" class="form-control" id="username" v-model="loginInfo.Username">
			</div>
        </div>

        <button type="button" class="btn btn-sm btn-primary" @click="login">
			Login
		</button>
	</div>

    <div v-show="token">
        <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Logout</h1>
		</div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

        <LoadingSpinner :loading="loading"></LoadingSpinner>

        <button type="button" class="btn btn-sm btn-primary" @click="logout">
			Logout
		</button>
	</div>


</template>

<style>
</style>
