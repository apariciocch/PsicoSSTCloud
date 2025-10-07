/***********************
 * PsicoSST Cloud – Core
 * Estado demo en LocalStorage (sin backend)
 ***********************/

/* ---------------------
   Utilidades básicas
----------------------*/
const $ = (sel, ctx = document) => ctx.querySelector(sel);
const $$ = (sel, ctx = document) => [...ctx.querySelectorAll(sel)];
const fmt = (x) => (x == null ? "" : String(x));

/* ---------------------
   Store (LocalStorage)
----------------------*/
const STORE_KEY = "psicosst_demo_v1";

const Seed = {
  empresa: { nombre: "AlbatrossCloud SAC", ruc: "2060XXXXXXX", instrumento: "MINSA 2024", firmaPsicologo: "Lic. Armando Aparicio – C.Ps.P XXXXX – Habilitado" },
  evaluaciones: [
    { id: crypto.randomUUID(), instrumento: "MINSA 2024", area: "Atención al Cliente", periodo: "2025-10", estado: "En curso", respuestas: 36, total: 50 },
    { id: crypto.randomUUID(), instrumento: "COPSOQ",      area: "Ventas",              periodo: "2025-09", estado: "Completado", respuestas: 42, total: 45 },
    { id: crypto.randomUUID(), instrumento: "MINSA 2024",  area: "Operaciones",         periodo: "2025-08", estado: "Completado", respuestas: 50, total: 50 },
  ],
  casos: [
    { id: "CASE-0012", tipo: "Acoso",      prioridad: "Alta",   updatedAt: "2025-10-02" },
    { id: "CASE-0013", tipo: "Conflicto",  prioridad: "Media",  updatedAt: "2025-10-03" },
  ],
  acciones: [
    { id: crypto.randomUUID(), riesgo: "Carga mental alta", accion: "Redistribuir turnos y pausas", responsable: "Jefe Operaciones", vence: "2025-10-20", estado: "En curso" },
    { id: crypto.randomUUID(), riesgo: "Liderazgo deficiente", accion: "Taller de feedback efectivo", responsable: "RRHH", vence: "2025-11-05", estado: "Pendiente" },
  ]
};

const Store = {
  _data: null,
  load() {
    const raw = localStorage.getItem(STORE_KEY);
    if (raw) {
      try { this._data = JSON.parse(raw); }
      catch { this._data = structuredClone(Seed); }
    } else {
      this._data = structuredClone(Seed);
      this.save();
    }
    return this._data;
  },
  save() { localStorage.setItem(STORE_KEY, JSON.stringify(this._data)); },
  get() { return this._data || this.load(); },

  // CRUD evaluaciones
  eval_list() { return this.get().evaluaciones; },
  eval_add(e) { this.get().evaluaciones.unshift({ id: crypto.randomUUID(), respuestas: 0, total: 0, estado: "En curso", ...e }); this.save(); },
  eval_update(id, patch) {
    const arr = this.get().evaluaciones;
    const i = arr.findIndex(x => x.id === id);
    if (i >= 0) { arr[i] = { ...arr[i], ...patch }; this.save(); }
  },
  eval_remove(id) {
    const arr = this.get().evaluaciones;
    const i = arr.findIndex(x => x.id === id);
    if (i >= 0) { arr.splice(i,1); this.save(); }
  }
};

/* ---------------------
   Router por hash
----------------------*/
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

async function loadView() {
  const hash = window.location.hash || '#/dashboard';
  const path = routes[hash] || routes['#/dashboard'];

  // activa item menú
  $$(".menu .item").forEach(a => a.classList.toggle('active', a.getAttribute('href') === hash));

  // migas
  const current = (menu.querySelector(`a[href="${hash}"]`)?.textContent || 'Dashboard').replace(/\s*\d*\s*pendientes/,'').trim();
  crumb.textContent = current;

  try {
    const res = await fetch(path, { cache: 'no-cache' });
    const html = await res.text();
    viewContainer.innerHTML = html;
    initPerView(hash);
  } catch (err) {
    viewContainer.innerHTML = `<div class="card"><h3>Error</h3><p class="muted">No se pudo cargar la vista: ${path}</p></div>`;
  }
}

