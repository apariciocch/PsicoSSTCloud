const STORAGE_KEY = 'psicosstcloud:evaluaciones';
const MONTH_LABELS = ['Ene', 'Feb', 'Mar', 'Abr', 'May', 'Jun', 'Jul', 'Ago', 'Sep', 'Oct', 'Nov', 'Dic'];

function createId() {
  return `eval-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 7)}`;
}

function formatPeriod(value) {
  if (!value || typeof value !== 'string') return 'Por definir';
  const [year, month] = value.split('-').map(Number);
  if (!year || !month || month < 1 || month > 12) return 'Por definir';
  return `${MONTH_LABELS[month - 1]}-${year}`;
}

function loadEvaluaciones() {
  if (typeof window === 'undefined' || !window.localStorage) return null;
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw);
    if (!Array.isArray(parsed)) return null;
    return parsed;
  } catch (error) {
    console.warn('No se pudieron cargar las evaluaciones almacenadas', error);
    return null;
  }
}

function persistEvaluaciones(evaluaciones) {
  if (typeof window === 'undefined' || !window.localStorage) return;
  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(evaluaciones));
  } catch (error) {
    console.warn('No se pudieron guardar las evaluaciones', error);
  }
}

function renderEmptyState(tableBody) {
  tableBody.innerHTML = `<tr><td colspan="5" class="muted">Aún no has creado evaluaciones con el MINSA 2024.</td></tr>`;
}

export function initEvaluaciones({ modalApi }) {
  const tableBody = document.getElementById('evaluaciones-body');
  const form = document.getElementById('evaluation-form');
  const areaInput = document.getElementById('evaluation-area');
  const periodInput = document.getElementById('evaluation-period');
  const modeSelect = document.getElementById('evaluation-mode');
  const notesInput = document.getElementById('evaluation-notes');
  const summaryTemplate = document.getElementById('evaluation-summary-template');
  const countBadge = document.getElementById('evaluaciones-count');

  const defaultEvaluaciones = [
    {
      id: createId(),
      instrument: 'MINSA 2024',
      area: 'Atención al Cliente',
      period: '2025-10',
      status: 'En curso',
      mode: 'Anónimo',
      notes: 'Seguimiento semanal con líderes de área.'
    },
    {
      id: createId(),
      instrument: 'MINSA 2024',
      area: 'Operaciones',
      period: '2025-08',
      status: 'Completado',
      mode: 'Identificado',
      notes: 'Resultados consolidados y compartidos con SST.'
    }
  ];

  let evaluaciones = loadEvaluaciones() ?? defaultEvaluaciones;

  const updateBadge = () => {
    if (!countBadge) return;
    const total = evaluaciones.length;
    countBadge.textContent = total === 1 ? '1 activa' : `${total} activas`;
  };

  const renderEvaluaciones = () => {
    if (!tableBody) return;

    if (!evaluaciones.length) {
      renderEmptyState(tableBody);
      updateBadge();
      return;
    }

    tableBody.innerHTML = evaluaciones
      .map((item) => {
        const periodLabel = formatPeriod(item.period);
        const statusClass = item.status === 'Completado' ? 'ok' : item.status === 'En curso' ? 'warn' : 'info';
        const statusLabel = item.status;
        return `
          <tr data-id="${item.id}">
            <td>${item.instrument}</td>
            <td>${item.area}</td>
            <td>${periodLabel}</td>
            <td><span class="pill ${statusClass}">${statusLabel}</span></td>
            <td><button class="btn" type="button" data-action="ver" data-id="${item.id}">Ver detalle</button></td>
          </tr>
        `;
      })
      .join('');

    updateBadge();
  };

  const resetForm = () => {
    if (form) {
      form.reset();
    }
    if (areaInput) areaInput.value = '';
    if (periodInput) periodInput.value = '';
    if (modeSelect) modeSelect.value = 'Anónimo';
    if (notesInput) notesInput.value = '';
  };

  const openCreation = () => {
    resetForm();
    if (modalApi?.openModal) {
      modalApi.openModal('Nueva evaluación MINSA 2024');
    }
    const focusArea = () => areaInput?.focus();
    if (typeof window !== 'undefined' && typeof window.requestAnimationFrame === 'function') {
      window.requestAnimationFrame(focusArea);
    } else {
      focusArea();
    }
  };

  form?.addEventListener('submit', (event) => {
    event.preventDefault();

    const area = areaInput?.value?.trim();
    const period = periodInput?.value;
    const mode = modeSelect?.value ?? 'Anónimo';
    const notes = notesInput?.value?.trim();

    if (!area) {
      areaInput?.focus();
      return;
    }

    if (!period) {
      periodInput?.focus();
      return;
    }

    const nuevaEvaluacion = {
      id: createId(),
      instrument: 'MINSA 2024',
      area,
      period,
      status: 'En curso',
      mode,
      notes
    };

    evaluaciones = [nuevaEvaluacion, ...evaluaciones];
    persistEvaluaciones(evaluaciones);
    renderEvaluaciones();

    if (modalApi?.closeModal) {
      modalApi.closeModal();
    }

    if (summaryTemplate) {
      const resumen = summaryTemplate.content.cloneNode(true);
      const areaEl = resumen.querySelector('[data-summary="area"]');
      const periodoEl = resumen.querySelector('[data-summary="periodo"]');
      const modoEl = resumen.querySelector('[data-summary="modo"]');
      const notasEl = resumen.querySelector('[data-summary="notas"]');
      if (areaEl) areaEl.textContent = area;
      if (periodoEl) periodoEl.textContent = formatPeriod(period);
      if (modoEl) modoEl.textContent = mode;
      if (notasEl) notasEl.textContent = notes || 'Sin notas adicionales.';
      const container = document.getElementById('evaluation-toast');
      if (container) {
        container.innerHTML = '';
        container.appendChild(resumen);
        container.classList.add('visible');
        setTimeout(() => container.classList.remove('visible'), 4200);
      }
    }
  });

  tableBody?.addEventListener('click', (event) => {
    const target = event.target;
    if (!(target instanceof HTMLElement)) return;
    if (target.dataset.action !== 'ver') return;
    const { id } = target.dataset;
    if (!id) return;
    const evalSeleccionada = evaluaciones.find((item) => item.id === id);
    if (!evalSeleccionada) return;
    const detalle = [
      `Instrumento: ${evalSeleccionada.instrument}`,
      `Área: ${evalSeleccionada.area}`,
      `Periodo: ${formatPeriod(evalSeleccionada.period)}`,
      `Modo: ${evalSeleccionada.mode ?? 'Anónimo'}`,
      `Estado: ${evalSeleccionada.status}`,
      `Notas: ${evalSeleccionada.notes?.trim() ? evalSeleccionada.notes : 'Sin notas registradas.'}`
    ].join('\n');
    alert(detalle);
  });

  renderEvaluaciones();

  return { openCreation };
}
