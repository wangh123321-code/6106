<script>
  import { api } from '$lib/api';
  import { onMount } from 'svelte';

  let users = [];
  let roleFilter = '';

  onMount(async () => {
    await loadUsers();
  });

  async function loadUsers() {
    const res = await api.get(`/users${roleFilter ? '?role=' + roleFilter : ''}`);
    if (res.data) users = res.data;
  }
</script>

<svelte:head>
  <title>用户管理 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>用户管理</h1>
</div>

<div class="card" style="margin-bottom:16px;">
  <div style="display:flex;gap:8px;align-items:center;">
    <label style="margin:0;">角色筛选:</label>
    <select bind:value={roleFilter} onchange={loadUsers}>
      <option value="">全部</option>
      <option value="player">选手</option>
      <option value="referee">裁判</option>
      <option value="committee">组委会</option>
    </select>
  </div>
</div>

<div class="card">
  <table>
    <thead>
      <tr>
        <th>用户名</th>
        <th>姓名</th>
        <th>角色</th>
        <th>排名积分</th>
      </tr>
    </thead>
    <tbody>
      {#each users as u}
        <tr>
          <td>{u.username}</td>
          <td>{u.display_name}</td>
          <td>{u.role}</td>
          <td>{u.ranking || '-'}</td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
