<script lang="ts">
    import RepositorySecrets from "$lib/RepositorySecrets.svelte";
    import type { PageData } from "./$types";

    let { data } = $props<{ data: PageData }>();
    let owner = $derived(data.owner);
    let repo = $derived(data.repo);
</script>

<RepositorySecrets {owner} {repo} secrets={data.secrets} csrfToken={data.csrfToken} />
