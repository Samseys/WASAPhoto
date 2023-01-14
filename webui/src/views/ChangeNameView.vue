<script>
export default {
    data: function () {
        return {
            errormsg: null,
            successmsg: null,
            loading: false,
            Username: {
                Username: ""
            },
            token: null
        }
    },
    methods: {
        async changeName() {
            this.loading = true;
            this.errormsg = null;
            this.successmsg = null
            try {
                if (this.Username.Username == "") {
                    this.errormsg = "The username is empty.";
                } else if (this.Username.Username == localStorage.username) {
                    this.errormsg = "You already have this name.";
                } else {
                    await this.$axios.put("/users/" + this.token + "/name", this.Username, {
                        headers: {
                            Authorization: 'Bearer ' + this.token
                        }
                    });
                    this.successmsg = "Name changed successfully.";
                    localStorage.setItem("username", this.Username.Username)
                }
            } catch (e) {
                if (e.response && e.response.status == '409') {
                    this.errormsg = "Another user has the same username.";
                } else if (e.response && e.response.status == '500') {
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
    }
}
</script>

<template>
    <div
        class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
        <h1 class="h2">Change Name</h1>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    <SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>

    <div v-if="token">
        <div v-if="!loading">
            <div class="mb-3">
                <label for="username" class="form-label">New Name</label>
                <input type="text" class="form-control" id="username" v-model="Username.Username" />
            </div>

            <button type="button" class="btn btn-sm btn-primary" @click="changeName">
                Change name
            </button>
        </div>
    </div>
    <div v-else>
        <div class="card">
            <div class="card-body">
                <p class="card-text">
                    You can't change name without being authenticated!
                </p>
            </div>
        </div>
    </div>
    <LoadingSpinner :loading="loading"></LoadingSpinner>
</template>

<style>

</style>
