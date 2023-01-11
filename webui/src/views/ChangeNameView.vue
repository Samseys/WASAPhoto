<script>
export default {
    components: {},
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
                if (this.Username.Username != "") {
                    await this.$axios.put("/users/" + localStorage.token + "/name", this.Username, {
                        headers: {
                            Authorization: 'Bearer ' + token
                        }
                    });
                    this.success = "Name changed successfully."
                } else {
                    this.errormsg = "The username is empty."
                }
            } catch (e) {
                if (e.response && e.response.status == '401') {
                    this.errormsg = "You have to be authenticated to change name."
                } else if (e.response && e.response.status == '409') {
                    this.errormsg = "Another user has the same username.";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
            this.loading = false;
        },
        async logout() {
            this.loading = true;
            this.token = null
            localStorage.removeItem("token")
            this.loading = false;
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
                <input type="text" class="form-control" id="username" v-model="Username.Username">
            </div>

            <button type="button" class="btn btn-sm btn-primary" @click="changeName">
                Change name
            </button>

        </div>
    </div>
    <div v-else>
        <h1>You can't do this without being authenticated!</h1>
    </div>
    <LoadingSpinner :loading="loading"></LoadingSpinner>

</template>

<style>

</style>
