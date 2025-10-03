const statusClassMap = {
  Activo: 'pill ok',
  Invitado: 'pill',
  Suspendido: 'pill warn'
};

const usuarios = [
  { name: 'Armando Aparicio', email: 'armando@empresa.pe', role: 'Psicólogo', status: 'Activo' },
  { name: 'María Torres', email: 'maria@empresa.pe', role: 'Comité SST', status: 'Invitado' }
];

export function initUsuarios() {
  const userTableBody = document.getElementById('user-table-body');
  const userModalBg = document.getElementById('user-modal-bg');
  const userModalClose = document.getElementById('user-modal-close');
  const userModalCancel = document.getElementById('user-modal-cancel');
  const userModalTitle = document.getElementById('user-modal-title');
  const userForm = document.getElementById('user-form');
  const userNameInput = document.getElementById('user-name');
  const userEmailInput = document.getElementById('user-email');
  const userRoleInput = document.getElementById('user-role');
  const userStatusInput = document.getElementById('user-status');
  const btnNuevoUser = document.getElementById('btnNuevoUser');

  if (!userTableBody || !userForm || !userModalBg || !userModalTitle || !userNameInput || !userEmailInput || !userRoleInput || !userStatusInput) {
    return;
  }

  let editingUserIndex = null;

  const renderUsers = () => {
    if (!usuarios.length) {
      userTableBody.innerHTML = '<tr><td colspan="5" class="muted" style="text-align:center">No hay usuarios registrados.</td></tr>';
      return;
    }

    userTableBody.innerHTML = usuarios.map((user, index) => {
      const pillClass = statusClassMap[user.status] || 'pill';
      return `
        <tr>
          <td>${user.name}</td>
          <td>${user.email}</td>
          <td>${user.role}</td>
          <td><span class="${pillClass}">${user.status}</span></td>
          <td><button class="btn btn-edit-user" data-index="${index}">Editar</button></td>
        </tr>
      `;
    }).join('');
  };

  const closeUserModal = () => {
    userModalBg.style.display = 'none';
    userForm.reset();
    userStatusInput.value = 'Activo';
    userRoleInput.value = 'Psicólogo';
    editingUserIndex = null;
  };

  const openUserModal = (index = null) => {
    editingUserIndex = index;
    if (index === null) {
      userForm.reset();
      userStatusInput.value = 'Activo';
      userRoleInput.value = 'Psicólogo';
      userModalTitle.textContent = 'Nuevo usuario';
    } else {
      const user = usuarios[index];
      userNameInput.value = user.name;
      userEmailInput.value = user.email;
      userRoleInput.value = user.role;
      userStatusInput.value = user.status;
      userModalTitle.textContent = 'Editar usuario';
    }
    userModalBg.style.display = 'grid';
    setTimeout(() => userNameInput.focus(), 50);
  };

  btnNuevoUser?.addEventListener('click', () => openUserModal());
  userModalClose?.addEventListener('click', closeUserModal);
  userModalCancel?.addEventListener('click', closeUserModal);

  userForm.addEventListener('submit', (event) => {
    event.preventDefault();
    const newUser = {
      name: userNameInput.value.trim(),
      email: userEmailInput.value.trim(),
      role: userRoleInput.value,
      status: userStatusInput.value
    };

    if (!newUser.name || !newUser.email) {
      return;
    }

    if (editingUserIndex === null) {
      usuarios.push(newUser);
    } else {
      usuarios[editingUserIndex] = { ...usuarios[editingUserIndex], ...newUser };
    }

    renderUsers();
    closeUserModal();
  });

  userTableBody.addEventListener('click', (event) => {
    const button = event.target.closest('.btn-edit-user');
    if (!button) return;
    const index = Number(button.dataset.index);
    if (Number.isNaN(index)) return;
    openUserModal(index);
  });

  renderUsers();
}
