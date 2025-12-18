// Navegación y carga dinámica de vistas
const menu = document.getElementById('menu');
const crumb = document.getElementById('crumb');
const main = document.getElementById('main-content');

const VIEW_PATHS = {
  dashboard: 'modules/dashboard/view.html',
  evaluaciones: 'modules/evaluaciones/view.html',
  'evaluaciones-minsa': 'modules/evaluaciones/minsa_form.html',
  evaluado: 'modules/evaluado/view.html',
  reportes: 'modules/reportes/view.html',
  carpeta: 'modules/carpeta/view.html',
  plan: 'modules/plan/view.html',
  capacitaciones: 'modules/capacitaciones/view.html',
  casos: 'modules/casos/view.html',
  usuarios: 'modules/usuarios/view.html',
  config: 'modules/config/view.html',
  resultados: 'modules/resultados/view.html',
  login: 'modules/login/view.html',
};

function isAuthed(){
  return localStorage.getItem('psico_auth') === 'admin';
}

function applyAuthLayout(){
  if(isAuthed()){
    document.body.classList.remove('unauth');
  }else{
    document.body.classList.add('unauth');
  }
}

async function loadView(name){
  const path = VIEW_PATHS[name];
  try{
    const res = await fetch(path, {cache:'no-store'});
    const html = await res.text();
    main.innerHTML = html;
    bindViewSpecific();
  }catch(err){
    main.innerHTML = `<div class="card"><h3>Error cargando vista</h3><p class="muted">${name}</p></div>`;
    console.error('Error al cargar vista', name, err);
  }
}

// Manejo de menú
menu.addEventListener('click', (e)=>{
  const item = e.target.closest('.item');
  if(!item) return;
  const target = item.dataset.view;
  if(!isAuthed() && target!=='evaluado'){
    crumb.textContent = 'Login';
    applyAuthLayout();
    loadView('login');
    return;
  }
  menu.querySelectorAll('.item').forEach(i=>i.classList.remove('active'));
  item.classList.add('active');
  crumb.textContent = item.textContent.replace(/\s*\d*\s*pendientes/,'').trim();
  loadView(target);
});

// Tabs en Reportes (delegación)
document.addEventListener('click', (e)=>{
  const t = e.target.closest('.tab');
  if(!t || !t.dataset.tab) return;
  const tabsContainer = t.parentElement;
  if(!tabsContainer) return;
  tabsContainer.querySelectorAll('.tab').forEach(x=>x.classList.remove('active'));
  t.classList.add('active');
  const tab = t.dataset.tab;
  const panelGlobal = document.getElementById('panel-global');
  const panelInd = document.getElementById('panel-individuales');
  const panelComp = document.getElementById('panel-comparativos');
  if(panelGlobal && panelInd && panelComp){
    panelGlobal.style.display = tab==='global'?'block':'none';
    panelInd.style.display = tab==='individuales'?'block':'none';
    panelComp.style.display = tab==='comparativos'?'block':'none';
  }
});

// Modal genérico
const modalBg = document.getElementById('modal-bg');
const openModal = (title='Nueva evaluación')=>{
  if(!isAuthed()){
    crumb.textContent = 'Login';
    applyAuthLayout();
    loadView('login');
    return;
  }
  document.getElementById('modal-title').textContent = title;
  modalBg.style.display='grid';
};
const closeModal = ()=> modalBg.style.display='none';
document.getElementById('modal-close').onclick = closeModal;
document.getElementById('modal-cancel').onclick = closeModal;

document.getElementById('btnNew').onclick = ()=> openModal('Nueva evaluación');
document.getElementById('btnZip').onclick = ()=>{
  alert('Se generará la Carpeta SUNAFIL (ZIP) con informes, plan, actas y evidencias.\n(Demo visual)');
};

