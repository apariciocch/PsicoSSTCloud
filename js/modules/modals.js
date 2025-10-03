export function initModales() {
  const modalBg = document.getElementById('modal-bg');
  const modalTitle = document.getElementById('modal-title');
  const modalClose = document.getElementById('modal-close');
  const modalCancel = document.getElementById('modal-cancel');

  const openModal = (title = 'Nueva evaluación') => {
    if (!modalBg) return;
    if (modalTitle) modalTitle.textContent = title;
    modalBg.style.display = 'grid';
  };

  const closeModal = () => {
    if (!modalBg) return;
    modalBg.style.display = 'none';
  };

  modalClose?.addEventListener('click', closeModal);
  modalCancel?.addEventListener('click', closeModal);

  return { openModal, closeModal };
}
