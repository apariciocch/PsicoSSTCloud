import 'dotenv/config';
import express from 'express';
import cors from 'cors';
import path from 'path';
import { fileURLToPath } from 'url';
import { query } from './db.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const app = express();
app.use(cors());
app.use(express.json());

// Servir el front (archivos estáticos desde la raíz del proyecto)
app.use(express.static(path.join(__dirname, '..')));

/* ------------------------------
   ENDPOINTS PRINCIPALES
------------------------------ */

// KPIs Dashboard
app.get('/api/kpis', async (_req, res) => {
  try {
    const { rows } = await query('SELECT * FROM app.v_dashboard_kpis ORDER BY tenant_nombre LIMIT 50;');
    res.json(rows);
  } catch (e) {
    console.error(e);
    res.status(500).json({ error: 'Error obteniendo KPIs' });
  }
});

// Evaluaciones - resumen
app.get('/api/evaluaciones', async (_req, res) => {
  try {
    const { rows } = await query(`
      SELECT id, tenant_id, area, periodo, estado, total_trab, instrumento, respondieron
      FROM app.v_evaluaciones_resumen
      ORDER BY periodo DESC, area ASC;
    `);
    res.json(rows);
  } catch (e) {
    console.error(e);
    res.status(500).json({ error: 'Error listando evaluaciones' });
  }
});

// Crear evaluación
app.post('/api/evaluaciones', async (req, res) => {
  try {
    const { tenant_id, instrumento, area, periodo, total_trab } = req.body;

    // Buscar instrumento_id por nombre
    const ins = await query(
      `SELECT id FROM app.instrumentos WHERE tenant_id = $1 AND nombre = $2 LIMIT 1;`,
      [tenant_id, instrumento]
    );
    if (ins.rowCount === 0) return res.status(400).json({ error: 'Instrumento no encontrado' });

    const instrumento_id = ins.rows[0].id;

    const { rows } = await query(
      `INSERT INTO app.evaluaciones (tenant_id, instrumento_id, area, periodo, total_trab, estado)
       VALUES ($1, $2, $3, $4, $5, 'En curso')
       RETURNING id;`,
      [tenant_id, instrumento_id, area, periodo + '-01', total_trab || 0]
    );
    res.status(201).json({ id: rows[0].id });
  } catch (e) {
    console.error(e);
    res.status(500).json({ error: 'Error creando evaluación' });
  }
});

// Finalizar evaluación
app.put('/api/evaluaciones/:id/finalizar', async (req, res) => {
  try {
    await query(`UPDATE app.evaluaciones SET estado = 'Completado' WHERE id = $1;`, [req.params.id]);
    res.json({ ok: true });
  } catch (e) {
    console.error(e);
    res.status(500).json({ error: 'Error actualizando evaluación' });
  }
});

// Eliminar evaluación
app.delete('/api/evaluaciones/:id', async (req, res) => {
  try {
    await query(`DELETE FROM app.evaluaciones WHERE id = $1;`, [req.params.id]);
    res.json({ ok: true });
  } catch (e) {
    console.error(e);
    res.status(500).json({ error: 'Error eliminando evaluación' });
  }
});

/* --------- Acciones (Plan) mínimas ---------- */

// Listar acciones pendientes
app.get('/api/acciones', async (_req, res) => {
  try {
    const { rows } = await query(`SELECT * FROM app.v_acciones_pendientes ORDER BY vence NULLS LAST;`);
    res.json(rows);
  } catch (e) {
    console.error(e);
    res.status(500).json({ error: 'Error listando acciones' });
  }
});

// Crear acción
app.post('/api/acciones', async (req, res) => {
  try {
    const { tenant_id, evaluacion_id, riesgo, accion, responsable, vence } = req.body;
    const { rows } = await query(
      `INSERT INTO app.acciones (tenant_id, evaluacion_id, riesgo, accion, responsable, vence)
       VALUES ($1,$2,$3,$4,$5,$6) RETURNING id;`,
      [tenant_id, evaluacion_id || null, riesgo, accion, responsable || null, vence || null]
    );
    res.status(201).json({ id: rows[0].id });
  } catch (e) {
    console.error(e);
    res.status(500).json({ error: 'Error creando acción' });
  }
});

// Actualizar estado de acción
app.put('/api/acciones/:id', async (req, res) => {
  try {
    const { estado } = req.body; // 'Pendiente' | 'En curso' | 'Completada'
    await query(`UPDATE app.acciones SET estado = $1 WHERE id = $2;`, [estado, req.params.id]);
    res.json({ ok: true });
  } catch (e) {
    console.error(e);
    res.status(500).json({ error: 'Error actualizando acción' });
  }
});

/* ------------------------------------------ */

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`✅ API corriendo en http://localhost:${PORT}`);
});
