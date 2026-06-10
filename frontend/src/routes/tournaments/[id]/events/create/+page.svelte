<script>
  import { api, getCurrentUser } from '$lib/api';
  import { page } from '$app/stores';
  import { onMount } from 'svelte';

  let tournament = null;
  let name = '';
  let type = 'ms';
  let drawMethod = 'snake';
  let seedCount = 4;
  let bracketSize = 64;
  let message = '';

  $: tournamentId = $page.params.id;

  onMount(async () => {
    const res = await api.get(`/tournaments/${tournamentId}`);
    if (res.data) tournament = res.data;
  });

  async function createEvent() {
    const res = await api.post(`/tournaments/${tournamentId}/events`, {
      name,
      type,
      draw_method: drawMethod,
      seed_count: seedCount,
      bracket_size: bracketSize
    });
    message = res.code === 0 ? '项目创建成功' : (res.message || '创建失败');
  }
</script>

<svelte:head>
  <title>创建项目 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>创建比赛项目</h1>
</div>

<div class="card" style="max-width:500px;">
  {#if message}
    <div class="alert">{message}</div>
  {/if}

  <div class="form-group">
    <label>项目名称</label>
    <input type="text" bind:value={name} placeholder="如：男子单打" />
  </div>

  <div class="form-group">
    <label>项目类型</label>
    <select bind:value={type}>
      <option value="ms">男子单打</option>
      <option value="ws">女子单打</option>
      <option value="xd">混合双打</option>
    </select>
  </div>

  <div class="form-group">
    <label>抽签方式</label>
    <select bind:value={drawMethod}>
      <option value="snake">积分蛇形分组</option>
      <option value="random">随机抽签</option>
    </select>
  </div>

  <div class="form-group">
    <label>种子选手数</label>
    <input type="number" bind:value={seedCount} min="0" />
  </div>

  <div class="form-group">
    <label>签位大小</label>
    <select bind:value={bracketSize}>
      <option value="8">8</option>
      <option value="16">16</option>
      <option value="32">32</option>
      <option value="64">64</option>
      <option value="128">128</option>
    </select>
  </div>

  <button class="btn btn-primary" style="width:100%;justify-content:center;" on:click={createEvent}>
    创建项目
  </button>
</div>

<style>
  .alert { padding: 10px 14px; border-radius: var(--radius); margin-bottom: 16px; font-size: 14px; background: #e8f0fe; color: #1a73e8; }
</style>
