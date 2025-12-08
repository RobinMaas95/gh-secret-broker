<script lang="ts">
  import { ModeWatcher } from "mode-watcher";
  import Login from "./lib/Login.svelte";
  import User from "./lib/User.svelte"; // This is now the Dashboard
  import Profile from "./lib/Profile.svelte";
  import RepositorySecrets from "$lib/RepositorySecrets.svelte";
  import Header from "$lib/components/Header.svelte";

  let page = $state<"login" | "dashboard" | "profile" | "secrets">("login");
  let currentRepo = $state<{ owner: string; repo: string } | null>(null);

  function handleHashChange() {
    const hash = window.location.hash;
    if (hash === "#/dashboard") {
      page = "dashboard";
    } else if (hash === "#/profile") {
      page = "profile";
    } else if (hash.startsWith("#/repo/")) {
      const parts = hash.split("/");
      // Expected format: #/repo/Owner/RepoName
      if (parts.length === 4) {
        currentRepo = { owner: parts[2], repo: parts[3] };
        page = "secrets";
      } else {
        // invalid format, go back to dashboard
        window.location.hash = "#/dashboard";
      }
    } else {
      page = "login";
    }
  }

  // Initialize on mount
  handleHashChange();
</script>

<ModeWatcher />
<svelte:window onhashchange={handleHashChange} />

<div class="min-h-screen bg-background">
  <Header />
  <main>
    {#if page === "login"}
      <Login />
    {:else if page === "dashboard"}
      <User />
    {:else if page === "profile"}
      <Profile />
    {:else if page === "secrets" && currentRepo}
      <RepositorySecrets owner={currentRepo.owner} repo={currentRepo.repo} />
    {/if}
  </main>
</div>

<style>
</style>
