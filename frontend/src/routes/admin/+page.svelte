<script>
  import { api, isLoggedIn, getCurrentUser } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  let user = null;
  let tournaments = [];

  onMount(async () => {
    if (!isLoggedIn()) { goto('/login'); return; }
    user = getCurrentUser();
    if (user?.role !== 'committee') { goto('/'); return; }
  });
</script>

<svelte:head>
  <title>管理后台 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>管理后台</h1>
</div>

<div class="admin-grid">
  <a href="/tournaments" class="card admin-card">
    <div class="icon">🏆</div>
    <h3>赛事管理</h3>
    <p>创建赛事、设置项目、管理状态</p>
  </a>

  <a href="/admin/users" class="card admin-card">
    <div class="icon">👥</div>
    <h3>用户管理</h3>
    <p>管理选手、裁判、设置排名积分</p>
  </a>

  <a href="/admin/referee-assign" class="card admin-card">
    <div class="icon">📋</div>
    <h3>裁判指派</h3>
    <p>为各场次指派裁判</p>
  </a>

  <a href="/audit" class="card admin-card">
    <div class="icon">📝</div>
    <h3>审计日志</h3>
    <p>查看所有操作修改记录</p>
  </a>

  <a href="/admin/medal-board" class="card admin-card">
    <div class="icon">🥇</div>
    <h3>奖牌榜</h3>
    <p>查看与发布总成绩排行</p>
  </a>

  <a href="/admin/certificates" class="card admin-card">
    <div class="icon">📜</div>
    <h3>证书导出</h3>
    <p>生成并导出PDF获奖证书</p>
  </a>
</div>

<style>
  .admin-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 16px;
  }

  .admin-card {
    text-decoration: none;
    color: inherit;
    text-align: center;
    padding: 32px 20px;
    transition: transform 0.2s, box-shadow 0.2s;
  }

  .admin-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  }

  .admin-card .icon { font-size: 36px; margin-bottom: 12px; }
  .admin-card h3 { font-size: 16px; margin-bottom: 8px; }
  .admin-card p { font-size: 13px; color: var(--text-secondary); }
</style>
