layui.use(['form', 'layer', 'table'], function () {
  const form = layui.form;
  const layer = layui.layer;
  const table = layui.table;

  let configCache = null;

  async function request(url, options) {
    const response = await fetch(url, options || {});
    const data = await response.json();
    if (!response.ok) {
      throw new Error(data.error || 'request failed');
    }
    return data;
  }

  function splitLines(value) {
    return value
      .split('\n')
      .map(item => item.trim())
      .filter(Boolean);
  }

  function applyConfig(config) {
    configCache = config;
    const formElement = document.getElementById('config-form');
    formElement.adb_path.value = config.adb_path || '';
    formElement.device_serial.value = config.device_serial || '';
    formElement.tesseract_path.value = config.tesseract_path || '';
    formElement.safe_mode.checked = !!config.safe_mode;
    formElement.auto_claim_reward.checked = !!config.auto_claim_reward;
    formElement.switch_room.checked = !!config.switch_room;
    formElement.debug_dry_run.checked = !!(config.debug && config.debug.dry_run);
    formElement.strategy_min_join_seconds.value = config.strategy ? config.strategy.min_join_seconds : 0;
    formElement.strategy_max_join_seconds.value = config.strategy ? config.strategy.max_join_seconds : 0;
    formElement.loop_interval_seconds.value = config.loop_interval_seconds || 0;
    formElement.debug_max_iterations.value = config.debug ? config.debug.max_iterations : 0;
    formElement.strategy_prize_whitelist.value = (config.strategy && config.strategy.prize_whitelist || []).join('\n');
    formElement.strategy_prize_blacklist.value = (config.strategy && config.strategy.prize_blacklist || []).join('\n');
    form.render();
  }

  async function loadConfig() {
    const config = await request('/api/config');
    applyConfig(config);
  }

  function collectConfig(fields) {
    const next = JSON.parse(JSON.stringify(configCache || {}));
    next.adb_path = fields.adb_path;
    next.device_serial = fields.device_serial;
    next.tesseract_path = fields.tesseract_path;
    next.safe_mode = !!fields.safe_mode;
    next.auto_claim_reward = !!fields.auto_claim_reward;
    next.switch_room = !!fields.switch_room;
    next.loop_interval_seconds = Number(fields.loop_interval_seconds || 0);
    next.debug = next.debug || {};
    next.debug.dry_run = !!fields.debug_dry_run;
    next.debug.max_iterations = Number(fields.debug_max_iterations || 0);
    next.strategy = next.strategy || {};
    next.strategy.min_join_seconds = Number(fields.strategy_min_join_seconds || 0);
    next.strategy.max_join_seconds = Number(fields.strategy_max_join_seconds || 0);
    next.strategy.prize_whitelist = splitLines(fields.strategy_prize_whitelist || '');
    next.strategy.prize_blacklist = splitLines(fields.strategy_prize_blacklist || '');
    return next;
  }

  async function saveConfig(fields) {
    const payload = collectConfig(fields);
    await request('/api/config', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    });
    layer.msg('配置已保存');
    await loadConfig();
  }

  async function loadTable(elem, url, cols) {
    const result = await request(url);
    table.render({
      elem,
      data: result.items || [],
      page: false,
      limit: result.items ? result.items.length : 0,
      cols: [cols]
    });
  }

  async function loadLogs() {
    const result = await request('/api/logs?limit=300');
    document.getElementById('logs-box').textContent = (result.items || []).join('\n');
  }

  function screenshotCard(item) {
    return `
      <div class="layui-col-md3 layui-col-sm6">
        <div class="shot-card">
          <a href="${item.url}" target="_blank" rel="noreferrer">
            <img src="${item.url}" alt="${item.name}">
          </a>
          <div class="shot-meta">
            <div>${item.name}</div>
            <div>${new Date(item.mod_time).toLocaleString()}</div>
            <div>${(item.size / 1024).toFixed(1)} KB</div>
          </div>
        </div>
      </div>
    `;
  }

  async function loadScreenshots() {
    const result = await request('/api/screenshots');
    const items = result.items || [];
    document.getElementById('screenshots-grid').innerHTML = items.map(screenshotCard).join('');
  }

  form.on('submit(save-config)', async function (data) {
    try {
      await saveConfig(data.field);
    } catch (error) {
      layer.alert(error.message);
    }
    return false;
  });

  document.getElementById('reload-config').addEventListener('click', function () {
    loadConfig().catch(error => layer.alert(error.message));
  });
  document.getElementById('reload-logs').addEventListener('click', function () {
    loadLogs().catch(error => layer.alert(error.message));
  });
  document.getElementById('reload-shots').addEventListener('click', function () {
    loadScreenshots().catch(error => layer.alert(error.message));
  });

  Promise.all([
    loadConfig(),
    loadTable('#records-table', '/api/records', [
      { field: 'record_id', title: 'ID', minWidth: 220 },
      { field: 'room_id', title: '直播间', minWidth: 140 }
      { field: 'prize_text', title: '奖品', minWidth: 180 },
      { field: 'result', title: '结果', width: 120 },
      { field: 'created_at', title: '时间', minWidth: 180 },
    ]),
    loadTable('#events-table', '/api/events', [
      { field: 'created_at', title: '时间', minWidth: 180 },
      { field: 'type', title: '类型', width: 160 },
      { field: 'message', title: '消息', minWidth: 260 },
      { field: 'device_id', title: '设备', minWidth: 140 }
    ]),
    loadLogs(),
    loadScreenshots()
  ]).catch(error => layer.alert(error.message));
});
