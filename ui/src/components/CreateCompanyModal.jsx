import React, { useState } from 'react';
import { X, Building2, Save } from 'lucide-react';
import { companyService } from '../services/api';

const CreateCompanyModal = ({ isOpen, onClose, onCompanyCreated }) => {
  const [formData, setFormData] = useState({
    name: '',
    code: '',
    factura_movil_company_id: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const validateRUT = (rut) => {
    // Basic RUT validation for Chilean format
    const rutRegex = /^[0-9]+-[0-9kK]$/;
    return rutRegex.test(rut);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.name.trim()) {
      setError('El nombre de la empresa es requerido');
      return;
    }

    if (!formData.code.trim()) {
      setError('El RUT de la empresa es requerido');
      return;
    }

    if (!validateRUT(formData.code)) {
      setError('El RUT debe tener el formato correcto (ej: 12345678-9)');
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const companyData = {
        name: formData.name.trim(),
        code: formData.code.trim(),
        factura_movil_company_id: formData.factura_movil_company_id ? 
          parseInt(formData.factura_movil_company_id) : undefined
      };

      const newCompany = await companyService.createCompany(companyData);
      
      // Reset form
      setFormData({
        name: '',
        code: '',
        factura_movil_company_id: ''
      });

      // Notify parent and close modal
      if (onCompanyCreated) {
        onCompanyCreated(newCompany);
      }
      onClose();

    } catch (err) {
      console.error('Error creating company:', err);
      setError(
        err.response?.data?.error || 
        'Error al crear la empresa. Verifica que el RUT no esté ya registrado.'
      );
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setFormData({
      name: '',
      code: '',
      factura_movil_company_id: ''
    });
    setError(null);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <h2 className="text-lg font-semibold text-gray-900 flex items-center">
            <Building2 className="w-5 h-5 mr-2 text-primary-600" />
            Crear Nueva Empresa
          </h2>
          <button
            onClick={handleClose}
            className="p-2 hover:bg-gray-100 rounded-lg"
          >
            <X className="w-5 h-5 text-gray-400" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-6">
          <div>
            <label htmlFor="name" className="label">
              Nombre de la Empresa *
            </label>
            <input
              type="text"
              id="name"
              name="name"
              value={formData.name}
              onChange={handleChange}
              className="input"
              placeholder="Ej: Mi Empresa S.A."
              required
            />
          </div>

          <div>
            <label htmlFor="code" className="label">
              RUT de la Empresa *
            </label>
            <input
              type="text"
              id="code"
              name="code"
              value={formData.code}
              onChange={handleChange}
              className="input"
              placeholder="Ej: 12345678-9"
              pattern="[0-9]+-[0-9kK]"
              required
            />
            <p className="text-sm text-gray-500 mt-1">
              Formato: 12345678-9 (incluye el guión y dígito verificador)
            </p>
          </div>

          <div>
            <label htmlFor="factura_movil_company_id" className="label">
              ID Factura Móvil (Opcional)
            </label>
            <input
              type="number"
              id="factura_movil_company_id"
              name="factura_movil_company_id"
              value={formData.factura_movil_company_id}
              onChange={handleChange}
              className="input"
              placeholder="Ej: 12345"
            />
            <p className="text-sm text-gray-500 mt-1">
              ID asignado por el sistema Factura Móvil (si aplica)
            </p>
          </div>

          {error && (
            <div className="p-4 bg-error-50 border border-error-200 rounded-lg">
              <p className="text-error-700 text-sm">{error}</p>
            </div>
          )}

          <div className="flex justify-end space-x-3 pt-6 border-t border-gray-200">
            <button
              type="button"
              onClick={handleClose}
              className="btn-secondary"
              disabled={loading}
            >
              Cancelar
            </button>
            <button
              type="submit"
              className={`btn-primary ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
              disabled={loading}
            >
              {loading ? (
                <div className="flex items-center space-x-2">
                  <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                  <span>Creando...</span>
                </div>
              ) : (
                <div className="flex items-center space-x-2">
                  <Save className="w-4 h-4" />
                  <span>Crear Empresa</span>
                </div>
              )}
            </button>
          </div>
        </form>

        {/* Information section */}
        <div className="px-6 pb-6">
          <div className="p-4 bg-blue-50 border border-blue-200 rounded-lg">
            <h3 className="font-medium text-blue-900 mb-2">Información Importante</h3>
            <ul className="text-sm text-blue-700 space-y-1">
              <li>• El RUT debe ser válido y único en el sistema</li>
              <li>• El ID Factura Móvil es opcional y puede agregarse después</li>
              <li>• Una vez creada, podrás cargar archivos CAF para esta empresa</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CreateCompanyModal; 