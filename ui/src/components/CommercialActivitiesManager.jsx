import { useState, useEffect } from 'react';
import { companyService } from '../services/api';

const CommercialActivitiesManager = ({ selectedCompany, refreshTrigger }) => {
    const [activities, setActivities] = useState([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);
    const [newActivity, setNewActivity] = useState({ code: '', description: '' });

    useEffect(() => {
        const loadActivities = async () => {
            if (!selectedCompany) return;
            
            setIsLoading(true);
            setError(null);
            try {
                const activitiesData = await companyService.getCommercialActivities(selectedCompany.id);
                setActivities(activitiesData);
            } catch (err) {
                setError('Error al cargar las actividades económicas');
                console.error('Error loading commercial activities:', err);
            } finally {
                setIsLoading(false);
            }
        };

        loadActivities();
    }, [selectedCompany, refreshTrigger]);

    const handleAddActivity = async () => {
        if (!newActivity.code || !newActivity.description) {
            setError('Por favor complete todos los campos');
            return;
        }

        setIsLoading(true);
        setError(null);
        try {
            await companyService.addCommercialActivity(selectedCompany.id, newActivity);
            setNewActivity({ code: '', description: '' });
            // Recargar actividades
            const activitiesData = await companyService.getCommercialActivities(selectedCompany.id);
            setActivities(activitiesData);
        } catch (err) {
            setError('Error al agregar la actividad económica');
            console.error('Error adding commercial activity:', err);
        } finally {
            setIsLoading(false);
        }
    };

    const handleRemoveActivity = async (activityId) => {
        if (!window.confirm('¿Está seguro de eliminar esta actividad económica?')) return;

        setIsLoading(true);
        setError(null);
        try {
            await companyService.removeCommercialActivity(selectedCompany.id, activityId);
            // Recargar actividades
            const activitiesData = await companyService.getCommercialActivities(selectedCompany.id);
            setActivities(activitiesData);
        } catch (err) {
            setError('Error al eliminar la actividad económica');
            console.error('Error removing commercial activity:', err);
        } finally {
            setIsLoading(false);
        }
    };

    if (!selectedCompany) return null;

    return (
        <div className="bg-white shadow rounded-lg p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Actividades Económicas</h2>
            
            {error && (
                <div className="mb-4 p-4 bg-red-50 text-red-700 rounded-md">
                    {error}
                </div>
            )}

            {/* Formulario para agregar actividad */}
            <div className="mb-6 p-4 bg-gray-50 rounded-lg">
                <h3 className="text-md font-medium text-gray-900 mb-3">Agregar Actividad Económica</h3>
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                    <div>
                        <label htmlFor="activityCode" className="block text-sm font-medium text-gray-700">
                            Código
                        </label>
                        <input
                            type="text"
                            id="activityCode"
                            value={newActivity.code}
                            onChange={(e) => setNewActivity({ ...newActivity, code: e.target.value })}
                            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
                            placeholder="Ej: 620100"
                        />
                    </div>
                    <div>
                        <label htmlFor="activityDescription" className="block text-sm font-medium text-gray-700">
                            Descripción
                        </label>
                        <input
                            type="text"
                            id="activityDescription"
                            value={newActivity.description}
                            onChange={(e) => setNewActivity({ ...newActivity, description: e.target.value })}
                            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
                            placeholder="Ej: Servicios de consultores en informática"
                        />
                    </div>
                </div>
                <div className="mt-4">
                    <button
                        onClick={handleAddActivity}
                        disabled={isLoading}
                        className="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
                    >
                        {isLoading ? 'Agregando...' : 'Agregar Actividad'}
                    </button>
                </div>
            </div>

            {/* Lista de actividades */}
            <div className="mt-6">
                <h3 className="text-md font-medium text-gray-900 mb-3">Actividades Registradas</h3>
                {isLoading ? (
                    <div className="text-center py-4">
                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
                    </div>
                ) : activities.length === 0 ? (
                    <p className="text-gray-500 text-center py-4">No hay actividades económicas registradas</p>
                ) : (
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-gray-200">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Código
                                    </th>
                                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Descripción
                                    </th>
                                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                        Acciones
                                    </th>
                                </tr>
                            </thead>
                            <tbody className="bg-white divide-y divide-gray-200">
                                {activities.map((activity) => (
                                    <tr key={activity.id}>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                            {activity.code}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {activity.description}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                            <button
                                                onClick={() => handleRemoveActivity(activity.id)}
                                                className="text-red-600 hover:text-red-900"
                                            >
                                                Eliminar
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </div>
    );
};

export default CommercialActivitiesManager; 