/* ---------------------
   Controladores de vista
----------------------*/
function initPerView(hash) {
  const data = Store.get();

  // DASHBOARD
  if (hash === '#/dashboard' || hash === '') {
    const totResp = data.evaluaciones.reduce((a,b)=> a + (b.respuestas||0), 0);
    const totCapacidad = data.evaluaciones.reduce((a,b)=> a + (b.total||0), 0) || 1;
    const tasa = Math.round((totResp / totCapacidad) * 100);

    $("#kpi-respondidas") && ($("#kpi-respondidas").textContent = fmt(totResp));
    $("#kpi-tasa") && ($("#kpi-tasa").textContent = `${tasa}% tasa`);

    // Riesgo global (demo): si hay alguna evaluación "Completado" con total>0 y respuestas<total => “Medio”
    const hayAlto = data.casos.some(c => c.prioridad === "Alta");
    $("#kpi-riesgo") && ($("#kpi-riesgo").textContent = hayAlto ? "Medio" : "Bajo");
    $("#kpi-riesgo-pill") && ($("#kpi-riesgo-pill").className = "pill " + (hayAlto ? "warn" : "ok"));

    const casosAbiertos = data.casos.length;
    $("#kpi-casos") && ($("#kpi-casos").textContent = fmt(casosAbiertos));
    $("#kpi-casos-pill") && ($("#kpi-casos-pill").textContent = casosAbiertos >= 3 ? "3 urgentes" : "OK");
    $("#kpi-casos-pill") && ($("#kpi-casos-pill").className = "pill " + (casosAbiertos >= 3 ? "danger" : "ok"));
  }

  // EVALUACIONES
  if (hash === '#/evaluaciones') {
    const tbody = $("#tbl-evaluaciones-body");
    if (tbody) {
      tbody.innerHTML = data.evaluaciones.map(ev => `
        <tr data-id="${ev.id}">
          <td>${ev.instrumento}</td>
          <td>${ev.area}</td>
          <td>${fmt(formatPeriodo(ev.periodo))}</td>
          <td><span class="pill ${ev.estado === "Completado" ? "ok" : "warn"}">${ev.estado}</span></td>
          <td>${ev.respuestas}/${ev.total}</td>
          <td>
            <button class="btn btn-mini" data-act="ver">Ver</button>
            <button class="btn btn-mini" data-act="pdf">PDF</button>
            <button class="btn btn-mini" data-act="fin">Finalizar</button>
            <button class="btn btn-mini" data-act="del">🗑</button>
          </td>
        </tr>
      `).join("");
    }

    // Crear nueva evaluación (form modal del topbar)
    const btnCrear = $("#btn-nueva-eval");
    btnCrear && (btnCrear.onclick = () => openModal("Nueva evaluación"));

    // Acciones por fila
    tbody?.addEventListener("click", (e)=>{
      const btn = e.target.closest("button[data-act]");
      if (!btn) return;
      const tr = e.target.closest("tr");
      const id = tr?.dataset.id;
      const act = btn.dataset.act;

      if (act === "del") {
        if (confirm("¿Eliminar evaluación?")) {
          Store.eval_remove(id);
          loadView();
        }
      }
      if (act === "fin") {
        Store.eval_update(id, { estado: "Completado" });
        loadView();
      }
      if (act === "ver") alert("Ver resultados (Demo)");
      if (act === "pdf") alert("Generar PDF (Demo)");
    });

    // Form inline de creación
    const form = $("#form-eval");
    form?.addEventListener("submit", (e)=>{
      e.preventDefault();
      const f = e.target;
      const instrumento = f.instrumento.value;
      const area = f.area.value.trim();
      const periodo = f.periodo.value;
      const total = parseInt(f.total.value || "0", 10);

      if (!instrumento || !area || !periodo) return alert("Completa los campos requeridos");
      Store.eval_add({ instrumento, area, periodo, total, respuestas: 0 });
      e.target.reset();
      loadView();
    });
  }

  // REPORTES (tabs demo) — se mantiene igual
  if (hash === '#/reportes') {
    const tabs = $("#tabs-reportes");
    const panels = {
      global: $("#panel-global"),
      individuales: $("#panel-individuales"),
      comparativos: $("#panel-comparativos"),
    };
    tabs?.addEventListener('click', (e)=>{
      const btn = e.target.closest('.tab');
      if (!btn) return;
      $$(".tab", tabs).forEach(t=>t.classList.remove('active'));
      btn.classList.add('active');
      const tab = btn.dataset.tab;
      Object.keys(panels).forEach(k => panels[k].style.display = (k===tab?'block':'none'));
    });
  }
}

/* ---------------------
   Helpers
----------------------*/
function formatPeriodo(yyyyMM) {
  // yyyy-mm (input de tipo month) → "Oct-2025"
  if (!yyyyMM) return "";
  const [y,m] = yyyyMM.split("-");
  const meses = ["Ene","Feb","Mar","Abr","May","Jun","Jul","Ago","Sep","Oct","Nov","Dic"];
  return `${meses[(+m||1)-1]}-${y}`;
}

/* ---------------------
   Modal genérico
----------------------*/
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

/* ---------------------
   Topbar / acciones
----------------------*/
document.getElementById('btnNew').onclick = ()=> openModal('Nueva evaluación');
document.getElementById('btnZip').onclick = ()=> alert('Generar Carpeta SUNAFIL (ZIP)… (Demo)');
// Búsqueda (demo)
document.getElementById('q').addEventListener('keydown', (e)=>{
  if (e.key === 'Enter') alert('Buscar: ' + e.target.value + ' (Demo)');
});

/* ---------------------
   Router listeners
----------------------*/
window.addEventListener('hashchange', loadView);
window.addEventListener('DOMContentLoaded', () => {
  Store.load();
  loadView();
});


// Dónde reemplazar más adelante:

// Store (LocalStorage) → 🔁 REEMPLAZAR POR API REAL llamando fetch('/api/evaluaciones'), etc.

// En eval_add, eval_update, eval_remove cambia a POST/PUT/DELETE.

// En initPerView('#/evaluaciones') las acciones ver/pdf pueden ir a tus endpoints de generación de PDF y consulta de resultados.
