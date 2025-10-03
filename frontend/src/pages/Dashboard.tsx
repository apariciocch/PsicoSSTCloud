export default function Dashboard() {
  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        {/* Tarjetas de indicadores */}
        <div className="bg-white p-4 rounded shadow flex flex-col items-center">
          <span className="text-lg font-semibold">Evaluaciones</span>
          <span className="text-3xl font-bold text-blue-600">120</span>
        </div>
        <div className="bg-white p-4 rounded shadow flex flex-col items-center">
          <span className="text-lg font-semibold">Empresas</span>
          <span className="text-3xl font-bold text-green-600">8</span>
        </div>
        <div className="bg-white p-4 rounded shadow flex flex-col items-center">
          <span className="text-lg font-semibold">Capacitaciones</span>
          <span className="text-3xl font-bold text-yellow-600">15</span>
        </div>
        <div className="bg-white p-4 rounded shadow flex flex-col items-center">
          <span className="text-lg font-semibold">Casos</span>
          <span className="text-3xl font-bold text-red-600">3</span>
        </div>
      </div>
      {/* Aquí irán los gráficos e indicadores */}
      <div className="bg-white p-6 rounded shadow">
        <p className="text-gray-500">Gráficos e indicadores próximamente...</p>
      </div>
    </div>
  );
}
