<script>
export default {
    data: function () {
        return {
            errormsg: null,
            successmsg: null,
            loading: false,
            profile: [],
            profileID: null,
            found: false,
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
            this.successmsg = null
            this.found = false;
            try {
                if (!this.profileID) {
                    this.errormsg = "The profile ID is empty"
                } else {
                    let response = await this.$axios.get("/users/" + this.profileID + "/profile", {
                        headers: {
                            Authorization: 'Bearer ' + this.token
                        }
                    });
                    this.profile = response.data;
                    this.found = true;
                }
            } catch (e) {
                if (e.response && e.response.status == '403') {
                    this.errormsg = "The owner of this profile banned you.";
                } if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no user with this User ID." + this.profileID;
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
        this.profileID = this.$route.params.id;
        this.token = localStorage.token;
        this.refresh();
    },
    watch: {
        token(token) {
            if (token) {
                localStorage.token = token;
            }
        }
    },
    created() {
        this.$watch(
            () => this.$route.params,
            (toParams, previousParams) => {
                if (toParams.id != previousParams.id) {
                    this.profileID = toParams.id;
                    this.refresh();
                }
            }
        )
    }
}
</script>
<template>
    <div
        class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
        <h1 class="h2">User Profile</h1>
        <div class="btn-toolbar mb-2 mb-md-0">
            <div class="btn-group me-2">
                <button type="button" class="btn btn-sm btn-outline-primary" @click="refresh">
                    Refresh
                </button>
            </div>
        </div>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    <SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>

    <div v-if="token">
        <div v-if="found && !loading">
            <div class="mb-3">
                <label>

                </label>
            </div>

            <button type="button" class="btn btn-sm btn-primary" @click="uploadPhoto">
                Upload Photo
            </button>
        </div>
    </div>
    <div v-else>
        <h2>You can't do this without being authenticated!</h2>
    </div>
    <LoadingSpinner :loading="loading"></LoadingSpinner>
</template>
<style scoped>
.card {
    margin-bottom: 20px;
}
</style>