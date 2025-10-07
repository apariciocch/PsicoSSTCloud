/***********************
 * PsicoSST Cloud – Front conectado a API
 ***********************/
const API_BASE = location.origin; // mismo host/puerto del server Express
const $ = (s, c=document)=>c.querySelector(s);
const $$ = (s, c=document)=>[...c.querySelectorAll(s)];
const fmt = x => (x==null?'':String(x));

/* -------- Cliente API -------- */
const Api = {
  async kpis() {
    const r = await fetch(`${API_BASE}/api/kpis`);
    if (!r.ok) throw new Error('KPIs');
    return r.json();
  },
  async evaluaciones_list() {
    const r = await fetch(`${API_BASE}/api/evaluaciones`);
    if (!r.ok) throw new Error('Evaluaciones');
    return r.json();
  },
  async evaluaciones_create({ tenant_id, instrumento, area, periodo, total_trab }) {
    const r = await fetch(`${API_BASE}/api/evaluaciones`, {
      method:'POST',
      headers:{'Content-Type':'application/json'},
      body: JSON.stringify({ tenant_id, instrumento, area, periodo, total_trab })
    });
    if (!r.ok) throw new Error('Crear evaluación');
    return r.json();
  },
  async evaluaciones_finalizar(id) {
    const r = await fetch(`${API_BASE}/api/evaluaciones/${id}/finalizar`, { method:'PUT' });
    if (!r.ok) throw new Error('Finalizar evaluación');
    return r.json();
  },
  async evaluaciones_delete(id) {
    const r = await fetch(`${API_BASE}/api/evaluaciones/${id}`, { method:'DELETE' });
    if (!r.ok) throw new Error('Eliminar evaluación');
    return r.json();
  },
  async acciones_list() {
    const r = await fetch(`${API_BASE}/api/acciones`);
    if (!r.ok) throw new Error('Acciones');
    return r.json();
  },
  async acciones_create(payload) {
    const r = await fetch(`${API_BASE}/api/acciones`, {
      method:'POST', headers:{'Content-Type':'application/json'},
      body: JSON.stringify(payload)
    });
    if (!r.ok) throw new Error('Crear acción');
    return r.json();
  },
  async acciones_update(id, payload) {
    const r = await fetch(`${API_BASE}/api/acciones/${id}`, {
      method:'PUT', headers:{'Content-Type':'application/json'},
      body: JSON.stringify(payload)
    });
    if (!r.ok) throw new Error('Actualizar acción');
    return r.json();
  }
};

/* -------- Router por hash -------- */
const routes = {
  '': 'views/dashboard.html',
  '#/dashboard': 'views/dashboard.html',
  '#/evaluaciones': 'views/evaluaciones.html',
  '#/reportes': 'views/reportes.html',
  '#/carpeta': 'views/carpeta.html',
  '#/plan': 'views/plan.html',
  '#/capacitaciones': 'views/capacitaciones.html',
  '#/casos': 'views/casos.html',
  '#/usuarios': 'views/usuarios.html',
  '#/config': 'views/config.html',
};
const viewContainer = document.getElementById('view-container');
const menu = document.getElementById('menu');
const crumb = document.getElementById('crumb');

/* -------- Carga de vistas -------- */
async function loadView() {
  const hash = window.location.hash || '#/dashboard';
  const path = routes[hash] || routes['#/dashboard'];

  $$('.menu .item').forEach(a => a.classList.toggle('active', a.getAttribute('href') === hash));
  const current = (menu.querySelector(`a[href="${hash}"]`)?.textContent || 'Dashboard').replace(/\s*\d*\s*pendientes/,'').trim();
  crumb.textContent = current;

  try {
    const res = await fetch(path, { cache:'no-cache' });
    const html = await res.text();
    viewContainer.innerHTML = html;
    await initPerView(hash);
  } catch (err) {
    console.error(err);
    viewContainer.innerHTML = `<div class="card"><h3>Error</h3><p class="muted">No se pudo cargar la vista: ${path}</p></div>`;
  }
}

