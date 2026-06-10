<script>
  import { api, isLoggedIn, getCurrentUser } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  let logs = [];
  let filterType = 'match';
  let filterId = '';

  onMount(async () => {
    const user = getCurrentUser();
    if (user?.role !== 'committee') { goto('/'); return; }
    await loadLogs();
  });

  async function loadLogs() {
    let path = '/audit-logs?';
    if (filterType) path += `target_type=${filterType}&`;
    if (filterId) path += `target_id=${filterId}&`;
    const res = await api.get(path);
    if (res.data) logs = res.data;
  }
</script>

<svelte:head>
  <title>审计日志 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>审计日志</h1>
  <button class="btn btn-outline" on:click={loadLogs}>刷新</button>
</div>

<div class="card" style="margin-bottom:16px;">
  <div style="display:flex;gap:12px;align-items:center;">
    <div class="form-group" style="margin:0;">
      <label>操作类型</label>
      <select bind:value={filterType} onchange={loadLogs}>
        <option value="match">比分修改</option>
        <option value="">全部</option>
      </select>
    </div>
    <div class="form-group" style="margin:0;">
      <label>目标ID</label>
      <input type="text" bind:value={filterId} placeholder="输入比赛ID" />
    </div>
    <button class="btn btn-primary" style="margin-top:20px;" on:click={loadLogs}>查询</button>
  </div>
</div>

<div class="card">
  <table>
    <thead>
      <tr>
        <th>时间</th>
        <th>操作人</th>
        <th>操作</th>
        <th>目标类型</th>
        <th>目标ID</th>
      </tr>
    </thead>
    <tbody>
      {#each logs as log}
        <tr>
          <td>{log.created_at?.substring(0,19).replace('T', ' ')}</td>
          <td>{log.operator_name}</td>
          <td>{log.action}</td>
          <td>{log.target_type}</td>
          <td style="font-size:12px;font-family:monospace;">{log.target_id}</td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
