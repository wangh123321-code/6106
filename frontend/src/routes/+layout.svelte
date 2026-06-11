<script>
  import { onMount } from 'svelte';
  import { isLoggedIn, getCurrentUser, logout } from '$lib/api';

  let user = null;
  let loggedIn = false;

  onMount(() => {
    loggedIn = isLoggedIn();
    user = getCurrentUser();
  });

  function handleLogout() {
    logout();
  }
</script>

<svelte:head>
  <title>乒乓球赛事管理平台</title>
</svelte:head>

<div class="app">
  <nav>
    <a href="/" class="nav-brand">🏓 赛事管理</a>
    {#if loggedIn}
      <a href="/tournaments">赛事列表</a>
      <a href="/appeals">成绩申诉</a>
      {#if user?.role === 'committee'}
        <a href="/admin">管理后台</a>
        <a href="/audit">审计日志</a>
      {/if}
      {#if user?.role === 'referee'}
        <a href="/referee">裁判工作台</a>
      {/if}
      {#if user?.role === 'player'}
        <a href="/my-registrations">我的报名</a>
      {/if}
      <div class="nav-right">
        <span>{user?.display_name || user?.username}</span>
        <button class="btn btn-outline" on:click={handleLogout}>退出</button>
      </div>
    {:else}
      <div class="nav-right">
        <a href="/login">登录</a>
        <a href="/register">注册</a>
      </div>
    {/if}
  </nav>

  <main class="container">
    <slot />
  </main>
</div>

<style>
  .app { min-height: 100vh; }
  main { padding-top: 24px; padding-bottom: 48px; }
</style>
