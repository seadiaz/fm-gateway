import React, { useState, useRef } from 'react';
import { Upload, File, CheckCircle, AlertCircle, X, FileText } from 'lucide-react';
import { cafService } from '../services/api';

const CAFUploader = ({ selectedCompany, onUploadSuccess }) => {
  const [dragActive, setDragActive] = useState(false);
  const [selectedFile, setSelectedFile] = useState(null);
  const [uploading, setUploading] = useState(false);
  const [uploadResult, setUploadResult] = useState(null);
  const [error, setError] = useState(null);
  const fileInputRef = useRef(null);

  const handleDrag = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  };

  const handleDrop = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleFile(e.dataTransfer.files[0]);
    }
  };

  const handleFileSelect = (e) => {
    if (e.target.files && e.target.files[0]) {
      handleFile(e.target.files[0]);
    }
  };

  const handleFile = (file) => {
    // Validate file type
    if (!file.name.toLowerCase().endsWith('.xml')) {
      setError('Por favor selecciona un archivo XML válido');
      return;
    }

    // Validate file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      setError('El archivo es demasiado grande. Máximo 5MB permitido');
      return;
    }

    setSelectedFile(file);
    setError(null);
    setUploadResult(null);
  };

  const uploadFile = async () => {
    if (!selectedFile || !selectedCompany) {
      setError('Selecciona un archivo y una empresa');
      return;
    }

    try {
      setUploading(true);
      setError(null);

      const result = await cafService.uploadCAF(selectedCompany.id, selectedFile);
      
      setUploadResult({
        type: 'success',
        message: 'CAF cargado exitosamente',
        data: result
      });

      // Reset form
      setSelectedFile(null);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }

      // Notify parent component
      if (onUploadSuccess) {
        onUploadSuccess(result);
      }

    } catch (err) {
      console.error('Upload error:', err);
      setError(
        err.response?.data?.error || 
        'Error al cargar el archivo CAF. Verifica que el archivo sea válido.'
      );
    } finally {
      setUploading(false);
    }
  };

  const clearSelection = () => {
    setSelectedFile(null);
    setError(null);
    setUploadResult(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  if (!selectedCompany) {
    return (
      <div className="card border-gray-200 bg-gray-50">
        <div className="text-center py-8">
          <Upload className="w-12 h-12 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500">Selecciona una empresa para cargar archivos CAF</p>
        </div>
      </div>
    );
  }

  return (
    <div className="card">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-lg font-semibold text-gray-900 flex items-center">
          <FileText className="w-5 h-5 mr-2 text-primary-600" />
          Cargar Archivo CAF
        </h2>
        <div className="text-sm text-gray-500">
          Empresa: <span className="font-medium text-gray-900">{selectedCompany.name}</span>
        </div>
      </div>

      {/* Upload Area */}
      <div
        className={`relative border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
          dragActive
            ? 'border-primary-400 bg-primary-50'
            : selectedFile
            ? 'border-success-300 bg-success-50'
            : 'border-gray-300 hover:border-gray-400'
        }`}
        onDragEnter={handleDrag}
        onDragLeave={handleDrag}
        onDragOver={handleDrag}
        onDrop={handleDrop}
      >
        <input
          ref={fileInputRef}
          type="file"
          accept=".xml"
          onChange={handleFileSelect}
          className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
        />

        {selectedFile ? (
          <div className="space-y-4">
            <CheckCircle className="w-12 h-12 text-success-500 mx-auto" />
            <div>
              <p className="text-lg font-medium text-success-700">Archivo seleccionado</p>
              <div className="mt-2 p-3 bg-white rounded border border-success-200">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    <File className="w-5 h-5 text-success-600" />
                    <div>
                      <p className="font-medium text-gray-900">{selectedFile.name}</p>
                      <p className="text-sm text-gray-500">
                        {(selectedFile.size / 1024).toFixed(1)} KB
                      </p>
                    </div>
                  </div>
                  <button
                    onClick={clearSelection}
                    className="p-1 hover:bg-gray-100 rounded"
                  >
                    <X className="w-4 h-4 text-gray-400" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <Upload className={`w-12 h-12 mx-auto ${dragActive ? 'text-primary-500' : 'text-gray-400'}`} />
            <div>
              <p className="text-lg font-medium text-gray-900">
                Arrastra tu archivo CAF aquí
              </p>
              <p className="text-gray-500 mt-1">
                o <span className="text-primary-600 font-medium">haz clic para seleccionar</span>
              </p>
              <p className="text-sm text-gray-400 mt-2">
                Archivos XML únicamente, máximo 5MB
              </p>
            </div>
          </div>
        )}
      </div>

      {/* Error Message */}
      {error && (
        <div className="mt-4 p-4 bg-error-50 border border-error-200 rounded-lg">
          <div className="flex items-center space-x-3">
            <AlertCircle className="w-5 h-5 text-error-500 flex-shrink-0" />
            <p className="text-error-700">{error}</p>
          </div>
        </div>
      )}

      {/* Success Message */}
      {uploadResult?.type === 'success' && (
        <div className="mt-4 p-4 bg-success-50 border border-success-200 rounded-lg">
          <div className="flex items-center space-x-3">
            <CheckCircle className="w-5 h-5 text-success-500 flex-shrink-0" />
            <div>
              <p className="text-success-700 font-medium">{uploadResult.message}</p>
              {uploadResult.data && (
                <div className="mt-2 text-sm text-success-600">
                  <p>Tipo de documento: {uploadResult.data.documentType}</p>
                  <p>Folios: {uploadResult.data.initialFolios} - {uploadResult.data.finalFolios}</p>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Upload Button */}
      <div className="mt-6 flex justify-end space-x-3">
        {selectedFile && (
          <button
            onClick={clearSelection}
            className="btn-secondary"
            disabled={uploading}
          >
            Cancelar
          </button>
        )}
        <button
          onClick={uploadFile}
          disabled={!selectedFile || uploading}
          className={`btn-primary ${uploading ? 'opacity-50 cursor-not-allowed' : ''}`}
        >
          {uploading ? (
            <div className="flex items-center space-x-2">
              <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
              <span>Cargando...</span>
            </div>
          ) : (
            'Cargar CAF'
          )}
        </button>
      </div>

      {/* CAF Information */}
      <div className="mt-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
        <h3 className="font-medium text-blue-900 mb-2 flex items-center">
          <FileText className="w-4 h-4 mr-2" />
          ¿Qué es un archivo CAF?
        </h3>
        <p className="text-sm text-blue-700">
          El CAF (Código de Autorización de Folios) es un archivo XML proporcionado por el SII 
          que autoriza a tu empresa a emitir documentos tributarios electrónicos con folios específicos.
        </p>
      </div>
    </div>
  );
};

export default CAFUploader; 