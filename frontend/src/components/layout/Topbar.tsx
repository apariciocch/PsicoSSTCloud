import { useLocation } from 'react-router-dom';

const breadcrumbMap: Record<string, string> = {
  '/dashboard': 'Dashboard',
  '/evaluaciones': 'Evaluaciones',
  '/reportes': 'Reportes',
  '/carpeta-legal': 'Carpeta Legal',
  '/plan-accion': 'Plan de Acción',
  '/capacitaciones': 'Capacitaciones',
  '/casos': 'Casos',
  '/usuarios': 'Usuarios',
  '/configuracion': 'Configuración',
};

export default function Topbar() {
  const location = useLocation();
  const breadcrumb = breadcrumbMap[location.pathname] || '';
  return (
    <header className="h-16 bg-white border-b flex items-center px-6 justify-between">
      <div className="text-gray-500 text-sm">{breadcrumb}</div>
      <div>
        {/* Aquí puedes agregar botones de acción globales */}
      </div>
    </header>
  );
}
