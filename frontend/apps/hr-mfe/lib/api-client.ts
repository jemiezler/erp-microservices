/**
 * HR Service API Client
 * Uses native Fetch API (no external HTTP libraries)
 * 
 * Note: All requests go through the API Gateway at http://localhost:8080
 * The Gateway routes requests to backend services
 */

import { toast } from '@erp/ui';

const API_BASE_URL = process.env.NEXT_PUBLIC_HR_API_URL || 'http://localhost:8080/api/v1';

// Response types
export interface ApiResponse<T> {
  success: boolean;
  code: string;
  message: string;
  data: T;
  timestamp: string;
}

export interface PaginatedResponse<T> {
  success: boolean;
  code: string;
  data: T[];
  pagination: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}

// Employee Types
export interface Employee {
  id: number;
  employee_id: string;
  name: string;
  email: string;
  position: string;
  department: string;
  status: string;
  manager_id?: number;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface EmployeeFormData {
  employee_id: string;
  name: string;
  email: string;
  position: string;
  department: string;
  status: string;
  manager_id?: number;
  role: string;
}

// Leave Types
export interface LeaveBalance {
  leave_type: string;
  available: number;
  used: number;
  pending: number;
}

export interface Leave {
  id: number;
  employee_id: string;
  leave_type: string;
  start_date: string;
  end_date: string;
  status: string;
  reason: string;
  created_at: string;
}

// Attendance Types
export interface Attendance {
  id: number;
  employee_id: string;
  attendance_date: string;
  check_in_time?: string;
  check_out_time?: string;
  status: string;
  working_hours: number;
}

export interface AttendanceStats {
  present: number;
  absent: number;
  on_leave: number;
  total_employees: number;
  attendance_percentage: number;
}

// Payroll Types
export interface Payroll {
  id: number;
  employee_id: string;
  month: string;
  gross_salary: number;
  deductions: number;
  net_salary: number;
  status: string;
}

// API Error
export class ApiError extends Error {
  constructor(
    public status: number,
    public code: string,
    message: string
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * Generic fetch wrapper
 */
async function fetchAPI<T>(
  endpoint: string,
  options: Omit<RequestInit, 'headers'> & { headers?: Record<string, string> } = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string> || {}),
  };

  // Add tenant ID if available
  const tenantId = localStorage.getItem('tenant_id');
  if (tenantId) {
    headers['X-Tenant-ID'] = tenantId;
  }

  // Add auth token if available
  const authToken = localStorage.getItem('auth_token');
  if (authToken) {
    headers['Authorization'] = `Bearer ${authToken}`;
  }

  try {
    const response = await fetch(url, {
      ...options,
      headers,
    });

    const data = await response.json();

    if (!response.ok) {
      const apiError = new ApiError(
        response.status,
        data.code || 'ERROR',
        data.message || `HTTP ${response.status}`
      );

      if (typeof window !== 'undefined') {
        toast.error(apiError.message);
      }

      throw apiError;
    }

    return data;
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }

    const networkError = new ApiError(500, 'NETWORK_ERROR', 'Network request failed');

    if (typeof window !== 'undefined') {
      toast.error(networkError.message);
    }

    throw networkError;
  }
}

/**
 * Employee API
 */
