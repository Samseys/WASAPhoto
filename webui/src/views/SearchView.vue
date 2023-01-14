<script>export default {

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