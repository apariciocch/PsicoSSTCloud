#!/usr/bin/env python3
"""
Script de inicialización de PocketBase para SBC Backend
Crea roles iniciales, usuario admin y datos de prueba
"""

import requests
import json
import time
import sys
from typing import Dict, Any

# Configuración
POCKETBASE_URL = "http://localhost:8090"
ADMIN_EMAIL = "admin@example.com"
ADMIN_PASSWORD = "secure_password"

# Colores para terminal
class Colors:
    HEADER = '\033[95m'
    BLUE = '\033[94m'
    CYAN = '\033[96m'
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    END = '\033[0m'
    BOLD = '\033[1m'

def print_status(message: str, status: str = "INFO"):
    """Imprimir mensaje de estado"""
    colors = {
        "INFO": Colors.BLUE,
        "SUCCESS": Colors.GREEN,
        "ERROR": Colors.RED,
        "WARNING": Colors.YELLOW,
    }
    color = colors.get(status, Colors.BLUE)
    print(f"{color}[{status}]{Colors.END} {message}")

def check_pocketbase_health() -> bool:
    """Verificar que PocketBase esté corriendo"""
    try:
        response = requests.get(f"{POCKETBASE_URL}/api/health", timeout=5)
        return response.status_code == 200
    except Exception as e:
        print_status(f"PocketBase no responde: {str(e)}", "ERROR")
        return False

def create_record(collection: str, data: Dict[str, Any], admin_token: str = None) -> Dict[str, Any]:
    """Crear un registro en una colección"""
    headers = {"Content-Type": "application/json"}
    if admin_token:
        headers["Authorization"] = f"Bearer {admin_token}"
    
    try:
        response = requests.post(
            f"{POCKETBASE_URL}/api/collections/{collection}/records",
            json=data,
            headers=headers,
            timeout=10
        )
        
        if response.status_code in [200, 201]:
            return response.json()
        else:
            print_status(
                f"Error creando {collection}: {response.status_code} - {response.text}",
                "ERROR"
            )
            return None
    except Exception as e:
        print_status(f"Excepción creando {collection}: {str(e)}", "ERROR")
        return None