/* -------- Controladores por vista -------- */
async function initPerView(hash) {
  // DASHBOARD
  if (hash === '#/dashboard' || hash === '') {
    try {
      const kpis = await Api.kpis();
      // si tienes multi-tenant, elige el primero por ahora:
      const k = kpis[0] || { total_respondido:0, total_programado:0, tasa_participacion_pct:0, casos_abiertos:0, riesgo_global_aprox:'Bajo' };
      $('#kpi-respondidas') && ($('#kpi-respondidas').textContent = fmt(k.total_respondido));
      $('#kpi-tasa') && ($('#kpi-tasa').textContent = `${k.tasa_participacion_pct}% tasa`);
      $('#kpi-riesgo') && ($('#kpi-riesgo').textContent = k.riesgo_global_aprox);
      $('#kpi-riesgo-pill') && ($('#kpi-riesgo-pill').className = 'pill ' + (k.riesgo_global_aprox === 'Bajo' ? 'ok' : 'warn'));
      $('#kpi-casos') && ($('#kpi-casos').textContent = fmt(k.casos_abiertos));
      $('#kpi-casos-pill') && ($('#kpi-casos-pill').textContent = k.casos_abiertos >= 3 ? '3 urgentes' : 'OK');
      $('#kpi-casos-pill') && ($('#kpi-casos-pill').className = 'pill ' + (k.casos_abiertos >= 3 ? 'danger' : 'ok'));
    } catch (e) {
      console.error(e);
    }
  }

  // EVALUACIONES
  if (hash === '#/evaluaciones') {
    const tbody = $('#tbl-evaluaciones-body');
    try {
      const list = await Api.evaluaciones_list();
      tbody.innerHTML = list.map(ev => `
        <tr data-id="${ev.id}">
          <td>${ev.instrumento}</td>
          <td>${ev.area}</td>
          <td>${formatPeriodoFromDate(ev.periodo)}</td>
          <td><span class="pill ${ev.estado === 'Completado' ? 'ok' : 'warn'}">${ev.estado}</span></td>
          <td>${ev.respondieron}/${ev.total_trab}</td>
          <td>
            <button class="btn btn-mini" data-act="ver">Ver</button>
            <button class="btn btn-mini" data-act="pdf">PDF</button>
            <button class="btn btn-mini" data-act="fin">Finalizar</button>
            <button class="btn btn-mini" data-act="del">🗑</button>
          </td>
        </tr>
      `).join('');

      // Acciones por fila
      tbody.addEventListener('click', async (e)=>{
        const btn = e.target.closest('button[data-act]');
        if (!btn) return;
        const tr = e.target.closest('tr');
        const id = tr?.dataset.id;
        const act = btn.dataset.act;

        if (act === 'del') {
          if (confirm('¿Eliminar evaluación?')) {
            await Api.evaluaciones_delete(id);
            await loadView();
          }
        }
        if (act === 'fin') {
          await Api.evaluaciones_finalizar(id);
          await loadView();
        }
        if (act === 'ver') alert('Ver resultados (implementa tu vista)');
        if (act === 'pdf') alert('Generar/descargar PDF (implementa tu endpoint)');
      });

      // Crear desde el form de la toolbar
      const form = $('#form-eval');
      form?.addEventListener('submit', async (e)=>{
        e.preventDefault();
        const f = e.target;
        const instrumento = f.instrumento.value;
        const area = f.area.value.trim();
        const periodo = f.periodo.value;              // yyyy-mm
        const total_trab = parseInt(f.total.value||'0',10);

        if (!instrumento || !area || !periodo) return alert('Completa los campos');
        // Por ahora asumimos 1 tenant (el primero de KPIs)
        const kpis = await Api.kpis();
        const tenant_id = (kpis[0] && kpis[0].tenant_id) || null;
        if (!tenant_id) return alert('No hay tenant en BD');

        await Api.evaluaciones_create({ tenant_id, instrumento, area, periodo, total_trab });
        f.reset();
        await loadView();
      });
    } catch (e) {
      console.error(e);
      tbody.innerHTML = `<tr><td colspan="6">Error cargando evaluaciones</td></tr>`;
    }
  }

  // REPORTES (tabs demo sin backend aún)
  if (hash === '#/reportes') {
    const tabs = $('#tabs-reportes');
    const panels = {
      global: $('#panel-global'),
      individuales: $('#panel-individuales'),
      comparativos: $('#panel-comparativos'),
    };
    tabs?.addEventListener('click', (e)=>{
      const btn = e.target.closest('.tab');
      if (!btn) return;
      $$('.tab', tabs).forEach(t=>t.classList.remove('active'));
      btn.classList.add('active');
      const tab = btn.dataset.tab;
      Object.keys(panels).forEach(k => panels[k].style.display = (k===tab?'block':'none'));
    });
  }
}

/* -------- Helpers -------- */
function formatPeriodoFromDate(dateStr) {
  // "2025-10-01" → "Oct-2025"
  if (!dateStr) return '';
  const d = new Date(dateStr);
  const meses = ['Ene','Feb','Mar','Abr','May','Jun','Jul','Ago','Sep','Oct','Nov','Dic'];
  return `${meses[d.getUTCMonth()]}-${d.getUTCFullYear()}`;
}

/* -------- Modal genérico (igual) -------- */
const modalBg = document.getElementById('modal-bg');
const openModal = (title='Nueva evaluación')=>{
  document.getElementById('modal-title').textContent = title;
  modalBg.style.display='grid';
  modalBg.setAttribute('aria-hidden','false');
};
const closeModal = ()=>{
  modalBg.style.display='none';
  modalBg.setAttribute('aria-hidden','true');
};
document.getElementById('modal-close').onclick = closeModal;
document.getElementById('modal-cancel').onclick = closeModal;

/* -------- Topbar -------- */
document.getElementById('btnNew').onclick = ()=> openModal('Nueva evaluación');
document.getElementById('btnZip').onclick = ()=> alert('Generar Carpeta SUNAFIL (ZIP) — implementa /api/zip cuando quieras');
document.getElementById('q').addEventListener('keydown', (e)=>{
  if (e.key === 'Enter') alert('Buscar: ' + e.target.value + ' (pendiente implementar)');
});

/* -------- Router -------- */
window.addEventListener('hashchange', loadView);
window.addEventListener('DOMContentLoaded', loadView);
