<script>
  import { api } from '$lib/api';
  import { onMount } from 'svelte';

  let registrations = [];

  onMount(async () => {
    const res = await api.get('/profile');
  });
</script>

<svelte:head>
  <title>我的报名 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>我的报名</h1>
</div>

<div class="card" style="text-align:center;padding:32px;">
  <p style="color:var(--text-secondary);">您可以通过赛事详情页进行报名</p>
  <a href="/tournaments" class="btn btn-primary" style="margin-top:12px;">浏览赛事</a>
</div>
