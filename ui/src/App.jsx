import React, { useState } from 'react';
import { FileText, Flag } from 'lucide-react';
import CompanySelector from './components/CompanySelector';
import CAFUploader from './components/CAFUploader';
import CAFList from './components/CAFList';
import CreateCompanyModal from './components/CreateCompanyModal';
import CommercialActivitiesManager from './components/CommercialActivitiesManager';

function App() {
  const [selectedCompany, setSelectedCompany] = useState(null);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleCompanySelect = (company) => {
    setSelectedCompany(company);
  };

  const handleCreateCompany = () => {
    setIsCreateModalOpen(true);
  };

  const handleCompanyCreated = (newCompany) => {
    setSelectedCompany(newCompany);
    // Trigger refresh of any dependent components
    setRefreshTrigger(prev => prev + 1);
  };

  const handleUploadSuccess = () => {
    // Trigger refresh of CAF list
    setRefreshTrigger(prev => prev + 1);
  };

  const handleCAFUploadSuccess = () => {
    // Trigger refresh of CAF list
    setRefreshTrigger(prev => prev + 1);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-3">
                <div className="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center">
                  <FileText className="w-5 h-5 text-white" />
                </div>
                <div>
                  <h1 className="text-xl font-bold text-gray-900">
                    Factura Móvil Gateway
                  </h1>
                  <p className="text-sm text-gray-500">
                    Sistema de gestión de CAF
                  </p>
                </div>
              </div>
            </div>
            
            <div className="flex items-center space-x-3">
              <div className="flex items-center space-x-2 text-sm text-gray-600">
                <Flag className="w-4 h-4 text-red-500" />
                <span>Chile</span>
              </div>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">
            Gestión de Archivos CAF
          </h2>
          <p className="text-gray-600">
            Administra los Códigos de Autorización de Folios (CAF) para tu empresa.
            Los archivos CAF son necesarios para emitir documentos tributarios electrónicos.
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Left Column */}
          <div className="space-y-8">
            {/* Company Selector */}
            <CompanySelector
              selectedCompany={selectedCompany}
              onCompanySelect={handleCompanySelect}
              onCreateCompany={handleCreateCompany}
            />

            {/* Company Details */}
            {selectedCompany && (
                <div className="space-y-6">
                    <div className="bg-white shadow rounded-lg p-6">
                        <h2 className="text-lg font-semibold text-gray-900 mb-4">Detalles de la Empresa</h2>
                        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                            <div>
                                <p className="text-sm font-medium text-gray-500">Nombre</p>
                                <p className="mt-1 text-sm text-gray-900">{selectedCompany.name}</p>
                            </div>
                            <div>
                                <p className="text-sm font-medium text-gray-500">RUT</p>
                                <p className="mt-1 text-sm text-gray-900">{selectedCompany.code}</p>
                            </div>
                            {selectedCompany.factura_movil_company_id && (
                                <div>
                                    <p className="text-sm font-medium text-gray-500">ID Factura Móvil</p>
                                    <p className="mt-1 text-sm text-gray-900">{selectedCompany.factura_movil_company_id}</p>
                                </div>
                            )}
                        </div>
                    </div>

                    {/* Commercial Activities Manager */}
                    <CommercialActivitiesManager
                        selectedCompany={selectedCompany}
                        refreshTrigger={refreshTrigger}
                    />

                    {/* CAF Uploader */}
                    <CAFUploader
                        selectedCompany={selectedCompany}
                        onUploadSuccess={handleCAFUploadSuccess}
                    />
                </div>
            )}
          </div>

          {/* Right Column */}
          <div className="space-y-8">
            {/* CAF List */}
            <CAFList
              selectedCompany={selectedCompany}
              refreshTrigger={refreshTrigger}
            />
          </div>
        </div>

        {/* Instructions Section */}
        {!selectedCompany && (
          <div className="mt-12">
            <div className="card bg-gradient-to-br from-primary-50 to-blue-50 border-primary-200">
              <div className="text-center">
                <div className="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
                  <FileText className="w-8 h-8 text-primary-600" />
                </div>
                <h3 className="text-xl font-semibold text-gray-900 mb-4">
                  ¿Cómo usar este sistema?
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6 text-left">
                  <div className="space-y-3">
                    <div className="w-8 h-8 bg-primary-600 text-white rounded-full flex items-center justify-center text-sm font-bold">
                      1
                    </div>
                    <h4 className="font-semibold text-gray-900">Selecciona o Crea una Empresa</h4>
                    <p className="text-sm text-gray-600">
                      Elige una empresa existente de la lista o crea una nueva con su RUT correspondiente.
                    </p>
                  </div>
                  <div className="space-y-3">
                    <div className="w-8 h-8 bg-primary-600 text-white rounded-full flex items-center justify-center text-sm font-bold">
                      2
                    </div>
                    <h4 className="font-semibold text-gray-900">Carga el Archivo CAF</h4>
                    <p className="text-sm text-gray-600">
                      Arrastra y suelta o selecciona el archivo XML del CAF proporcionado por el SII.
                    </p>
                  </div>
                  <div className="space-y-3">
                    <div className="w-8 h-8 bg-primary-600 text-white rounded-full flex items-center justify-center text-sm font-bold">
                      3
                    </div>
                    <h4 className="font-semibold text-gray-900">Administra tus CAFs</h4>
                    <p className="text-sm text-gray-600">
                      Visualiza el estado, folios disponibles y fechas de vencimiento de todos tus CAFs.
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-gray-200 mt-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <div className="chile-accent w-6 h-4 rounded"></div>
              <div>
                <p className="text-sm text-gray-600">
                  Sistema de Facturación Electrónica - Chile
                </p>
                <p className="text-xs text-gray-500">
                  Compatible con normativas del SII (Servicio de Impuestos Internos)
                </p>
              </div>
            </div>
            <div className="text-xs text-gray-500">
              v1.0.0
            </div>
          </div>
        </div>
      </footer>

      {/* Create Company Modal */}
      <CreateCompanyModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onCompanyCreated={handleCompanyCreated}
      />
    </div>
  );
}

export default App; 