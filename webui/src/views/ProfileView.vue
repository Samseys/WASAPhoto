<script>
export default {
    data: function () {
        return {
            errormsg: null,
            successmsg: null,
            loading: false,
            profile: [],
            photos: [],
            comments: [],
            profileID: null,
            found: false,
            token: null
        }
    },
    methods: {
        load() {
            return load
        },
        async getImage(id) {
            let response = await this.$axios.get("/photos/" + id + "/file", {
                headers: {
                    Authorization: 'Bearer ' + this.token
                }, responseType: 'blob'
            });
            return response.data;

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
                    if (this.profile.Photos != null) {
                        for (const p of this.profile.Photos) {
                            var blob = await this.getImage(p.PhotoID)
                            this.photos[p.PhotoID] = window.URL.createObjectURL(blob)
                            this.comments[p.PhotoID] = " "
                        }
                    }
                    this.found = true;
                }
            } catch (e) {
                if (e.response && e.response.status == '403') {
                    this.errormsg = "The owner of this profile banned you.";
                } if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no user with this User ID " + this.profileID + ".";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
            this.loading = false;
        },
        getLikeQty(p) {
            if (p.Likes != null) {
                return p.Likes.length;
            } else {
                return 0;
            }
        },
        async like(p) {
            this.errormsg = ""
            try {
                await this.$axios.put("/photos/" + p.PhotoID + "/likes/" + this.token, "", {
                    headers: {
                        Authorization: 'Bearer ' + this.token
                    }
                });
                if (p.Likes == null) {
                    p.Likes = []
                }
                p.Likes.push({ "UserID": this.token, "Username": localStorage.getItem("Username") })
            } catch (e) {
                if (e.response && e.response.status == '403') {
                    this.errormsg = "The owner of this profile banned you.";
                } if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no photo with this id: " + p.PhotoID + ".";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
        },
        async unlike(p) {
            this.errormsg = ""
            try {
                await this.$axios.delete("/photos/" + p.PhotoID + "/likes/" + this.token, {
                    headers: {
                        Authorization: 'Bearer ' + this.token
                    }
                });
                p.Likes.splice(p.Likes.findIndex(item => item.UserID == this.token), 1)
            } catch (e) {
                if (e.response && e.response.status == '403') {
                    this.errormsg = "The owner of this profile banned you.";
                } if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no photo with this id: " + p.PhotoID + ".";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
        },
        async deletePhoto(p) {
            try {
                var response = await this.$axios.delete("/photos/" + p.PhotoID, {
                    headers: {
                        Authorization: 'Bearer ' + this.token
                    }
                });
                this.profile.Photos.splice(this.profile.Photos.findIndex(item => item.PhotoID == p.PhotoID), 1)
            } catch (e) {
                if (e.response && e.response.status == '403') {
                    this.errormsg = "The owner of this profile banned you.";
                } if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no photo with this id or the comment doesn't exist.";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
        },
        async postComment(p) {
            this.errormsg = ""
            if (this.comments[p.PhotoID] == null || this.comments[p.PhotoID] == "") {
                this.errormsg = "You can't post an empty comment";
            }
            try {
                var response = await this.$axios.post("/photos/" + p.PhotoID + "/comments", { Comment: this.comments[p.PhotoID] }, {
                    headers: {
                        Authorization: 'Bearer ' + this.token
                    }
                });
                if (p.Comments == null) {
                    p.Comments = []
                }
                p.Comments.push(response.data);
            } catch (e) {
                if (e.response && e.response.status == '403') {
                    this.errormsg = "The owner of this profile banned you.";
                } if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no photo with this id: " + p.PhotoID + ".";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
        },
        async deleteComment(p, c) {
            try {
                var response = await this.$axios.delete("/photos/" + p.PhotoID + "/comments/" + c.CommentID, {
                    headers: {
                        Authorization: 'Bearer ' + this.token
                    }
                });
                p.Comments.splice(p.Comments.findIndex(item => item.CommentID == c.CommentID), 1)
            } catch (e) {
                if (e.response && e.response.status == '403') {
                    this.errormsg = "The owner of this profile banned you.";
                } if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no photo with this id or the comment doesn't exist.";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
        }
    },
    mounted() {
        this.profileID = this.$route.params.id;
        this.token = localStorage.token;
        this.refresh();
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
            <div v-if="this.profile.Photos && this.profile.Photos.length != 0">
                <div class="card" v-for="p in this.profile.Photos">
                    <div class="card-header d-flex justify-content-between align-items-center">
                        <span>
                            <RouterLink :to="'/profile/' + p.Owner.UserID">
                                {{ p.Owner.Username }}
                            </RouterLink>
                            - {{ p.CreationDate }}
                        </span>
                        <button type="button" class="btn btn-danger" v-if="p.Owner.UserID == this.token"
                            @click="deletePhoto(p)">Delete
                            photo</button>
                    </div>
                    <div class="card-body">
                        <p class="card-text">
                        <div v-if=p.Comment>
                            {{ p.Comment }}<br /> <br />
                        </div>
                        </p>

                        <img :src=this.photos[p.PhotoID]>
                        <br />
                        Likes: {{ getLikeQty(p) }}<br /><br />
                        <div class="btn-toolbar">
                            <button type="button" class="btn btn-danger" @click="unlike(p)"
                                v-if="p.Likes != null && p.Likes.some(like => like.UserID == this.token)">Unlike</button>

                            <button type="button" class="btn btn-primary" @click="like(p)" v-else>Like</button>
                        </div>
                        <br />
                        <div class="card">
                            <div class="card-header">
                                Comment Section
                            </div>
                            <div class="card-body">
                                <div class="card" v-for="c in p.Comments">
                                    <div class="card-header d-flex justify-content-between align-items-center">
                                        <span>
                                            <RouterLink :to="'/profile/' + c.Owner.UserID">
                                                {{ c.Owner.Username }}
                                            </RouterLink>
                                            - {{ c.CreationDate }}
                                        </span>
                                        <button type="button" class="btn btn-danger small" @click="deleteComment(p, c)"
                                            v-if="c.Owner.UserID == this.token">Delete
                                            comment</button>

                                    </div>
                                    <div class="card-body">
                                        {{ c.Comment }}
                                    </div>
                                </div>
                                <div class="card">
                                    <div class="card-header">
                                        Post a comment
                                    </div>
                                    <div class="card-body d-flex justify-content-between align-items-center">
                                        <textarea v-model="this.comments[p.PhotoID]"></textarea><br />

                                        <button type="button" class="btn btn-primary"
                                            @click="postComment(p)">Comment</button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div v-else>
                <div class=" card">
                    <div class="card-body">
                        <p class="card-text">
                            This user hasn't uploaded any photo yet.
                        </p>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div v-else>
        <div class="card">
            <div class="card-body">
                <p class="card-text">
                    You can't access profiles without being authenticated!
                </p>
            </div>
        </div>
    </div>
    <LoadingSpinner :loading="loading"></LoadingSpinner>
</template>
<style scoped>
.card {
    margin-bottom: 20px;
}

img {
    display: block;
    max-height: 300px;
    width: auto;
    height: auto;
}

textarea {
    resize: none;
    width: 50%;
    height: 15vh;
    display: block;
}
</style>