<script>
import Photo from '../components/Photo.vue'
export default {
    components: {
        Photo
    },
    data: function () {
        return {
            errormsg: null,
            successmsg: null,
            header: "User Profile",
            loading: false,
            banned: false,
            profile: {
                UserID: Number,
                Username: String,
                Photos: [],
                Followers: [{
                    UserID: Number,
                    Username: String
                }],
                Following: [{
                    UserID: Number,
                    Username: String
                }],
                Banned: [{
                    UserID: Number,
                    Username: String
                }],
                BannedBy: [{
                    UserID: Number,
                    Username: String
                }],
            },
            profileID: null,
            found: false,
            token: null,
        }
    },
    methods: {
        load() {
            return load
        },
        deletePhoto(photo) {
            let index = this.profile.Photos.indexOf(photo);
            if (index != -1) {
                this.profile.Photos.splice(index, 1);
                this.successmsg = "Photo deleted.";
            } else {
                this.refresh();
            }
        },
        async refresh() {
            if (!this.token)
                return;
            this.loading = true;
            this.errormsg = null;
            this.successmsg = null;
            this.found = false;
            this.header = "User Profile";
            this.banned = false;
            try {
                if (!this.profileID) {
                    this.errormsg = "The profile ID is empty"
                } else {
                    let response = await this.$axios.get("/users/" + this.profileID + "/profile");
                    this.profile = response.data;
                    this.header = this.profile.Username;
                    this.found = true;
                }
            } catch (e) {
                if (e.response && e.response.status == '403') {
                    this.errormsg = "The owner of this profile banned you.";
                    this.banned = true;
                } else if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no user with this User ID " + this.profileID + ".";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
            this.loading = false;
        },
        async follow() {
            this.errormsg = ""
            try {
                await this.$axios.put("/users/" + this.token + "/followed/" + this.profile.UserID, "");
                if (!this.profile.Followers) {
                    this.profile.Followers = [];
                }

                this.profile.Followers.push({ "UserID": this.token, "Username": localStorage.username })
            } catch (e) {
                if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no user with this id: " + this.profile.UserID + ".";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
        },
        async unfollow() {
            this.errormsg = ""
            try {
                await this.$axios.delete("/users/" + this.token + "/followed/" + this.profile.UserID);
                this.profile.Followers.splice(this.profile.Followers.findIndex(item => item.UserID == this.token), 1)
            } catch (e) {
                if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no user with this id: " + this.profile.UserID + ".";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
        },
        async ban() {
            this.errormsg = ""
            try {
                await this.$axios.put("/users/" + this.token + "/banned/" + this.profile.UserID, "");

                if (!this.profile.BannedBy) {
                    this.profile.BannedBy = [];
                }

                this.profile.BannedBy.push({ "UserID": this.token, "Username": localStorage.username })
            } catch (e) {
                if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no user with this id: " + this.profile.UserID + ".";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
        },
        async unban() {
            this.errormsg = ""
            try {
                await this.$axios.delete("/users/" + this.token + "/banned/" + this.profile.UserID);
                this.profile.BannedBy.splice(this.profile.BannedBy.findIndex(item => item.UserID == this.token), 1)
            } catch (e) {
                if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no user with this id: " + this.profile.UserID + ".";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
        }
    },
    mounted() {
        if (this.$route.params.id == "me") {
            this.profileID = localStorage.token;
        } else {
            this.profileID = this.$route.params.id;
        }
        this.token = localStorage.token;
        this.refresh();
    },
    created() {
        this.$watch(
            () => this.$route.params,
            (toParams) => {
                let id;
                if (toParams.id == "me") {
                    id = localStorage.token;
                } else {
                    id = toParams.id;
                }

                if (id != this.profileID) {
                    this.profileID = id;
                    this.refresh();
                }
            }
        )
    }
}
</script>
<template>
    <div>
        <div v-if="!loading">
            <div
                class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
                <h1 class="h2">{{ this.header }}</h1>
                <div class="btn-toolbar mb-2 mb-md-0">
                    <div class="btn-group me-2" v-if="this.token && !this.banned && this.profile.UserID != this.token">
                        <button type="button" class="btn btn-sm btn-outline-primary" @click="unfollow"
                            v-if="(this.profile.Followers ?? []).some(follower => follower.UserID == this.token)">
                            Unfollow
                        </button>
                        <button type="button" class="btn btn-sm btn-outline-primary" @click="follow" v-else>
                            Follow
                        </button>
                        <button type="button" class="btn btn-sm btn-outline-primary" @click="unban"
                            v-if="(this.profile.BannedBy ?? []).some(bannedby => bannedby.UserID == this.token)">
                            Unban
                        </button>
                        <button type="button" class="btn btn-sm btn-outline-primary" @click="ban" v-else>
                            Ban
                        </button>
                    </div>
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
                <div v-if="found">
                    <div class="card">
                        <div class="card-header">User Info</div>
                        <div class="card-body">
                            <p class="card-text">Number of uploaded photos: {{ (this.profile.Photos ?? []).length }}</p>
                            <p class="card-text">
                                Following: <span v-for="(followed, index) in (this.profile.Following ?? [])"
                                    :key="followed.UserID">
                                    <RouterLink :to="'/profile/' + followed.UserID">
                                        {{ followed.Username }}
                                    </RouterLink>
                                    <span v-if="index != Object.keys(this.profile.Following ?? []).length - 1">, </span>
                                </span>
                            </p>
                            <p class="card-text">
                                Followers: <span v-for="(follower, index) in (this.profile.Followers ?? [])"
                                    :key="follower.UserID">
                                    <RouterLink :to="'/profile/' + follower.UserID">
                                        {{ follower.Username }}
                                    </RouterLink>
                                    <span v-if="index != Object.keys(this.profile.Followers ?? []).length - 1">, </span>
                                </span>
                            </p>
                        </div>
                    </div>
                    <div v-if=this.profile.Photos>
                        <Photo :photoID=photo v-for="photo in this.profile.Photos" :key="photo"
                            @delete-photo="deletePhoto(photo)"></Photo>
                    </div>
                    <div class="card" v-else>
                        <div class="card-body">
                            <p class="card-text">This user hasn't uploaded any photo yet.</p>
                        </div>
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
