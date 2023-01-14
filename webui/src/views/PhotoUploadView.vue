<script>
export default {
    components: {},
    data: function () {
        return {
            errormsg: null,
            successmsg: null,
            loading: false,
            Image: {
                MainComment: "",
                UploadedPhoto: null
            },
            token: null
        }
    },
    methods: {
        onChange($event) {
            this.Image.UploadedPhoto = $event.target.files[0];
        },
        async uploadPhoto() {
            this.loading = true;
            this.errormsg = null;
            this.successmsg = null;
            try {
                if (this.Image.UploadedPhoto != null) {
                    await this.$axios.post("/photos", this.Image, {
                        headers: {
                            'Authorization': 'Bearer ' + localStorage.token,
                            'Content-Type': 'multipart/form-data'
                        }
                    });
                    this.successmsg = "Photo uploaded successfully."
                    this.Image.MainComment = null
                    this.Image.UploadedPhoto = null
                } else {
                    this.errormsg = "Select an image to upload."
                }
            } catch (e) {
                if (e.response && e.response.status == '401') {
                    this.errormsg = "You have be authenticated to upload a photo."
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
    <div v-if="!loading">
        <div
            class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
            <h1 class="h2">Upload Photo</h1>
        </div>

        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        <SuccessMsg v-if="successmsg" :msg="successmsg"></SuccessMsg>

        <div v-if="token">
            <div class="mb-3">
                <label>
                    <label for="comment" class="form-label">Insert a comment to write with the photo</label>
                    <textarea class="form-control" id="comment" v-model="Image.MainComment"></textarea>
                    <br />
                    <input type="file" id="photo" name="photo" accept="image/png, image/jpeg" v-on:change="onChange" />
                </label>
            </div>

            <button type="button" class="btn btn-sm btn-primary" @click="uploadPhoto">
                Upload Photo
            </button>
        </div>
        <div v-else>
            <div class="card">
                <div class="card-body">
                    <p class="card-text">
                        You can't upload photos without being authenticated!
                    </p>
                </div>
            </div>
        </div>
    </div>
    <LoadingSpinner :loading="loading"></LoadingSpinner>
</template>

<style>
textarea {
    resize: none;
    width: 100%;
    height: 20vh;
    display: block;
}
</style>
