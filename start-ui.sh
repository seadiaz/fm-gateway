#!/bin/bash

echo "🚀 Iniciando Factura Móvil Gateway - UI CAF"
echo "========================================"

# Verificar si Node.js está instalado
if ! command -v node &> /dev/null; then
    echo "❌ Node.js no está instalado. Por favor instala Node.js 16+ desde https://nodejs.org/"
    exit 1
fi

# Verificar versión de Node.js
NODE_VERSION=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
if [ "$NODE_VERSION" -lt 16 ]; then
    echo "❌ Se requiere Node.js 16+. Versión actual: $(node -v)"
    exit 1
fi

echo "✅ Node.js $(node -v) detectado"

# Cambiar al directorio UI
cd ui

# Verificar si package.json existe
if [ ! -f "package.json" ]; then
    echo "❌ No se encontró package.json en el directorio ui/"
    exit 1
fi

# Instalar dependencias si no existen
if [ ! -d "node_modules" ]; then
    echo "📦 Instalando dependencias..."
    npm install
    if [ $? -ne 0 ]; then
        echo "❌ Error al instalar dependencias"
        exit 1
    fi
else
    echo "✅ Dependencias ya instaladas"
fi

# Verificar si el backend está corriendo
echo "🔍 Verificando conexión con el backend..."
if curl -f -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Backend detectado en http://localhost:8080"
else
    echo "⚠️  Backend no detectado en http://localhost:8080"
    echo "   Asegúrate de que el servidor Go esté corriendo"
    echo "   Puedes continuar, pero la UI no funcionará completamente sin el backend"
    read -p "¿Continuar de todas formas? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo ""
echo "🎉 Iniciando servidor de desarrollo..."
echo "   La aplicación estará disponible en: http://localhost:3000"
echo "   Presiona Ctrl+C para detener el servidor"
echo ""

# Iniciar el servidor de desarrollo
npm start 