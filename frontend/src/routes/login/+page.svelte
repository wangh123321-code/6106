<script>
  import { api } from '$lib/api';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  let username = '';
  let password = '';
  let error = '';

  async function handleLogin() {
    error = '';
    const res = await api.post('/login', { username, password });
    if (res.data?.token) {
      localStorage.setItem('token', res.data.token);
      localStorage.setItem('user', JSON.stringify(res.data.user));
      goto('/');
    } else {
      error = res.message || '登录失败';
    }
  }
</script>

<svelte:head>
  <title>登录 - 赛事管理平台</title>
</svelte:head>

<div class="login-page">
  <div class="login-card card">
    <h2>登录</h2>
    {#if error}
      <div class="alert alert-danger">{error}</div>
    {/if}
    <form on:submit|preventDefault={handleLogin}>
      <div class="form-group">
        <label for="username">用户名</label>
        <input id="username" type="text" bind:value={username} placeholder="请输入用户名" required />
      </div>
      <div class="form-group">
        <label for="password">密码</label>
        <input id="password" type="password" bind:value={password} placeholder="请输入密码" required />
      </div>
      <button type="submit" class="btn btn-primary" style="width:100%;justify-content:center;margin-top:8px;">登录</button>
    </form>
    <p style="text-align:center;margin-top:16px;font-size:14px;">
      还没有账号？<a href="/register">立即注册</a>
    </p>
  </div>
</div>

<style>
  .login-page { display: flex; justify-content: center; padding-top: 80px; }
  .login-card { width: 400px; }
  .login-card h2 { text-align: center; margin-bottom: 24px; }
  .alert { padding: 10px 14px; border-radius: var(--radius); margin-bottom: 16px; font-size: 14px; }
  .alert-danger { background: #fce8e6; color: #c5221f; }
</style>