def initialize_pocketbase():
    """Inicializar PocketBase con datos"""
    
    print(f"\n{Colors.HEADER}{Colors.BOLD}🚀 Inicializador de PocketBase - SBC Backend{Colors.END}\n")
    
    # 1. Verificar PocketBase
    print_status("Verificando conexión a PocketBase...", "INFO")
    if not check_pocketbase_health():
        print_status("PocketBase no está disponible. Inicia con: docker-compose up -d pocketbase", "ERROR")
        sys.exit(1)
    print_status("✅ PocketBase conectado", "SUCCESS")
    
    # 2. Crear roles
    print_status("\nCreando roles del sistema...", "INFO")
    
    roles_data = [
        {
            "name": "Administrador",
            "slug": "admin",
            "description": "Acceso total al sistema y gestión de usuarios",
            "permissions": ["*"],
            "level": 1,
            "is_system_role": True
        },
        {
            "name": "Supervisor",
            "slug": "supervisor",
            "description": "Gestión de equipo, observaciones y reportes",
            "permissions": [
                "users.view", "users.edit",
                "observations.view", "observations.create", "observations.edit",
                "incidents.view", "incidents.edit",
                "reports.view", "reports.create"
            ],
            "level": 2,
            "is_system_role": True
        },
        {
            "name": "Observador",
            "slug": "observer",
            "description": "Crear observaciones y ver incidentes",
            "permissions": [
                "observations.create", "observations.view",
                "incidents.view"
            ],
            "level": 3,
            "is_system_role": True
        },
        {
            "name": "Trabajador",
            "slug": "worker",
            "description": "Solo ver sus propias observaciones",
            "permissions": [
                "observations.view"
            ],
            "level": 4,
            "is_system_role": True
        },
        {
            "name": "Auditor",
            "slug": "auditor",
            "description": "Revisar auditoría y reportes",
            "permissions": [
                "audit_logs.view",
                "reports.view"
            ],
            "level": 5,
            "is_system_role": True
        }
    ]
    
    roles_created = {}
    for role_data in roles_data:
        result = create_record("roles", role_data)
        if result:
            role_slug = role_data["slug"]
            roles_created[role_slug] = result["id"]
            print_status(f"  ✅ Rol '{role_data['name']}' creado", "SUCCESS")
        else:
            print_status(f"  ❌ Error creando rol '{role_data['name']}'", "ERROR")
    
    if not roles_created:
        print_status("No se pudieron crear los roles. Abortando.", "ERROR")
        sys.exit(1)
    
    # 3. Crear usuario admin
    print_status("\nCreando usuario administrador...", "INFO")
    
    import hashlib
    import os
    from datetime import datetime
    
    # Usar bcrypt para hash de contraseña (simulado)
    # En producción, usar: python -m bcrypt password
    password_hash = hashlib.sha256(ADMIN_PASSWORD.encode()).hexdigest()
    
    admin_user_data = {
        "email": ADMIN_EMAIL,
        "username": "admin",
        "password_hash": password_hash,
        "role_id": roles_created["admin"],
        "is_active": True,
        "failed_attempts": 0
    }
    
    admin_user = create_record("users", admin_user_data)
    if admin_user:
        print_status(f"  ✅ Usuario admin creado: {ADMIN_EMAIL}", "SUCCESS")
    else:
        print_status(f"  ❌ Error creando usuario admin", "ERROR")
    
    # 4. Crear áreas de trabajo de ejemplo
    print_status("\nCreando áreas de trabajo...", "INFO")
    
    work_areas = [
        {
            "name": "Producción",
            "description": "Planta de producción principal",
            "location": "Edificio A - Piso 1",
            "risk_level": "high",
            "is_active": True
        },
        {
            "name": "Almacén",
            "description": "Almacén de materia prima",
            "location": "Edificio B",
            "risk_level": "medium",
            "is_active": True
        },
        {
            "name": "Oficinas",
            "description": "Área administrativa",
            "location": "Edificio A - Piso 2",
            "risk_level": "low",
            "is_active": True
        },
        {
            "name": "Mantenimiento",
            "description": "Taller de mantenimiento",
            "location": "Edificio C",
            "risk_level": "critical",
            "is_active": True
        }
    ]
    
    work_areas_created = {}
    for area_data in work_areas:
        result = create_record("work_areas", area_data)
        if result:
            work_areas_created[area_data["name"]] = result["id"]
            print_status(f"  ✅ Área '{area_data['name']}' creada", "SUCCESS")
        else:
            print_status(f"  ❌ Error creando área '{area_data['name']}'", "ERROR")
    
    # 5. Crear categorías de conductas
    print_status("\nCreando categorías de conductas...", "INFO")
    
    behavior_categories = [
        {
            "name": "Uso de EPE",
            "description": "Uso correcto de Equipos de Protección Personal",
            "category_type": "safe",
            "color_code": "#00AA00"
        },
        {
            "name": "No usar EPE",
            "description": "Trabajar sin equipos de protección cuando es requerido",
            "category_type": "unsafe",
            "color_code": "#CC0000"
        },
        {
            "name": "Orden y Limpieza",
            "description": "Mantener área de trabajo ordenada y limpia",
            "category_type": "safe",
            "color_code": "#00AA00"
        },
        {
            "name": "Área Desordenada",
            "description": "Área de trabajo con desorden o suciedad",
            "category_type": "unsafe",
            "color_code": "#CC0000"
        },
        {
            "name": "Comunicación Efectiva",
            "description": "Comunicarse claramente en actividades de riesgo",
            "category_type": "safe",
            "color_code": "#00AA00"
        },
        {
            "name": "Falta de Comunicación",
            "description": "No comunicar adecuadamente riesgos",
            "category_type": "unsafe",
            "color_code": "#CC0000"
        }
    ]
    
    categories_created = {}
    for cat_data in behavior_categories:
        result = create_record("behavior_categories", cat_data)
        if result:
            categories_created[cat_data["name"]] = result["id"]
            print_status(f"  ✅ Categoría '{cat_data['name']}' creada", "SUCCESS")
        else:
            print_status(f"  ❌ Error creando categoría '{cat_data['name']}'", "ERROR")
    
    # 6. Crear empleados de ejemplo
    print_status("\nCreando empleados de ejemplo...", "INFO")
    
    if work_areas_created:
        first_area_id = list(work_areas_created.values())[0]
        
        employees = [
            {
                "first_name": "Carlos",
                "last_name": "García",
                "email": "carlos.garcia@empresa.com",
                "document_number": "12345678",
                "phone": "+34 912 345 678",
                "work_area_id": first_area_id,
                "job_position": "Operario",
                "hire_date": "2023-01-15",
                "is_active": True
            },
            {
                "first_name": "María",
                "last_name": "López",
                "email": "maria.lopez@empresa.com",
                "document_number": "87654321",
                "phone": "+34 912 345 679",
                "work_area_id": first_area_id,
                "job_position": "Supervisora",
                "hire_date": "2022-06-10",
                "is_active": True
            }
        ]
        
        for emp_data in employees:
            result = create_record("employees", emp_data)
            if result:
                print_status(f"  ✅ Empleado '{emp_data['first_name']} {emp_data['last_name']}' creado", "SUCCESS")
            else:
                print_status(f"  ❌ Error creando empleado '{emp_data['first_name']}'", "ERROR")
    
    # Resumen final
    print(f"\n{Colors.GREEN}{Colors.BOLD}✅ Inicialización completada{Colors.END}\n")
    print(f"{Colors.CYAN}📋 Resumen:{Colors.END}")
    print(f"  • Roles creados: {len(roles_created)}/5")
    print(f"  • Áreas creadas: {len(work_areas_created)}/4")
    print(f"  • Categorías creadas: {len(categories_created)}/6")
    print(f"\n{Colors.CYAN}🔐 Credenciales de acceso:{Colors.END}")
    print(f"  • Email: {ADMIN_EMAIL}")
    print(f"  • Password: {ADMIN_PASSWORD}")
    print(f"  • URL Admin: {POCKETBASE_URL}/_/")
    print(f"\n{Colors.CYAN}🚀 Próximos pasos:{Colors.END}")
    print(f"  1. Iniciar backend: docker-compose up -d backend")
    print(f"  2. Acceder a PocketBase: {POCKETBASE_URL}/_/")
    print(f"  3. Probar API: curl http://localhost:8080/api/v1/auth/login")
    print()

if __name__ == "__main__":
    try:
        initialize_pocketbase()
    except KeyboardInterrupt:
        print(f"\n{Colors.YELLOW}⚠️ Inicialización cancelada por usuario{Colors.END}\n")
        sys.exit(0)
    except Exception as e:
        print_status(f"Error inesperado: {str(e)}", "ERROR")
        sys.exit(1)
