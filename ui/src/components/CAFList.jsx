import React, { useState, useEffect } from 'react';
import { FileText, Calendar, Hash, Building2, AlertTriangle, CheckCircle, Clock } from 'lucide-react';
import { cafService } from '../services/api';

const CAFList = ({ selectedCompany, refreshTrigger }) => {
  const [cafs, setCafs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (selectedCompany) {
      loadCAFs();
    }
  }, [selectedCompany, refreshTrigger]);

  const loadCAFs = async () => {
    if (!selectedCompany) return;

    try {
      setLoading(true);
      setError(null);
      const cafsData = await cafService.getCompanyCAFs(selectedCompany.id);
      setCafs(cafsData);
    } catch (err) {
      setError('Error al cargar los CAFs');
      console.error('Error loading CAFs:', err);
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('es-CL', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const getDocumentTypeName = (type) => {
    const types = {
      33: 'Factura Electrónica',
      34: 'Factura Exenta',
      39: 'Boleta Electrónica',
      41: 'Boleta Exenta',
      52: 'Guía de Despacho',
      56: 'Nota de Débito',
      61: 'Nota de Crédito'
    };
    return types[type] || `Tipo ${type}`;
  };

  const getStatusColor = (expirationDate, currentFolios, finalFolios) => {
    const now = new Date();
    const expiry = new Date(expirationDate);
    const daysToExpiry = (expiry - now) / (1000 * 60 * 60 * 24);
    const foliosUsed = (currentFolios - 1); // Assuming initial is 1
    const totalFolios = finalFolios;
    const usagePercentage = (foliosUsed / totalFolios) * 100;

    if (daysToExpiry < 0) return 'expired';
    if (daysToExpiry < 30 || usagePercentage > 90) return 'warning';
    if (usagePercentage > 70) return 'caution';
    return 'active';
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'expired':
        return <AlertTriangle className="w-4 h-4 text-error-500" />;
      case 'warning':
        return <AlertTriangle className="w-4 h-4 text-orange-500" />;
      case 'caution':
        return <Clock className="w-4 h-4 text-yellow-500" />;
      default:
        return <CheckCircle className="w-4 h-4 text-success-500" />;
    }
  };

  const getStatusText = (status) => {
    switch (status) {
      case 'expired':
        return 'Expirado';
      case 'warning':
        return 'Por vencer';
      case 'caution':
        return 'Cuidado';
      default:
        return 'Activo';
    }
  };

  if (!selectedCompany) {
    return (
      <div className="card border-gray-200 bg-gray-50">
        <div className="text-center py-8">
          <FileText className="w-12 h-12 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500">Selecciona una empresa para ver sus CAFs</p>
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="card">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-semibold text-gray-900 flex items-center">
            <FileText className="w-5 h-5 mr-2 text-primary-600" />
            CAFs de la Empresa
          </h2>
        </div>
        <div className="space-y-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="border border-gray-200 rounded-lg p-4">
              <div className="animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-1/3 mb-2"></div>
                <div className="h-3 bg-gray-200 rounded w-1/2"></div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="card border-error-200 bg-error-50">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <AlertTriangle className="w-5 h-5 text-error-500" />
            <p className="text-error-700">{error}</p>
          </div>
          <button 
            onClick={loadCAFs}
            className="btn-secondary text-sm"
          >
            Reintentar
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="card">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-lg font-semibold text-gray-900 flex items-center">
          <FileText className="w-5 h-5 mr-2 text-primary-600" />
          CAFs de la Empresa
        </h2>
        <div className="text-sm text-gray-500">
          {cafs.length} CAF{cafs.length !== 1 ? 's' : ''} encontrado{cafs.length !== 1 ? 's' : ''}
        </div>
      </div>

      {cafs.length === 0 ? (
        <div className="text-center py-12">
          <FileText className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">No hay CAFs cargados</h3>
          <p className="text-gray-500 mb-6">
            Esta empresa aún no tiene archivos CAF cargados. 
            Carga tu primer CAF para comenzar a generar documentos tributarios.
          </p>
        </div>
      ) : (
        <div className="space-y-4">
          {cafs.map((caf) => {
            const status = getStatusColor(caf.expirationDate, caf.currentFolios, caf.finalFolios);
            const foliosUsed = caf.currentFolios - caf.initialFolios;
            const totalFolios = caf.finalFolios - caf.initialFolios + 1;
            const usagePercentage = (foliosUsed / totalFolios) * 100;

            return (
              <div 
                key={caf.id} 
                className={`border rounded-lg p-4 ${
                  status === 'expired' ? 'border-error-200 bg-error-50' :
                  status === 'warning' ? 'border-orange-200 bg-orange-50' :
                  status === 'caution' ? 'border-yellow-200 bg-yellow-50' :
                  'border-gray-200 bg-white'
                }`}
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center space-x-3 mb-3">
                      <div className="flex items-center space-x-2">
                        {getStatusIcon(status)}
                        <span className={`text-sm font-medium ${
                          status === 'expired' ? 'text-error-700' :
                          status === 'warning' ? 'text-orange-700' :
                          status === 'caution' ? 'text-yellow-700' :
                          'text-success-700'
                        }`}>
                          {getStatusText(status)}
                        </span>
                      </div>
                      <div className="flex items-center space-x-2 text-sm text-gray-500">
                        <Building2 className="w-4 h-4" />
                        <span>{getDocumentTypeName(caf.documentType)}</span>
                      </div>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                      <div>
                        <div className="flex items-center space-x-2 text-sm text-gray-600 mb-1">
                          <Hash className="w-4 h-4" />
                          <span className="font-medium">Folios</span>
                        </div>
                        <p className="text-gray-900">
                          {caf.initialFolios.toLocaleString()} - {caf.finalFolios.toLocaleString()}
                        </p>
                        <p className="text-xs text-gray-500">
                          Actual: {caf.currentFolios.toLocaleString()}
                        </p>
                      </div>

                      <div>
                        <div className="flex items-center space-x-2 text-sm text-gray-600 mb-1">
                          <Calendar className="w-4 h-4" />
                          <span className="font-medium">Autorización</span>
                        </div>
                        <p className="text-gray-900">{formatDate(caf.authorizationDate)}</p>
                        <p className="text-xs text-gray-500">
                          Expira: {formatDate(caf.expirationDate)}
                        </p>
                      </div>

                      <div>
                        <div className="text-sm text-gray-600 mb-1 font-medium">Uso de Folios</div>
                        <div className="flex items-center space-x-3">
                          <div className="flex-1 bg-gray-200 rounded-full h-2">
                            <div 
                              className={`h-2 rounded-full ${
                                usagePercentage > 90 ? 'bg-error-500' :
                                usagePercentage > 70 ? 'bg-orange-500' :
                                'bg-success-500'
                              }`}
                              style={{ width: `${Math.min(usagePercentage, 100)}%` }}
                            ></div>
                          </div>
                          <span className="text-xs text-gray-600 min-w-[3rem]">
                            {usagePercentage.toFixed(1)}%
                          </span>
                        </div>
                        <p className="text-xs text-gray-500 mt-1">
                          {foliosUsed.toLocaleString()} de {totalFolios.toLocaleString()} folios usados
                        </p>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Additional info */}
                <div className="mt-4 pt-3 border-t border-gray-200">
                  <div className="flex items-center justify-between text-xs text-gray-500">
                    <span>ID: {caf.id}</span>
                    <span>Código empresa: {caf.companyCode}</span>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};

export default CAFList; 