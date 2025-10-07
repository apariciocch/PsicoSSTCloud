// Rutas → vista.html
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

// Cargar vista según hash
async function loadView() {
  const hash = window.location.hash || '#/dashboard';
  const path = routes[hash] || routes['#/dashboard'];

  // Activar item de menú
  [...menu.querySelectorAll('.item')].forEach(a => {
    a.classList.toggle('active', a.getAttribute('href') === hash);
  });

  // Breadcrumb simple
  const current = (menu.querySelector(`a[href="${hash}"]`)?.textContent || 'Dashboard').replace(/\s*\d*\s*pendientes/,'').trim();
  crumb.textContent = current;

  try {
    const res = await fetch(path, { cache: 'no-cache' });
    const html = await res.text();
    viewContainer.innerHTML = html;

    // Inicializaciones por vista (comportamientos mínimos de demo)
    initPerView(hash);
  } catch (err) {
    viewContainer.innerHTML = `<div class="card"><h3>Error</h3><p class="muted">No se pudo cargar la vista: ${path}</p></div>`;
  }
}

// Comportamientos demo por vista
function initPerView(hash) {
  if (hash === '#/reportes') {
    const tabs = document.getElementById('tabs-reportes');
    const panels = {
      global: document.getElementById('panel-global'),
      individuales: document.getElementById('panel-individuales'),
      comparativos: document.getElementById('panel-comparativos'),
    };
    if (tabs) {
      tabs.addEventListener('click', (e)=>{
        const btn = e.target.closest('.tab');
        if (!btn) return;
        tabs.querySelectorAll('.tab').forEach(t=>t.classList.remove('active'));
        btn.classList.add('active');
        const tab = btn.dataset.tab;
        Object.keys(panels).forEach(k => panels[k].style.display = (k===tab?'block':'none'));
      });
    }
  }
}

// Modal genérico
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

// Botones topbar
document.getElementById('btnNew').onclick = ()=> openModal('Nueva evaluación');
document.getElementById('btnZip').onclick = ()=> alert('Se generará la Carpeta SUNAFIL (ZIP) con informes, plan, actas y evidencias. (Demo)');
document.getElementById('modal-close').onclick = closeModal;
document.getElementById('modal-cancel').onclick = closeModal;

// Búsqueda (demo)
document.getElementById('q').addEventListener('keydown', (e)=>{
  if (e.key === 'Enter') alert('Buscar: ' + e.target.value + ' (Demo)');
});

// Router
window.addEventListener('hashchange', loadView);
window.addEventListener('DOMContentLoaded', loadView);
