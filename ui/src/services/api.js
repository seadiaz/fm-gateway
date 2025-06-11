import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || '';

// Mock data for development/demo purposes
const MOCK_COMPANIES = [
    {
        id: '1',
        name: 'Empresa Demo S.A.',
        code: '11111111-1',
        factura_movil_company_id: 12345
    },
    {
        id: '2',
        name: 'Tecnología Digital Ltda.',
        code: '22222222-2',
        factura_movil_company_id: 67890
    },
    {
        id: '3',
        name: 'Servicios Profesionales SpA',
        code: '33333333-3'
    }
];

const MOCK_CAFS = {
    '1': [
        {
            id: 'caf-1',
            companyCode: '11111111-1',
            documentType: 33,
            initialFolios: 1,
            finalFolios: 1000,
            currentFolios: 250,
            authorizationDate: '2024-01-01T00:00:00Z',
            expirationDate: '2024-12-31T23:59:59Z'
        },
        {
            id: 'caf-2',
            companyCode: '11111111-1',
            documentType: 39,
            initialFolios: 1,
            finalFolios: 500,
            currentFolios: 450,
            authorizationDate: '2024-01-01T00:00:00Z',
            expirationDate: '2024-06-30T23:59:59Z'
        }
    ],
    '2': [
        {
            id: 'caf-3',
            companyCode: '22222222-2',
            documentType: 33,
            initialFolios: 1001,
            finalFolios: 2000,
            currentFolios: 1100,
            authorizationDate: '2024-02-01T00:00:00Z',
            expirationDate: '2025-01-31T23:59:59Z'
        }
    ],
    '3': []
};

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
    timeout: 5000, // 5 second timeout
});

// Add request interceptor for logging
api.interceptors.request.use((config) => {
    console.log(`🔄 ${config.method?.toUpperCase()} ${config.url}`);
    return config;
});

// Add response interceptor for error handling
api.interceptors.response.use(
    (response) => {
        console.log(`✅ ${response.config.method?.toUpperCase()} ${response.config.url} - ${response.status}`);
        return response;
    },
    (error) => {
        console.error(`❌ ${error.config?.method?.toUpperCase()} ${error.config?.url} - ${error.response?.status}`);
        return Promise.reject(error);
    }
);

// Helper function to simulate API delay
const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

// Helper function to check if backend is available
const isBackendAvailable = async () => {
    try {
        await api.get('/healthz', { timeout: 2000 });
        return true;
    } catch (error) {
        console.warn('Backend not available, using mock data');
        return false;
    }
};

export const companyService = {
    // Get all companies
    getCompanies: async () => {
        try {
            const response = await api.get('/companies');
            return response.data;
        } catch (error) {
            console.warn('Using mock companies data');
            await delay(500); // Simulate network delay
            return MOCK_COMPANIES;
        }
    },

    // Get company by ID
    getCompany: async (companyId) => {
        try {
            const response = await api.get(`/companies/${companyId}`);
            return response.data;
        } catch (error) {
            console.warn('Using mock company data');
            await delay(300);
            const company = MOCK_COMPANIES.find(c => c.id === companyId);
            if (!company) {
                throw new Error('Company not found');
            }
            return company;
        }
    },

    // Create new company
    createCompany: async (companyData) => {
        try {
            const response = await api.post('/companies', companyData);
            return response.data;
        } catch (error) {
            console.warn('Using mock company creation');
            await delay(800);

            // Check if RUT already exists in mock data
            const existingCompany = MOCK_COMPANIES.find(c => c.code === companyData.code);
            if (existingCompany) {
                throw new Error('RUT already exists');
            }

            const newCompany = {
                id: String(MOCK_COMPANIES.length + 1),
                ...companyData
            };
            MOCK_COMPANIES.push(newCompany);
            return newCompany;
        }
    },
};

export const cafService = {
    // Upload CAF for a company
    uploadCAF: async (companyId, cafFile) => {
        try {
            const response = await api.post(`/companies/${companyId}/cafs`, cafFile, {
                headers: {
                    'Content-Type': 'application/xml',
                },
            });
            return response.data;
        } catch (error) {
            console.warn('Using mock CAF upload');
            await delay(1200);

            // Simulate CAF parsing
            const mockCAF = {
                id: `caf-${Date.now()}`,
                companyCode: MOCK_COMPANIES.find(c => c.id === companyId)?.code || '12345678-9',
                documentType: 33, // Factura Electrónica
                initialFolios: 1,
                finalFolios: 1000,
                currentFolios: 1,
                authorizationDate: new Date().toISOString(),
                expirationDate: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000).toISOString() // 1 year from now
            };

            // Add to mock data
            if (!MOCK_CAFS[companyId]) {
                MOCK_CAFS[companyId] = [];
            }
            MOCK_CAFS[companyId].push(mockCAF);

            return mockCAF;
        }
    },

    // Get CAFs for a company
    getCompanyCAFs: async (companyId) => {
        try {
            const response = await api.get(`/companies/${companyId}/cafs`);
            return response.data;
        } catch (error) {
            console.warn('Using mock CAFs data');
            await delay(400);
            return MOCK_CAFS[companyId] || [];
        }
    },
};

export const stampService = {
    // Generate stamp
    generateStamp: async (companyId, stampData, options = {}) => {
        const { format, includeBarcode } = options;
        let url = `/companies/${companyId}/stamps`;

        const params = new URLSearchParams();
        if (format) params.append('format', format);
        if (includeBarcode) params.append('include_barcode', 'true');

        if (params.toString()) {
            url += `?${params.toString()}`;
        }

        try {
            const response = await api.post(url, stampData, {
                responseType: format === 'pdf417' ? 'blob' : 'json',
            });
            return response.data;
        } catch (error) {
            console.warn('Stamp generation not available in demo mode');
            throw new Error('Stamp generation requires backend connection');
        }
    },
};

export default api; 