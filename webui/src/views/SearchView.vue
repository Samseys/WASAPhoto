<script>

export default {
    data: function () {
        return {
            username: "",
            errormsg: null,
            successmsg: null,
            token: null
        }
    },
    methods: {
        async searchUsername() {
            try {
                if (username != null) {
                    let response = await this.$axios.get("/userids/" + this.username);
                    let responseID = response.data.UserID;
                    this.$router.push('/profile/' + responseID)
                }
            } catch (e) {
                if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no user with this Username " + this.username + ".";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
        }
    },
    mounted() {
        this.token = localStorage.token;
    }
}
</script>
<template>
    <div
        class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
        <h1 class="h2">Search</h1>
    </div>
    <SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    <div v-if="token">
        <div class="mb-3">
            <label for="username" class="form-label">Username</label>
            <input type="text" class="form-control" id="username" v-model="this.username" />
        </div>

        <button type="button" class="btn btn-sm btn-primary" @click="searchUsername()">
            Search
        </button>
    </div>
    <div v-else>
        <div class="card">
            <div class="card-body">
                <p class="card-text">
                    You can't search users without being authenticated!
                </p>
            </div>
        </div>
    </div>
    <LoadingSpinner :loading="loading"></LoadingSpinner>
</template>

<style>
.card {
    margin-bottom: 20px;
}
</style>