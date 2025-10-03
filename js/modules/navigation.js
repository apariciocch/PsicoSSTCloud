export function initNavigation({ onViewChange } = {}) {
  const menu = document.getElementById('menu');
  const views = Array.from(document.querySelectorAll('.view'));
  if (!menu || !views.length) return null;

  const activateItem = (item) => {
    if (!item) return;
    menu.querySelectorAll('.item').forEach((node) => node.classList.remove('active'));
    item.classList.add('active');
    const target = item.dataset.view;
    views.forEach((view) => {
      view.classList.toggle('active', view.id === `view-${target}`);
    });
    if (typeof onViewChange === 'function') {
      onViewChange(item.dataset.label || item.textContent.trim());
    }
  };

  menu.addEventListener('click', (event) => {
    const item = event.target.closest('.item');
    if (!item) return;
    activateItem(item);
  });

  const resetActiveView = () => {
    const defaultItem = menu.querySelector('.item[data-view="dashboard"]') || menu.querySelector('.item');
    activateItem(defaultItem);
  };

  resetActiveView();

  return { resetActiveView };
}
