import { Routes, Route, Navigate } from 'react-router-dom';
import Sidebar from './components/layout/Sidebar';
import Topbar from './components/layout/Topbar';
import Dashboard from './pages/Dashboard';
import Evaluaciones from './pages/Evaluaciones';
import Reportes from './pages/Reportes';
import CarpetaLegal from './pages/CarpetaLegal';
import PlanAccion from './pages/PlanAccion';
import Capacitaciones from './pages/Capacitaciones';
import Casos from './pages/Casos';
import Usuarios from './pages/Usuarios';
import Configuracion from './pages/Configuracion';

function App() {
  return (
    <div className="flex h-screen bg-gray-100">
      <Sidebar />
      <div className="flex-1 flex flex-col">
        <Topbar />
        <main className="flex-1 overflow-y-auto p-4">
          <Routes>
            <Route path="/" element={<Navigate to="/dashboard" />} />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/evaluaciones" element={<Evaluaciones />} />
            <Route path="/reportes" element={<Reportes />} />
            <Route path="/carpeta-legal" element={<CarpetaLegal />} />
            <Route path="/plan-accion" element={<PlanAccion />} />
            <Route path="/capacitaciones" element={<Capacitaciones />} />
            <Route path="/casos" element={<Casos />} />
            <Route path="/usuarios" element={<Usuarios />} />
            <Route path="/configuracion" element={<Configuracion />} />
          </Routes>
        </main>
      </div>
    </div>
  );
}

export default App;
