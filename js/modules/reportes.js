export function initReportes() {
  const tabs = document.getElementById('tabs-reportes');
  if (!tabs) return;

  const panels = {
    global: document.getElementById('panel-global'),
    individuales: document.getElementById('panel-individuales'),
    comparativos: document.getElementById('panel-comparativos')
  };

  const showTab = (tab) => {
    Object.entries(panels).forEach(([key, panel]) => {
      if (!panel) return;
      panel.style.display = key === tab ? 'block' : 'none';
    });
  };

  tabs.addEventListener('click', (event) => {
    const tab = event.target.closest('.tab');
    if (!tab) return;
    tabs.querySelectorAll('.tab').forEach((node) => node.classList.remove('active'));
    tab.classList.add('active');
    showTab(tab.dataset.tab);
  });
}
