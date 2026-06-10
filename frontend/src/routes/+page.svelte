<script>
  import { isLoggedIn, api } from '$lib/api';
  import { onMount } from 'svelte';

  let tournaments = [];

  onMount(async () => {
    if (isLoggedIn()) {
      const res = await api.get('/tournaments');
      if (res.data) tournaments = res.data;
    }
  });
</script>

<svelte:head>
  <title>乒乓球赛事管理平台 - 首页</title>
</svelte:head>

<div class="hero">
  <h1>市级乒乓球锦标赛管理系统</h1>
  <p>从报名、抽签到成绩发布，一站式赛事管理平台</p>
  {#if !isLoggedIn()}
    <div style="margin-top: 24px; display: flex; gap: 12px; justify-content: center;">
      <a href="/login" class="btn btn-primary">登录</a>
      <a href="/register" class="btn btn-outline">注册</a>
    </div>
  {/if}
</div>

{#if tournaments.length > 0}
<div class="section">
  <h2>当前赛事</h2>
  <div class="tournament-grid">
    {#each tournaments as t}
      <a href="/tournaments/{t.id}" class="card tournament-card">
        <h3>{t.name}</h3>
        <p class="text-secondary">{t.location}</p>
        <p class="text-secondary">{t.start_date?.substring(0,10)} ~ {t.end_date?.substring(0,10)}</p>
        <span class="badge badge-info">{t.status}</span>
      </a>
    {/each}
  </div>
</div>
{/if}

<style>
  .hero {
    text-align: center;
    padding: 80px 20px;
    background: linear-gradient(135deg, #1a73e8 0%, #34a853 100%);
    color: #fff;
    border-radius: 12px;
    margin-bottom: 40px;
  }
  .hero h1 { font-size: 36px; margin-bottom: 12px; }
  .hero p { font-size: 18px; opacity: 0.9; }

  .section h2 { font-size: 20px; margin-bottom: 16px; }

  .tournament-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 16px;
  }

  .tournament-card {
    text-decoration: none;
    color: inherit;
    transition: transform 0.2s, box-shadow 0.2s;
  }
  .tournament-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  }
  .tournament-card h3 { font-size: 18px; margin-bottom: 8px; }

  .text-secondary { color: var(--text-secondary); font-size: 14px; margin-bottom: 4px; }
</style>
