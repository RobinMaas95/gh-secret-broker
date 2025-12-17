<script lang="ts">
    import User from "$lib/User.svelte";

    let { data } = $props();
</script>

<User {data} />
