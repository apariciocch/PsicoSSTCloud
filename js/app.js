import { initLogin } from './modules/login.js';
import { initNavigation } from './modules/navigation.js';
import { initReportes } from './modules/reportes.js';
import { initUsuarios } from './modules/usuarios.js';
import { initModales } from './modules/modals.js';

document.addEventListener('DOMContentLoaded', () => {
  const crumb = document.getElementById('crumb');
  const searchInput = document.getElementById('q');
  const btnZip = document.getElementById('btnZip');
  const btnNew = document.getElementById('btnNew');
  const btnNuevaEval = document.getElementById('btnNuevaEval');
  const btnNuevaAccion = document.getElementById('btnNuevaAccion');

  const modalApi = initModales();
  const navigation = initNavigation({
    onViewChange: (label) => {
      if (crumb) {
        crumb.textContent = label;
      }
    }
  });

  initReportes();
  initUsuarios();
  initLogin({ onLogout: navigation?.resetActiveView });

  btnNew?.addEventListener('click', () => modalApi.openModal('Nueva evaluación'));
  btnNuevaEval?.addEventListener('click', () => modalApi.openModal('Nueva evaluación'));
  btnNuevaAccion?.addEventListener('click', () => modalApi.openModal('Nueva acción'));

  btnZip?.addEventListener('click', () => {
    alert('Se generará la Carpeta SUNAFIL (ZIP) con informes, plan, actas y evidencias.\n(Demo visual)');
  });

  searchInput?.addEventListener('keydown', (event) => {
    if (event.key === 'Enter') {
      alert(`Buscar: ${searchInput.value}\n(Demo visual)`);
    }
  });
});