// Logout
document.getElementById('btnLogout').onclick = ()=>{
  localStorage.removeItem('psico_auth');
  menu.querySelectorAll('.item').forEach(i=>i.classList.remove('active'));
  crumb.textContent = 'Login';
  applyAuthLayout();
  loadView('login');
};

// Búsqueda (demo)
const q = document.getElementById('q');
q.addEventListener('keydown', (e)=>{
  if(e.key==='Enter') alert('Buscar: ' + q.value + '\n(Demo visual)');
});

// Enlaces específicos de cada vista (si existen)
function bindViewSpecific(){
  // --- Evaluado: seleccionar evaluación y abrir MINSA ---
  const evalSelect = document.getElementById('eval-select');
  const btnIniciarEval = document.getElementById('btnIniciarEval');
  const evaladoMsg = document.getElementById('evalado-msg');
  function loadEvals(){
    try{ return JSON.parse(localStorage.getItem('psico_evals')||'[]'); }catch{ return [] }
  }
  if(evalSelect){
    const evals = loadEvals();
    evalSelect.innerHTML = evals.length ? evals.map(e=>`<option value="${e.id}">${e.nombre} — ${e.area||'-'} (${e.periodo||'-'})</option>`).join('') : '<option value="">No hay evaluaciones disponibles</option>';
  }
  if(btnIniciarEval){
    btnIniciarEval.onclick = ()=>{
      const id = evalSelect?.value;
      if(!id){ if(evaladoMsg){ evaladoMsg.style.display='inline-block'; evaladoMsg.textContent='Selecciona una evaluación'; } return; }
      localStorage.setItem('psico_current_eval', id);
      localStorage.setItem('psico_form_origin','evaluado');
      crumb.textContent = 'Evaluado';
      loadView('evaluaciones-minsa');
    };
  }
  // --- Evaluaciones: CRUD + abrir formulario MINSA ---
  const evalsList = document.getElementById('evals-list');
  const btnGuardarEval = document.getElementById('btnGuardarEval');
  const evalNombreEl = document.getElementById('eval-nombre');
  const evalAreaEl = document.getElementById('eval-area');
  const evalPeriodoEl = document.getElementById('eval-periodo');
  const evalMsgEl = document.getElementById('eval-msg');
  const evalCountEl = document.getElementById('eval-count');
  const evalFiltroEl = document.getElementById('eval-filtro');
  function loadEvals(){
    try{ return JSON.parse(localStorage.getItem('psico_evals')||'[]'); }catch{ return [] }
  }
  function saveEvals(arr){
    localStorage.setItem('psico_evals', JSON.stringify(arr));
  }
  function renderEvals(){
    if(!evalsList) return;
    const all = loadEvals();
    const term = (evalFiltroEl?.value||'').toLowerCase();
    const rows = all.filter(e=> !term || e.area.toLowerCase().includes(term) || e.nombre.toLowerCase().includes(term));
    if(evalCountEl) evalCountEl.textContent = `${all.length} registradas`;
    if(rows.length===0){
      evalsList.innerHTML = `<tr><td colspan="5" class="muted">Sin evaluaciones</td></tr>`;
      return;
    }
    evalsList.innerHTML = rows.map(e=>{
      const estadoPill = e.estado==='Completado'? 'ok' : (e.estado==='En curso'? 'warn':'');
      return `<tr>
        <td>${e.nombre}</td>
        <td>${e.area||'-'}</td>
        <td>${e.periodo||'-'}</td>
        <td><span class="pill ${estadoPill}">${e.estado||'-'}</span></td>
        <td>
          <button class="btn" data-eval-aplicar="${e.id}">Aplicar</button>
          <button class="btn" data-eval-del="${e.id}">Eliminar</button>
        </td>
      </tr>`;
    }).join('');
  }
  if(evalsList){
    renderEvals();
    evalsList.addEventListener('click', (ev)=>{
      const ap = ev.target.closest('[data-eval-aplicar]');
      const del = ev.target.closest('[data-eval-del]');
      if(ap){
        const id = ap.getAttribute('data-eval-aplicar');
        localStorage.setItem('psico_current_eval', id);
        crumb.textContent = 'Evaluación';
        loadView('evaluaciones-minsa');
      }else if(del){
        const id = del.getAttribute('data-eval-del');
        const evals = loadEvals().filter(x=> String(x.id)!==String(id));
        saveEvals(evals);
        renderEvals();
      }
    });
    if(evalFiltroEl){ evalFiltroEl.oninput = renderEvals; }
  }
  if(btnGuardarEval){
    btnGuardarEval.onclick = ()=>{
      const nombre = (evalNombreEl?.value||'').trim();
      const area = (evalAreaEl?.value||'').trim();
      const periodo = evalPeriodoEl?.value||'';
      if(!nombre){ if(evalMsgEl){ evalMsgEl.textContent='Ingresa el nombre'; evalMsgEl.classList.add('danger'); evalMsgEl.style.display='inline-block'; } return; }
      const evals = loadEvals();
      const id = Date.now();
      evals.push({id, nombre, area, periodo, estado:'En curso'});
      saveEvals(evals);
      if(evalMsgEl){ evalMsgEl.textContent='Guardado'; evalMsgEl.classList.remove('danger'); evalMsgEl.style.display='inline-block'; setTimeout(()=> evalMsgEl.style.display='none', 1500); }
      renderEvals();
      if(evalNombreEl) evalNombreEl.value='';
      if(evalAreaEl) evalAreaEl.value='';
      if(evalPeriodoEl) evalPeriodoEl.value='';
    };
    const btnNuevaEval = document.getElementById('btnNuevaEval');
    if(btnNuevaEval){ btnNuevaEval.onclick = ()=> evalNombreEl?.focus(); }
  }

  // --- Formulario MINSA: guardado y exportación ---
  const minsaForm = document.getElementById('minsa-form');
  const btnMinsaGuardar = document.getElementById('btnMinsaGuardar');
  const btnMinsaExport = document.getElementById('btnMinsaExport');
  const btnMinsaVolver = document.getElementById('btnMinsaVolver');
  const origin = localStorage.getItem('psico_form_origin');
  if(origin==='evaluado' && btnMinsaExport){ btnMinsaExport.style.display='none'; }
  function loadResp(){ try{ return JSON.parse(localStorage.getItem('psico_respuestas')||'[]'); }catch{ return [] } }
  function saveResp(arr){ localStorage.setItem('psico_respuestas', JSON.stringify(arr)); }
  function collectAnswers(){
    const ans = {};
    for(let i=1;i<=44;i++){
      const sel = document.querySelector(`input[name="q${i}"]:checked`);
      ans[`q${i}`] = sel ? Number(sel.value) : '';
    }
    return ans;
  }
  function currentEval(){
    const id = localStorage.getItem('psico_current_eval');
    const all = loadEvals();
    return all.find(x=> String(x.id)===String(id));
  }
  if(btnMinsaVolver){
    btnMinsaVolver.onclick = ()=>{
      const origin = localStorage.getItem('psico_form_origin');
      if(origin==='evaluado'){
        crumb.textContent = 'Evaluado';
        loadView('evaluado');
      }else{
        crumb.textContent = 'Evaluaciones';
        loadView('evaluaciones');
      }
      localStorage.removeItem('psico_form_origin');
    };
  }
  if(btnMinsaGuardar){
    btnMinsaGuardar.onclick = ()=>{
      const ev = currentEval();
      const ans = collectAnswers();
      const meta = {
        edad: (document.getElementById('f-edad')?.value||''),
        genero: (document.getElementById('f-genero')?.value||''),
        puesto: (document.getElementById('f-puesto')?.value||''),
        fecha: (document.getElementById('f-fecha')?.value||'')
      };
      const arr = loadResp();
      arr.push({evalId: ev?.id || null, nombre: ev?.nombre || 'MINSA 2024', respuestas: ans, meta, ts: Date.now()});
      saveResp(arr);
      alert('Respuestas guardadas');
    };
  }
  if(btnMinsaExport){
    btnMinsaExport.onclick = ()=>{
      const ev = currentEval();
      const ans = collectAnswers();
      const meta = {
        edad: (document.getElementById('f-edad')?.value||''),
        genero: (document.getElementById('f-genero')?.value||''),
        puesto: (document.getElementById('f-puesto')?.value||''),
        fecha: (document.getElementById('f-fecha')?.value||'')
      };
      // Generar CSV con valores numéricos Q1..Q44
      const headers = ['Evaluación','Área','Periodo','Edad','Género','Puesto','Fecha', ...Array.from({length:44},(_,i)=>`Q${i+1}`)];
      const row = [ev?.nombre||'', ev?.area||'', ev?.periodo||'', meta.edad, meta.genero, meta.puesto, meta.fecha, ...Array.from({length:44},(_,i)=> ans[`q${i+1}`])];
      const csv = [headers.join(','), row.join(',')].join('\n');
      const blob = new Blob([csv], {type:'text/csv;charset=utf-8;'});
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `evaluacion_minsa_${ev?.area||'general'}_${Date.now()}.csv`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    };
  }

  // --- Resultados: selección y cálculo ---
  const resSelect = document.getElementById('res-select');
  const btnVerResultado = document.getElementById('btnVerResultado');
  const resTotalEl = document.getElementById('res-total');
  const resDimsEl = document.getElementById('res-dimensiones');

  function dimConfig(){
    return [
      {key:'carga_trabajo', label:'Carga de trabajo', items:[1,2]},
      {key:'desarrollo_competencias', label:'Desarrollo de competencias', items:[3,4]},
      {key:'liderazgo', label:'Liderazgo', items:[5,6,7,8]},
      {key:'claridad_rol', label:'Claridad del rol', items:[9,10]},
      {key:'organizacion_trabajo', label:'Organización del trabajo', items:[11,12,13,14]},
      {key:'exigencias_trabajo', label:'Exigencias del trabajo', items:[15,16,17,18]},
      {key:'soporte_apoyo', label:'Soporte y apoyo', items:[19,20,21,22]},
      {key:'apartamiento', label:'Apartamiento', items:[23,24,25,26]},
      {key:'inclusion_laboral', label:'Inclusión laboral', items:[27,28]},
      {key:'motivacion_trabajo', label:'Motivación en el trabajo', items:[29,30]},
      {key:'autonomia_trabajo', label:'Autonomía del trabajo', items:[31,32,33,34]},
      {key:'condiciones_trabajo', label:'Condiciones de trabajo', items:[35,36]},
      {key:'interfase_trabajo_familia', label:'Interfase trabajo/familia', items:[37,38]},
      {key:'seguridad', label:'Seguridad', items:[39,40]},
      {key:'salud_percibida', label:'Salud percibida', items:[41,42]},
      {key:'satisfaccion', label:'Satisfacción con el trabajo', items:[43,44]}
    ];
  }
  function classify(sum, min, max){
    const range = max - min + 1;
    const band = Math.floor(range/3);
    const lowMax = min + band - 1;
    const midMax = lowMax + band;
    if(sum <= lowMax) return 'Alto riesgo';
    if(sum <= midMax) return 'Medio';
    return 'Bajo';
  }
  function calculateRisk(respuestas){
    const faltantes = [];
    for(let i=1;i<=44;i++){
      const v = respuestas[`q${i}`];
      if(v==='' || v===undefined || isNaN(Number(v))) faltantes.push(i);
    }
    if(faltantes.length){
      return {estado:'incompleto', faltantes};
    }
    const dims = dimConfig();
    const puntajesPorDimension = {};
    const nivelPorDimension = {};
    let total = 0;
    dims.forEach(d=>{
      const sum = d.items.reduce((acc,i)=> acc + Number(respuestas[`q${i}`]), 0);
      const min = d.items.length * 1;
      const max = d.items.length * 4;
      puntajesPorDimension[d.key] = sum;
      nivelPorDimension[d.key] = classify(sum, min, max);
      total += sum;
    });
    const nivelGeneral = classify(total, 44*1, 44*4);
    return {estado:'ok', puntajesPorDimension, nivelPorDimension, puntajeTotal: total, nivelGeneral};
  }
  try{
    const minAns = Object.fromEntries(Array.from({length:44},(_,i)=>[`q${i+1}`,1]));
    const maxAns = Object.fromEntries(Array.from({length:44},(_,i)=>[`q${i+1}`,4]));
    const midAns = Object.fromEntries(Array.from({length:44},(_,i)=>[`q${i+1}`,2]));
    console.log('[TEST] MINSA min', calculateRisk(minAns));
    console.log('[TEST] MINSA max', calculateRisk(maxAns));
    console.log('[TEST] MINSA mid', calculateRisk(midAns));
  }catch(e){ console.warn('Pruebas MINSA no ejecutadas', e); }

  function loadAllResp(){ try{ return JSON.parse(localStorage.getItem('psico_respuestas')||'[]'); }catch{ return [] } }
  if(resSelect){
    const arr = loadAllResp();
    resSelect.innerHTML = arr.length ? arr.map((r,i)=>{
      const evName = r.nombre || 'MINSA 2024';
      const metaTxt = [r.meta?.puesto, r.meta?.fecha].filter(Boolean).join(' · ');
      return `<option value="${i}">${evName} — ${metaTxt || 'sin meta'}</option>`;
    }).join('') : '<option value="">Sin respuestas</option>';
  }
  function renderResultado(r){
    if(!r || !resTotalEl || !resDimsEl) return;
    const calc = calculateRisk(r.respuestas||{});
    if(calc.estado!=='ok'){
      resTotalEl.innerHTML = `<span>Puntaje total</span><b>—</b><span class="pill danger">Incompleto: faltan ${calc.faltantes.join(', ')}</span>`;
      resDimsEl.innerHTML = '';
      return;
    }
    resTotalEl.innerHTML = `<span>Puntaje total</span><b>${calc.puntajeTotal}</b><span class="pill">${calc.nivelGeneral}</span>`;
    const dims = dimConfig();
    resDimsEl.innerHTML = dims.map(d=>{
      const score = calc.puntajesPorDimension[d.key];
      const niv = calc.nivelPorDimension[d.key];
      const pillClass = niv==='Alto riesgo' ? 'danger' : (niv==='Medio' ? 'warn' : 'ok');
      return `<div class="card" style="grid-column: span 6">
        <div class="kpi"><span>${d.label}</span><b>${score}</b><span class="pill ${pillClass}">${niv}</span></div>
      </div>`;
    }).join('');
  }
  if(btnVerResultado){
    btnVerResultado.onclick = ()=>{
      const idx = resSelect?.value;
      if(idx==='' || idx===undefined){ alert('Selecciona una respuesta'); return; }
      const arr = loadAllResp();
      const r = arr[Number(idx)];
      renderResultado(r);
    };
  }

  // --- Configuración: Empresas ---
  const btnGuardarEmpresa = document.getElementById('btnGuardarEmpresa');
  const listEmpresas = document.getElementById('empresas-list');
  function loadEmpresas(){
    try{
      return JSON.parse(localStorage.getItem('psico_empresas')||'[]');
    }catch{ return [] }
  }
  function saveEmpresas(arr){
    localStorage.setItem('psico_empresas', JSON.stringify(arr));
  }
  function renderEmpresas(){
    if(!listEmpresas) return;
    const empresas = loadEmpresas();
    if(empresas.length===0){
      listEmpresas.innerHTML = `<tr><td colspan="4" class="muted">Sin empresas registradas</td></tr>`;
      return;
    }
    listEmpresas.innerHTML = empresas.map((e,i)=>
      `<tr>
        <td>${e.nombre}</td>
        <td>${e.ruc}</td>
        <td>${e.instrumento||'-'}</td>
        <td><button class="btn" data-del-index="${i}">Eliminar</button></td>
      </tr>`
    ).join('');
  }
  if(listEmpresas){
    renderEmpresas();
    listEmpresas.addEventListener('click',(ev)=>{
      const btn = ev.target.closest('[data-del-index]');
      if(!btn) return;
      const idx = parseInt(btn.getAttribute('data-del-index'),10);
      const empresas = loadEmpresas();
      if(idx>=0 && idx<empresas.length){
        empresas.splice(idx,1);
        saveEmpresas(empresas);
        renderEmpresas();
      }
    });
  }
  if(btnGuardarEmpresa){
    const nombreEl = document.getElementById('empresa-nombre');
    const rucEl = document.getElementById('empresa-ruc');
    const instrEl = document.getElementById('empresa-instrumento');
    const firmaEl = document.getElementById('empresa-firma');
    const msgEl = document.getElementById('empresa-msg');
    btnGuardarEmpresa.onclick = ()=>{
      const nombre = (nombreEl?.value||'').trim();
      const ruc = (rucEl?.value||'').trim();
      const instrumento = instrEl?.value||'';
      const firma = (firmaEl?.value||'').trim();
      if(!nombre || !ruc){
        if(msgEl){ msgEl.textContent='Completa nombre y RUC'; msgEl.classList.add('danger'); msgEl.style.display='inline-block'; }
        return;
      }
      const empresas = loadEmpresas();
      empresas.push({nombre, ruc, instrumento, firma, ts: Date.now()});
      saveEmpresas(empresas);
      if(msgEl){ msgEl.textContent='Guardado'; msgEl.classList.remove('danger'); msgEl.style.display='inline-block'; setTimeout(()=> msgEl.style.display='none', 1500); }
      renderEmpresas();
    };
  }
  // Login form
  const form = document.getElementById('login-form');
  if(form){
    crumb.textContent = 'Login';
    form.addEventListener('submit', (e)=>{
      e.preventDefault();
      const u = (document.getElementById('login-user')?.value || '').trim();
      const p = (document.getElementById('login-pass')?.value || '').trim();
      const err = document.getElementById('login-error');
      if(u==='admin' && p==='admin'){
        localStorage.setItem('psico_auth','admin');
        if(err) err.style.display='none';
        crumb.textContent = 'Dashboard';
        applyAuthLayout();
        const first = menu.querySelector('.item[data-view="dashboard"]');
        if(first){
          menu.querySelectorAll('.item').forEach(i=>i.classList.remove('active'));
          first.classList.add('active');
        }
        loadView('dashboard');
      }else{
        if(err) err.style.display='block';
      }
    });
  }
  const btnNuevaEval = document.getElementById('btnNuevaEval');
  if (btnNuevaEval) btnNuevaEval.onclick = ()=> openModal('Nueva evaluación');
  const btnNuevaAccion = document.getElementById('btnNuevaAccion');
  if (btnNuevaAccion) btnNuevaAccion.onclick = ()=> openModal('Nueva acción');
  const btnNuevaCap = document.getElementById('btnNuevaCap');
  if (btnNuevaCap) btnNuevaCap.onclick = ()=> openModal('Nueva capacitación');

  // --- Capacitaciones: CRUD y adjunto ---
  const capsList = document.getElementById('caps-list');
  const btnGuardarCap = document.getElementById('btnGuardarCap');
  const capCursoEl = document.getElementById('cap-curso');
  const capFechaEl = document.getElementById('cap-fecha');
  const capAreaEl = document.getElementById('cap-area');
  const capInsEl = document.getElementById('cap-inscritos');
  const capFileEl = document.getElementById('cap-file');
  const capFileNameEl = document.getElementById('cap-file-name');
  const capMsgEl = document.getElementById('cap-msg');
  function loadCaps(){
    try{ return JSON.parse(localStorage.getItem('psico_caps')||'[]'); }catch{ return [] }
  }
  function saveCaps(arr){
    localStorage.setItem('psico_caps', JSON.stringify(arr));
  }
  function renderCaps(){
    if(!capsList) return;
    const caps = loadCaps();
    if(caps.length===0){
      capsList.innerHTML = `<tr><td colspan="6" class="muted">Sin capacitaciones registradas</td></tr>`;
      return;
    }
    capsList.innerHTML = caps.map((c,i)=>{
      const hasDoc = c.docName && c.docData;
      const docCell = hasDoc ? `<a class="btn" href="${c.docData}" download="${c.docName}">Descargar</a>` : `<span class="muted">Sin documento</span>`;
      return `<tr>
        <td>${c.curso}</td>
        <td>${c.fecha||'-'}</td>
        <td>${c.area||'-'}</td>
        <td>${c.inscritos||0}</td>
        <td>${docCell}</td>
        <td>
          <button class="btn" data-cap-del="${i}">Eliminar</button>
        </td>
      </tr>`;
    }).join('');
  }
  if(capsList){
    renderCaps();
    capsList.addEventListener('click',(ev)=>{
      const delBtn = ev.target.closest('[data-cap-del]');
      if(!delBtn) return;
      const idx = parseInt(delBtn.getAttribute('data-cap-del'),10);
      const caps = loadCaps();
      if(idx>=0 && idx<caps.length){
        caps.splice(idx,1);
        saveCaps(caps);
        renderCaps();
      }
    });
  }
  if(capFileEl && capFileNameEl){
    capFileEl.addEventListener('change', ()=>{
      const f = capFileEl.files?.[0];
      capFileNameEl.textContent = f ? `${f.name} (${Math.round(f.size/1024)} KB)` : '';
    });
  }
  if(btnGuardarCap){
    btnGuardarCap.onclick = ()=>{
      const curso = (capCursoEl?.value||'').trim();
      const fecha = capFechaEl?.value||'';
      const area = (capAreaEl?.value||'').trim();
      const inscritos = parseInt(capInsEl?.value||'0',10) || 0;
      if(!curso){
        if(capMsgEl){ capMsgEl.textContent='Completa el nombre del curso'; capMsgEl.classList.add('danger'); capMsgEl.style.display='inline-block'; }
        return;
      }
      const addCap = (docName, docData)=>{
        const caps = loadCaps();
        caps.push({curso, fecha, area, inscritos, docName, docData, ts: Date.now()});
        saveCaps(caps);
        if(capMsgEl){ capMsgEl.textContent='Guardado'; capMsgEl.classList.remove('danger'); capMsgEl.style.display='inline-block'; setTimeout(()=> capMsgEl.style.display='none', 1500); }
        renderCaps();
        if(capCursoEl) capCursoEl.value='';
        if(capFechaEl) capFechaEl.value='';
        if(capAreaEl) capAreaEl.value='';
        if(capInsEl) capInsEl.value='';
        if(capFileEl) capFileEl.value='';
        if(capFileNameEl) capFileNameEl.textContent='';
      };
      const file = capFileEl?.files?.[0];
      if(file){
        // Convertir a dataURL (atención: límite de tamaño de localStorage ~5MB)
        const reader = new FileReader();
        reader.onload = ()=> addCap(file.name, reader.result);
        reader.onerror = ()=> addCap(undefined, undefined);
        reader.readAsDataURL(file);
      }else{
        addCap(undefined, undefined);
      }
    };
  }

  // --- Usuarios: CRUD simple ---
  const usersList = document.getElementById('users-list');
  const btnGuardarUser = document.getElementById('btnGuardarUser');
  const btnNuevoUser = document.getElementById('btnNuevoUser');
  function loadUsuarios(){
    try{ return JSON.parse(localStorage.getItem('psico_usuarios')||'[]'); }catch{ return [] }
  }
  function saveUsuarios(arr){
    localStorage.setItem('psico_usuarios', JSON.stringify(arr));
  }
  function renderUsuarios(){
    if(!usersList) return;
    const usuarios = loadUsuarios();
    if(usuarios.length===0){
      usersList.innerHTML = `<tr><td colspan="8" class="muted">Sin usuarios registrados</td></tr>`;
      return;
    }
    usersList.innerHTML = usuarios.map((u,i)=>
      `<tr>
        <td>${u.nombre}</td>
        <td>${u.dni||'-'}</td>
        <td>${u.correo}</td>
        <td>${u.genero||'-'}</td>
        <td>${u.puesto||'-'}</td>
        <td>${u.rol}</td>
        <td><span class="pill ${u.estado==='Activo'?'ok':''}">${u.estado}</span></td>
        <td><button class="btn" data-user-del="${i}">Eliminar</button></td>
      </tr>`
    ).join('');
  }
  if(usersList){
    renderUsuarios();
    usersList.addEventListener('click', (ev)=>{
      const btn = ev.target.closest('[data-user-del]');
      if(!btn) return;
      const idx = parseInt(btn.getAttribute('data-user-del'),10);
      const usuarios = loadUsuarios();
      if(idx>=0 && idx<usuarios.length){
        usuarios.splice(idx,1);
        saveUsuarios(usuarios);
        renderUsuarios();
      }
    });
  }
  if(btnGuardarUser){
    const nombreEl = document.getElementById('user-nombre');
    const dniEl = document.getElementById('user-dni');
    const correoEl = document.getElementById('user-correo');
    const generoEl = document.getElementById('user-genero');
    const puestoEl = document.getElementById('user-puesto');
    const rolEl = document.getElementById('user-rol');
    const estadoEl = document.getElementById('user-estado');
    const msgEl = document.getElementById('user-msg');
    btnGuardarUser.onclick = ()=>{
      const nombre = (nombreEl?.value||'').trim();
      const dni = (dniEl?.value||'').trim();
      const correo = (correoEl?.value||'').trim();
      const genero = generoEl?.value||'';
      const puesto = (puestoEl?.value||'').trim();
      const rol = rolEl?.value||'';
      const estado = estadoEl?.value||'Invitado';
      if(!nombre || !correo){
        if(msgEl){ msgEl.textContent='Completa nombre y correo'; msgEl.classList.add('danger'); msgEl.style.display='inline-block'; }
        return;
      }
      const usuarios = loadUsuarios();
      usuarios.push({nombre, dni, correo, genero, puesto, rol, estado, ts: Date.now()});
      saveUsuarios(usuarios);
      if(msgEl){ msgEl.textContent='Guardado'; msgEl.classList.remove('danger'); msgEl.style.display='inline-block'; setTimeout(()=> msgEl.style.display='none', 1500); }
      renderUsuarios();
      if(nombreEl) nombreEl.value='';
      if(dniEl) dniEl.value='';
      if(correoEl) correoEl.value='';
      if(generoEl) generoEl.value='';
      if(puestoEl) puestoEl.value='';
    };
  }
  if(btnNuevoUser){
    btnNuevoUser.onclick = ()=>{
      document.getElementById('user-nombre')?.focus();
    };
  }
}

// Carga inicial
applyAuthLayout();
if(isAuthed()){
  crumb.textContent = 'Dashboard';
  loadView('dashboard');
}else{
  crumb.textContent = 'Login';
  loadView('login');
}
