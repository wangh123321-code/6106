<script>
  let playerName = '';
  let eventName = '';
  let tournamentName = '';
  let medal = 'gold';
  let date = '';

  function exportPDF() {
    const params = new URLSearchParams({
      player_name: playerName,
      event_name: eventName,
      tournament_name: tournamentName,
      medal,
      date: date || new Date().toISOString().substring(0, 10)
    });
    window.open(`/api/v1/certificate/export?${params.toString()}`, '_blank');
  }
</script>

<svelte:head>
  <title>证书导出 - 赛事管理平台</title>
</svelte:head>

<div class="page-header">
  <h1>证书导出</h1>
</div>

<div class="card" style="max-width:500px;">
  <div class="form-group">
    <label>选手姓名</label>
    <input type="text" bind:value={playerName} placeholder="输入选手姓名" />
  </div>
  <div class="form-group">
    <label>赛事名称</label>
    <input type="text" bind:value={tournamentName} placeholder="输入赛事名称" />
  </div>
  <div class="form-group">
    <label>项目名称</label>
    <input type="text" bind:value={eventName} placeholder="如：男子单打" />
  </div>
  <div class="form-group">
    <label>奖项</label>
    <select bind:value={medal}>
      <option value="gold">金牌</option>
      <option value="silver">银牌</option>
      <option value="bronze">铜牌</option>
    </select>
  </div>
  <div class="form-group">
    <label>日期</label>
    <input type="date" bind:value={date} />
  </div>
  <button class="btn btn-primary" style="width:100%;justify-content:center;" on:click={exportPDF}>
    导出PDF证书
  </button>
</div>