export const employeeApi = {
  getAll: (page = 1, pageSize = 20) =>
    fetchAPI<PaginatedResponse<Employee>>(
      `/hr/employees?page=${page}&page_size=${pageSize}`
    ),

  getById: (id: number) =>
    fetchAPI<ApiResponse<Employee>>(`/hr/employees/${id}`),

  create: (data: EmployeeFormData) =>
    fetchAPI<ApiResponse<Employee>>('/hr/employees', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  update: (id: number, data: Partial<EmployeeFormData>) =>
    fetchAPI<ApiResponse<Employee>>(`/hr/employees/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    }),

  delete: (id: number) =>
    fetchAPI<ApiResponse<void>>(`/hr/employees/${id}`, {
      method: 'DELETE',
    }),

  getHierarchy: (id: number) =>
    fetchAPI<ApiResponse<Employee[]>>(`/hr/employees/${id}/hierarchy`),
};

/**
 * Leave API
 */
export const leaveApi = {
  getBalance: (employeeId: string) =>
    fetchAPI<ApiResponse<LeaveBalance[]>>(
      `/hr/leaves/balance?employee_id=${employeeId}`
    ),

  getLeaves: (employeeId: string, page = 1) =>
    fetchAPI<PaginatedResponse<Leave>>(
      `/hr/leaves?employee_id=${employeeId}&page=${page}`
    ),

  getPending: () =>
    fetchAPI<PaginatedResponse<Leave>>('/hr/leaves/pending'),

  apply: (data: {
    employee_id: string;
    leave_type: string;
    start_date: string;
    end_date: string;
    reason: string;
  }) =>
    fetchAPI<ApiResponse<Leave>>('/hr/leaves', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  approve: (leaveId: number, comment?: string) =>
    fetchAPI<ApiResponse<Leave>>(`/hr/leaves/${leaveId}/approve`, {
      method: 'PATCH',
      body: JSON.stringify({ comment }),
    }),

  reject: (leaveId: number, reason?: string) =>
    fetchAPI<ApiResponse<Leave>>(`/hr/leaves/${leaveId}/reject`, {
      method: 'PATCH',
      body: JSON.stringify({ reason }),
    }),

  getHolidays: () =>
    fetchAPI<ApiResponse<unknown[]>>('/hr/leaves/holidays'),
};

/**
 * Attendance API
 */
export const attendanceApi = {
  checkIn: (latitude: number, longitude: number, deviceId: string) =>
    fetchAPI<ApiResponse<Attendance>>('/hr/attendance/check-in', {
      method: 'POST',
      body: JSON.stringify({ latitude, longitude, device_id: deviceId }),
    }),

  checkOut: (latitude: number, longitude: number, deviceId: string) =>
    fetchAPI<ApiResponse<Attendance>>('/hr/attendance/check-out', {
      method: 'POST',
      body: JSON.stringify({ latitude, longitude, device_id: deviceId }),
    }),

  getTodayAttendance: () =>
    fetchAPI<ApiResponse<Attendance>>('/hr/attendance/today'),

  getHistory: (employeeId: string, startDate: string, endDate: string) =>
    fetchAPI<PaginatedResponse<Attendance>>(
      `/hr/attendance/history?employee_id=${employeeId}&start_date=${startDate}&end_date=${endDate}`
    ),

  getStats: (departmentId?: string) =>
    fetchAPI<ApiResponse<AttendanceStats>>(
      `/hr/attendance/stats${departmentId ? `?department_id=${departmentId}` : ''}`
    ),
};

/**
 * Payroll API
 */
export const payrollApi = {
  getByMonth: (month: string) =>
    fetchAPI<PaginatedResponse<Payroll>>(
      `/hr/payroll?month=${month}`
    ),

  getHistory: (employeeId: string) =>
    fetchAPI<PaginatedResponse<Payroll>>(
      `/hr/payroll/history?employee_id=${employeeId}`
    ),

  getPending: () =>
    fetchAPI<PaginatedResponse<Payroll>>('/hr/payroll/pending'),

  approve: (payrollId: number) =>
    fetchAPI<ApiResponse<Payroll>>(`/hr/payroll/${payrollId}/approve`, {
      method: 'PATCH',
    }),

  post: (payrollId: number) =>
    fetchAPI<ApiResponse<Payroll>>(`/hr/payroll/${payrollId}/post`, {
      method: 'PATCH',
    }),
};

const apiClient = {
  employeeApi,
  leaveApi,
  attendanceApi,
  payrollApi,
};

export default apiClient;
