# Factura Móvil Gateway - UI CAF

UI en React para la gestión de archivos CAF (Código de Autorización de Folios) del sistema de facturación electrónica chileno.

## 🇨🇱 Características

- **Gestión de Empresas**: Crear y seleccionar empresas para gestionar sus CAFs
- **Carga de CAF**: Interfaz drag & drop para cargar archivos XML del SII
- **Visualización de Estado**: Dashboard con información detallada de los CAFs
- **Validación RUT**: Validación automática del formato de RUT chileno
- **Responsive**: Diseño adaptable para diferentes dispositivos
- **Indicadores Visuales**: Estados de CAF (activo, por vencer, expirado)

## 🇨🇱 Inicio Rápido

### Requisitos

- Node.js 16+ y npm
- Backend API corriendo en `http://localhost:3000`

### Instalación

```bash
# Instalar dependencias
npm install

# Iniciar en modo desarrollo
npm start
```

La aplicación estará disponible en `http://localhost:3002`

### Scripts Disponibles

```bash
npm start    # Servidor de desarrollo
npm build    # Build para producción
npm test     # Ejecutar tests
```

## 🎨 Tecnologías

- **React 18** - Framework frontend
- **Tailwind CSS** - Estilos y diseño
- **Axios** - Cliente HTTP
- **Lucide React** - Iconos
- **PostCSS** - Procesamiento CSS

## 📱 Características de la UI

### Selector de Empresas
- Dropdown con búsqueda
- Información detallada de cada empresa
- Opción para crear nuevas empresas
- Validación de RUT chileno

### Cargador de CAF
- Drag & drop de archivos XML
- Validación de tipo y tamaño de archivo
- Feedback visual del proceso de carga
- Mensajes de error informativos

### Lista de CAFs
- Información completa de cada CAF
- Estados visuales (activo, advertencia, expirado)
- Progreso de uso de folios
- Fechas de autorización y vencimiento

### Modal de Creación de Empresa
- Formulario con validación
- Formato RUT chileno automático
- Campo opcional para ID Factura Móvil
- Información contextual

## 🔗 Integración con API

La aplicación se conecta con el backend Go a través de:

```
GET    /companies              # Listar empresas
POST   /companies              # Crear empresa
GET    /companies/{id}/cafs    # Listar CAFs de empresa
POST   /companies/{id}/cafs    # Cargar CAF
```

## 🎯 Flujo de Usuario

1. **Seleccionar/Crear Empresa**: El usuario elige una empresa existente o crea una nueva
2. **Cargar CAF**: Arrastra y suelta el archivo XML del CAF del SII
3. **Verificar Carga**: El sistema valida y procesa el archivo
4. **Gestionar CAFs**: Visualiza el estado y información de todos los CAFs

## 🔧 Configuración

### Variables de Entorno

Puedes configurar la URL del backend usando variables de entorno:

```bash
# .env.local
REACT_APP_API_URL=http://localhost:3000
```

### Configuración del Proxy

El proyecto está configurado para hacer proxy de las peticiones API al backend:

```json
"proxy": "http://localhost:3000"
```

## 🌟 Características Técnicas

### Manejo de Estado
- React Hooks para estado local
- Context API para estado global cuando sea necesario
- Refresh triggers para sincronización de datos

### Validaciones
- RUT chileno con formato `12345678-9`
- Archivos XML únicamente
- Tamaño máximo de archivo (5MB)
- Campos requeridos en formularios

### UX/UI
- Diseño responsive con Tailwind CSS
- Indicadores de carga y estados
- Mensajes de error informativos
- Confirmaciones visuales de acciones

### Optimizaciones
- Componentes reutilizables
- Lazy loading cuando corresponda
- Manejo eficiente de re-renders
- Axios interceptors para logging

## 🚀 Despliegue

### Build de Producción

```bash
npm run build
```

### Docker (Opcional)

```dockerfile
FROM node:16-alpine as builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 📚 Estructura del Proyecto

```
ui/
├── src/
│   ├── components/          # Componentes React
│   │   ├── CompanySelector.jsx
│   │   ├── CAFUploader.jsx
│   │   ├── CAFList.jsx
│   │   └── CreateCompanyModal.jsx
│   ├── services/           # Servicios API
│   │   └── api.js
│   ├── App.jsx            # Componente principal
│   ├── index.js           # Punto de entrada
│   └── index.css          # Estilos globales
├── public/                # Archivos públicos
└── package.json          # Dependencias
```

## 🤝 Contribución

1. Fork del repositorio
2. Crear branch para feature (`git checkout -b feature/AmazingFeature`)
3. Commit de cambios (`git commit -m 'Add AmazingFeature'`)
4. Push al branch (`git push origin feature/AmazingFeature`)
5. Abrir Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver `LICENSE` para más detalles. 