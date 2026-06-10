<script>
  import { api } from '$lib/api';
  import { goto } from '$app/navigation';

  let username = '';
  let password = '';
  let displayName = '';
  let role = 'player';
  let error = '';

  async function handleRegister() {
    error = '';
    if (password.length < 6) {
      error = '密码至少6位';
      return;
    }
    const res = await api.post('/register', { username, password, display_name: displayName, role });
    if (res.code === 0) {
      goto('/login');
    } else {
      error = res.message || '注册失败';
    }
  }
</script>

<svelte:head>
  <title>注册 - 赛事管理平台</title>
</svelte:head>

<div class="login-page">
  <div class="login-card card">
    <h2>注册</h2>
    {#if error}
      <div class="alert alert-danger">{error}</div>
    {/if}
    <form on:submit|preventDefault={handleRegister}>
      <div class="form-group">
        <label>用户名</label>
        <input type="text" bind:value={username} required />
      </div>
      <div class="form-group">
        <label>姓名</label>
        <input type="text" bind:value={displayName} required />
      </div>
      <div class="form-group">
        <label>密码</label>
        <input type="password" bind:value={password} required />
      </div>
      <div class="form-group">
        <label>角色</label>
        <select bind:value={role}>
          <option value="player">选手</option>
          <option value="referee">裁判</option>
          <option value="committee">组委会</option>
        </select>
      </div>
      <button type="submit" class="btn btn-primary" style="width:100%;justify-content:center;margin-top:8px;">注册</button>
    </form>
    <p style="text-align:center;margin-top:16px;font-size:14px;">
      已有账号？<a href="/login">去登录</a>
    </p>
  </div>
</div>

<style>
  .login-page { display: flex; justify-content: center; padding-top: 60px; }
  .login-card { width: 400px; }
  .login-card h2 { text-align: center; margin-bottom: 24px; }
  .alert { padding: 10px 14px; border-radius: var(--radius); margin-bottom: 16px; font-size: 14px; }
  .alert-danger { background: #fce8e6; color: #c5221f; }
</style>
