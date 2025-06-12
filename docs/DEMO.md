# 🎯 Demo - UI CAF Factura Móvil Gateway

Esta guía te muestra cómo usar la interfaz React para gestionar archivos CAF.

## 🚀 Inicio Rápido

### 1. Iniciar Backend
Primero asegúrate de que el backend Go esté corriendo:
```bash
# En la raíz del proyecto
go run cmd/main.go
```

### 2. Iniciar Frontend
Desde la raíz del proyecto:
```bash
./start-ui.sh
```

O manualmente:
```bash
cd ui
npm start
```

La aplicación estará disponible en: **http://localhost:3002**

## 📱 Flujo de Usuario Completo

### Paso 1: Pantalla Inicial
- Al abrir la aplicación verás una pantalla de bienvenida
- Instrucciones de uso en 3 pasos
- Header con branding chileno
- Footer con información del sistema SII

### Paso 2: Crear o Seleccionar Empresa

#### Crear Nueva Empresa
1. Clic en "Nueva Empresa" en el selector
2. Completar formulario:
   - **Nombre**: "Mi Empresa de Prueba S.A."
   - **RUT**: "12345678-9" (formato chileno)
   - **ID Factura Móvil**: 12345 (opcional)
3. Clic en "Crear Empresa"

#### Seleccionar Empresa Existente
1. Clic en el dropdown del selector
2. Elegir empresa de la lista
3. Ver información detallada

### Paso 3: Cargar Archivo CAF

#### Preparar Archivo de Prueba
Crea un archivo XML de prueba llamado `test-caf.xml`:

```xml
<?xml version="1.0" encoding="ISO-8859-1"?>
<AUTORIZACION version="1.0">
    <CAF version="1.0">
        <DA>
            <RE>12345678-9</RE>
            <RS>Mi Empresa de Prueba S.A.</RS>
            <TD>33</TD>
            <RNG>
                <D>1</D>
                <H>1000</H>
            </RNG>
            <FA>2024-01-01</FA>
            <RSAPK>
                <M>nYtqeqtSYkOWUbEh6YqjNjLRo6R7MQJz9A6YwbGFJ4HKLr5JCb2kC9PGhRfY6XdH8JW2s9A4E3V8Z1T5RhJf3kGfKlLqNz7Cy5Xh8ZP9R4WgS3AqEfB6V2YwJcKmEz7DtY1HrF9</M>
                <E>Aw==</E>
            </RSAPK>
            <IDK>12345</IDK>
        </DA>
        <FRMA algoritmo="SHA1withRSA">abcdef123456789</FRMA>
    </CAF>
    <RSASK>-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQC9i2p6q1JiQ5ZRsSHpiqM2MtGjpHsxAnP0DpjBsYUngcouzkl
...
-----END RSA PRIVATE KEY-----</RSASK>
    <RSAPUBK>-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC9i2p6q1JiQ5ZRsSHpiqM2Mt
...
-----END PUBLIC KEY-----</RSAPUBK>
</AUTORIZACION>
```

#### Cargar CAF
1. **Drag & Drop**: Arrastra el archivo XML al área de carga
2. **O Clic**: Clic en el área para seleccionar archivo
3. **Validación**: El sistema valida tipo y tamaño
4. **Confirmación**: Clic en "Cargar CAF"
5. **Resultado**: Mensaje de éxito con detalles del CAF

### Paso 4: Ver Lista de CAFs

Después de cargar, verás:
- **Estado del CAF**: Activo, Por vencer, Expirado
- **Información de Folios**: Rango inicial-final, folio actual
- **Fechas**: Autorización y expiración
- **Progreso**: Barra visual de uso de folios
- **Detalles**: Tipo de documento, código empresa

## 🎨 Características de la UI

### Diseño Responsive
- **Desktop**: Layout de 2 columnas
- **Tablet**: Layout adaptativo
- **Mobile**: Stack vertical

### Estados Visuales
- **Cargando**: Spinners y placeholders
- **Error**: Mensajes informativos en rojo
- **Éxito**: Confirmaciones en verde
- **Advertencia**: Notificaciones en amarillo

### Validaciones
- **RUT Chileno**: Formato 12345678-9
- **Archivos**: Solo XML, máximo 5MB
- **Formularios**: Campos requeridos marcados

### Indicadores de Estado CAF
- 🟢 **Activo**: CAF válido con folios disponibles
- 🟡 **Cuidado**: Más del 70% de folios usados
- 🟠 **Por vencer**: Menos de 30 días o más del 90% usado
- 🔴 **Expirado**: CAF vencido

## 🔧 Funcionalidades Avanzadas

### Selector de Empresa
- Dropdown con información completa
- Búsqueda y filtrado
- Creación inline de empresas
- Estado persistente

### Cargador de CAF
- Drag & drop avanzado
- Validación en tiempo real
- Progress indicators
- Manejo de errores detallado

### Lista de CAFs
- Información completa de cada CAF
- Estados visuales intuitivos
- Progreso de uso de folios
- Fechas formateadas en español

## 🚨 Casos de Error Comunes

### Backend No Disponible
```
⚠️ Backend no detectado en http://localhost:8080
```
**Solución**: Iniciar el servidor Go

### Archivo Inválido
```
❌ Por favor selecciona un archivo XML válido
```
**Solución**: Usar archivo .xml del SII

### RUT Inválido
```
❌ El RUT debe tener el formato correcto (ej: 12345678-9)
```
**Solución**: Usar formato chileno con guión

### Empresa Duplicada
```
❌ Error al crear la empresa. Verifica que el RUT no esté ya registrado.
```
**Solución**: Usar RUT único o seleccionar empresa existente

## 📊 Datos de Prueba

### Empresas de Ejemplo
```json
{
  "name": "Empresa Demo S.A.",
  "code": "11111111-1", 
  "factura_movil_company_id": 12345
}
```

### CAF de Prueba
- **Tipo**: 33 (Factura Electrónica)
- **Folios**: 1 - 1000
- **RUT**: Debe coincidir con empresa seleccionada

## 🎉 Resultado Final

Una vez completado el flujo, tendrás:
1. ✅ Empresa creada/seleccionada
2. ✅ CAF cargado y validado
3. ✅ Dashboard con información completa
4. ✅ Sistema listo para generar documentos

## 📱 Screenshots y Videos

### Pantalla Principal
- Layout limpio y profesional
- Branding chileno con colores de la bandera
- Información contextual del SII

### Formularios
- Validación en tiempo real
- Mensajes de error informativos
- Confirmaciones visuales

### Dashboard de CAFs
- Estados color-coded
- Información completa y organizada
- Responsive en todos los dispositivos

## 🔄 Refresh y Sincronización

La UI se sincroniza automáticamente:
- **Crear empresa** → Selecciona automáticamente
- **Cargar CAF** → Actualiza lista inmediatamente  
- **Cambiar empresa** → Carga CAFs correspondientes

## 🌟 Próximos Pasos

Después del demo, puedes:
1. Explorar la generación de timbres
2. Integrar con sistemas de facturación
3. Configurar alertas de vencimiento
4. Añadir más tipos de documentos 