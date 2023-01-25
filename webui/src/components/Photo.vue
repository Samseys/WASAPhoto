<script>
export default {
    props: {
        photoID: Number
    },
    data: function data() {
        return {
            photo: {
                PhotoID: Number,
                Owner: {
                    UserID: Number,
                    Username: String
                },
                CreationDate: String,
                Comment: String,
                Comments: [{
                    CommentID: Number,
                    Comment: String,
                    Owner: {
                        UserID: Number,
                        Username: String
                    },
                    CreationDate: String
                }],
                Likes: [{
                    UserID: Number,
                    Username: String
                }]
            },
            errormsg: null,
            comment: "",
            token: null,
            photoFile: null,
            loading: false
        };
    },
    methods: {
        async getImage() {
            this.loading = true;
            let response = await this.$axios.get("/photos/" + this.photoID + "");
            this.photo = response.data;

            response = await this.$axios.get("/photos/" + this.photoID + "/file", { responseType: 'blob' });
            this.photoFile = window.URL.createObjectURL(response.data);
            this.loading = false;
        },
        async deletePhoto() {
            try {
                await this.$axios.delete("/photos/" + this.photoID);
                this.$emit("delete-photo");
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
        async like() {
            this.errormsg = ""
            try {
                await this.$axios.put("/photos/" + this.photoID + "/likes/" + this.token, "");
                if (this.photo.Likes == null) {
                    this.photo.Likes = []
                }
                this.photo.Likes.push({ "UserID": this.token, "Username": localStorage.username })
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
        async unlike() {
            this.errormsg = ""
            try {
                await this.$axios.delete("/photos/" + this.photoID + "/likes/" + this.token);
                this.photo.Likes.splice(this.photo.Likes.findIndex(item => item.UserID == this.token), 1)
            } catch (e) {
                if (e.response && e.response.status == '403') {
                    this.errormsg = "The owner of this profile banned you.";
                } if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no photo with this id: " + this.photoID + ".";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
        },
        async postComment() {
            this.errormsg = ""
            if (this.comment == null || this.comment.replace(/\s+/g, "") == "") {
                this.errormsg = "You can't post an empty comment";
                return;
            }
            try {
                let response = await this.$axios.post("/photos/" + this.photoID + "/comments", { Comment: this.comment });
                if (this.photo.Comments == null) {
                    this.photo.Comments = []
                }
                this.photo.Comments.push(response.data);
                this.comment = ""
            } catch (e) {
                if (e.response && e.response.status == '403') {
                    this.errormsg = "The owner of this profile banned you.";
                } if (e.response && e.response.status == '404') {
                    this.errormsg = "There is no photo with this id: " + this.photoID + ".";
                } else if (e.response && e.response.status == '500') {
                    this.errormsg = "An internal error has occured.";
                } else {
                    this.errormsg = e.toString();
                }
            }
        },
        async deleteComment(c) {
            try {
                await this.$axios.delete("/photos/" + this.photoID + "/comments/" + c.CommentID);
                this.photo.Comments.splice(this.photo.Comments.findIndex(item => item.CommentID == c.CommentID), 1)
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
        this.token = localStorage.token;
        this.getImage()
    }
}

</script>
<template>
    <div v-if="!this.loading">
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        <div class="card">
            <div class="card-header d-flex justify-content-between align-items-center">
                <span>
                    <RouterLink :to="'/profile/' + this.photo.Owner.UserID">
                        {{ this.photo.Owner.Username }}
                    </RouterLink>
                    - {{ this.photo.CreationDate }}
                </span>
                <button type="button" class="btn btn-danger" v-if="this.photo.Owner.UserID == this.token"
                    @click="deletePhoto()">
                    Delete photo
                </button>
            </div>
            <div class="card-body">
                <p class="card-text" style="white-space:pre-line" v-if="this.photo.Comment">
                    {{ this.photo.Comment }}
                    <br />
                </p>
                <img :src="this.photoFile" />
                <br />
                <p class="card-text">
                    Likes: {{ (this.photo.Likes ?? []).length }}
                </p>

                <div class="btn-toolbar">
                    <button type="button" class="btn btn-danger" @click="unlike()"
                        v-if="(this.photo.Likes ?? []).some(like => like.UserID == this.token)">
                        Unlike
                    </button>

                    <button type="button" class="btn btn-primary" @click="like()" v-else>
                        Like
                    </button>
                </div>
                <br />
                <div class="card">
                    <div class="card-header">Comments: {{ (this.photo.Comments ?? []).length}}</div>
                    <div class="card-body">
                        <div class="card" v-for="c in this.photo.Comments" :key="c.CommentID">
                            <div class="card-header d-flex justify-content-between align-items-center">
                                <span>
                                    <RouterLink :to="'/profile/' + c.Owner.UserID">
                                        {{ c.Owner.Username }}
                                    </RouterLink>
                                    - {{ c.CreationDate }}
                                </span>
                                <button type="button" class="btn btn-danger small" @click="deleteComment(c)"
                                    v-if="c.Owner.UserID == this.token">
                                    Delete comment
                                </button>
                            </div>
                            <div class="card-body">
                                <p class="card-text" style="white-space:pre-line">
                                    {{ c.Comment }}
                                </p>
                            </div>
                        </div>
                        <div class="card">
                            <div class="card-header">Post a comment</div>
                            <div class="card-body d-flex justify-content-between align-items-center">
                                <textarea class="form-control" v-model="this.comment"></textarea>

                                <button type="button" class="btn btn-primary" @click="postComment()">
                                    Comment
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
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