<script>
export default {
    props: ['msg', 'details']
}
</script>

<template>
    <div class="alert alert-success" role="alert">
        {{ msg }}
        <small v-if="details"><br /><br />Details: {{ details }}</small>
    </div>
</template>

<style>

</style>