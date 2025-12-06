<script lang="ts">
  import Login from "./lib/Login.svelte";
  import User from "./lib/User.svelte";

  let page = $state<"login" | "user">("login");

  function handleHashChange() {
    const hash = window.location.hash;
    if (hash === "#/userpage") {
      page = "user";
    } else {
      page = "login";
    }
  }

  // Initialize on mount
  handleHashChange();
</script>

<svelte:window onhashchange={handleHashChange} />

<main>
  {#if page === "login"}
    <Login />
  {:else if page === "user"}
    <User />
  {/if}
</main>

<style>
</style>
