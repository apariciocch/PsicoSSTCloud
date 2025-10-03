import { NavLink } from 'react-router-dom';

const menu = [
  { to: '/dashboard', label: 'Dashboard' },
  { to: '/evaluaciones', label: 'Evaluaciones' },
  { to: '/reportes', label: 'Reportes' },
  { to: '/carpeta-legal', label: 'Carpeta Legal' },
  { to: '/plan-accion', label: 'Plan de Acción' },
  { to: '/capacitaciones', label: 'Capacitaciones' },
  { to: '/casos', label: 'Casos' },
  { to: '/usuarios', label: 'Usuarios' },
  { to: '/configuracion', label: 'Configuración' },
];

export default function Sidebar() {
  return (
    <aside className="w-64 bg-white shadow h-full flex flex-col">
      <div className="h-16 flex items-center justify-center font-bold text-xl border-b">PsicoSST Cloud</div>
      <nav className="flex-1 p-4">
        <ul className="space-y-2">
          {menu.map((item) => (
            <li key={item.to}>
              <NavLink
                to={item.to}
                className={({ isActive }) =>
                  `block px-4 py-2 rounded hover:bg-blue-100 transition ${isActive ? 'bg-blue-500 text-white' : 'text-gray-700'}`
                }
              >
                {item.label}
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
}
