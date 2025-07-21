import React, { useState, useEffect } from 'react';
import { Building2, ChevronDown, Plus } from 'lucide-react';
import { companyService } from '../services/api';

const CompanySelector = ({ selectedCompany, onCompanySelect, onCreateCompany }) => {
  const [companies, setCompanies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);

  useEffect(() => {
    loadCompanies();
  }, []);

  const loadCompanies = async () => {
    try {
      setLoading(true);
      const companiesData = await companyService.getCompanies();
      // Ensure we always have an array
      if (Array.isArray(companiesData)) {
        setCompanies(companiesData);
      } else {
        console.warn('Companies data is not an array:', companiesData);
        setCompanies([]);
      }
    } catch (err) {
      setError('Error al cargar empresas');
      console.error('Error loading companies:', err);
      setCompanies([]); // Ensure we have an empty array on error
    } finally {
      setLoading(false);
    }
  };

  const handleCompanySelect = (company) => {
    onCompanySelect(company);
    setIsDropdownOpen(false);
  };

  if (loading) {
    return (
      <div className="card">
        <div className="flex items-center space-x-3">
          <Building2 className="w-5 h-5 text-gray-400" />
          <div className="flex-1">
            <div className="h-4 bg-gray-200 rounded animate-pulse"></div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="card border-error-200 bg-error-50">
        <div className="flex items-center space-x-3">
          <Building2 className="w-5 h-5 text-error-500" />
          <p className="text-error-700">{error}</p>
          <button
            onClick={loadCompanies}
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
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-gray-900 flex items-center">
          <Building2 className="w-5 h-5 mr-2 text-primary-600" />
          Seleccionar Empresa
        </h2>
        <button
          onClick={onCreateCompany}
          className="btn-secondary text-sm flex items-center space-x-2"
        >
          <Plus className="w-4 h-4" />
          <span>Nueva Empresa</span>
        </button>
      </div>

      <div className="relative">
        <button
          onClick={() => setIsDropdownOpen(!isDropdownOpen)}
          className="w-full input flex items-center justify-between"
        >
          <span className={selectedCompany ? 'text-gray-900' : 'text-gray-500'}>
            {selectedCompany ? (
              <div className="flex flex-col items-start">
                <span className="font-medium">{selectedCompany.name}</span>
                <span className="text-sm text-gray-500">RUT: {selectedCompany.code}</span>
              </div>
            ) : (
              'Seleccionar empresa...'
            )}
          </span>
          <ChevronDown className={`w-4 h-4 text-gray-400 transition-transform ${isDropdownOpen ? 'rotate-180' : ''}`} />
        </button>

        {isDropdownOpen && (
          <div className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-auto">
            {companies.length === 0 ? (
              <div className="p-4 text-gray-500 text-center">
                <Building2 className="w-8 h-8 mx-auto mb-2 text-gray-300" />
                <p>No hay empresas registradas</p>
                <button
                  onClick={() => {
                    setIsDropdownOpen(false);
                    onCreateCompany();
                  }}
                  className="btn-primary text-sm mt-2"
                >
                  Crear primera empresa
                </button>
              </div>
            ) : (
              companies.map((company) => (
                <button
                  key={company.id}
                  onClick={() => handleCompanySelect(company)}
                  className="w-full p-3 text-left hover:bg-gray-50 focus:bg-gray-50 focus:outline-none border-b border-gray-100 last:border-b-0"
                >
                  <div className="flex flex-col">
                    <span className="font-medium text-gray-900">{company.name}</span>
                    <span className="text-sm text-gray-500">RUT: {company.code}</span>
                    {company.factura_movil_company_id && (
                      <span className="text-xs text-primary-600">
                        ID FM: {company.factura_movil_company_id}
                      </span>
                    )}
                  </div>
                </button>
              ))
            )}
          </div>
        )}
      </div>

      {selectedCompany && (
        <div className="mt-4 p-3 bg-primary-50 rounded-lg border border-primary-200">
          <h3 className="font-medium text-primary-900 mb-2">Empresa Seleccionada</h3>
          <div className="space-y-1 text-sm">
            <p><span className="font-medium">Nombre:</span> {selectedCompany.name}</p>
            <p><span className="font-medium">RUT:</span> {selectedCompany.code}</p>
            {selectedCompany.factura_movil_company_id && (
              <p><span className="font-medium">ID Factura Móvil:</span> {selectedCompany.factura_movil_company_id}</p>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default CompanySelector; 