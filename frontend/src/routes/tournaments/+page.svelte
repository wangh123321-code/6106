<script>
  import { api, isLoggedIn } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  let tournaments = [];

  onMount(async () => {
    if (!isLoggedIn()) { goto('/login'); return; }
    const res = await api.get('/tournaments');
    if (res.data) tournaments = res.data;
  });

  function statusLabel(s) {
    const map = { draft: '草稿', open: '报名中', registration_closed: '报名截止', drawn: '已抽签', in_progress: '进行中', completed: '已结束', published: '已发布' };
    return map[s] || s;
  }

  function statusClass(s) {
    const map = { open: 'badge-success', registration_closed: 'badge-warning', drawn: 'badge-info', in_progress: 'badge-info', completed: 'badge-success', published: 'badge-success' };
    return map[s] || 'badge-info';
  }
</script>

<svelte:head>
  <title>赛事列表 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>赛事列表</h1>
</div>

<div class="card">
  <table>
    <thead>
      <tr>
        <th>赛事名称</th>
        <th>地点</th>
        <th>开始日期</th>
        <th>结束日期</th>
        <th>状态</th>
        <th>操作</th>
      </tr>
    </thead>
    <tbody>
      {#each tournaments as t}
        <tr>
          <td>{t.name}</td>
          <td>{t.location}</td>
          <td>{t.start_date?.substring(0,10)}</td>
          <td>{t.end_date?.substring(0,10)}</td>
          <td><span class="badge {statusClass(t.status)}">{statusLabel(t.status)}</span></td>
          <td><a href="/tournaments/{t.id}" class="btn btn-outline" style="padding:4px 12px;">查看</a></td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>
