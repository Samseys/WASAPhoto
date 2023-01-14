<script>
export default {
    components: {},
    data: function () {
        return {
            errormsg: null,
            successmsg: null,
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
            this.successmsg = null;
            try {
                if (this.loginInfo.Username != "") {
                    let response = await this.$axios.post("/session", this.loginInfo);
                    this.token = response.data.UserID;
                    localStorage.setItem("username", this.loginInfo.Username);
                    this.successmsg = "Logged in with User ID " + this.token;
                } else {
                    this.errormsg = "The username is empty."
                }
            } catch (e) {
                if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
            this.loading = false;
        },
        logout() {
            this.error = null;
            this.loading = true;
            this.token = null;
            localStorage.removeItem("token");
            localStorage.removeItem("username");
            this.loginInfo.Username = ""
            this.successmsg = "Logged out successfully.";
            this.loading = false;
        }
    },
    mounted() {
        this.token = localStorage.token;
        this.loginInfo.Username = localStorage.username;
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
    <div v-if="!loading">
        <div
            class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            <div v-if="token">
                <h1 class="h2">Logout</h1>
            </div>
            <div v-else>
                <h1 class="h2">Login</h1>
            </div>
        </div>
        <SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        <div v-if="!token">
            <div class="mb-3">
                <label for="username" class="form-label">Username</label>
                <input type="text" class="form-control" id="username" v-model="loginInfo.Username" />
            </div>
        </div>
        <div v-else>
            <div class="card">
                <div class="card-body">
                    <p class="card-text">
                        Username: {{ this.loginInfo.Username }}<br />
                        ID: {{ this.token }}

                    </p>
                </div>
            </div>
        </div>

        <div v-if="token">
            <button type="button" class="btn btn-sm btn-primary" @click="logout">
                Logout
            </button>
        </div>
        <div v-else>
            <button type="button" class="btn btn-sm btn-primary" @click="login">
                Login
            </button>
        </div>
    </div>

    <LoadingSpinner :loading="loading"></LoadingSpinner>
</template>

<style>
.card {
    margin-bottom: 20px;
}
</style>
