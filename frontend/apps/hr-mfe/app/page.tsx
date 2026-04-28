'use client';

import { useState, useEffect } from 'react';
import { employeeApi, leaveApi, attendanceApi, payrollApi, Employee, Leave, Attendance, AttendanceStats, Payroll, ApiError } from '@/lib/api-client';

interface DashboardStats {
  totalEmployees: number;
  activeEmployees: number;
  pendingLeaves: number;
  todayPresent: number;
  pendingPayrolls: number;
}

export default function HRDashboard() {
  const [stats, setStats] = useState<DashboardStats>({
    totalEmployees: 0,
    activeEmployees: 0,
    pendingLeaves: 0,
    todayPresent: 0,
    pendingPayrolls: 0,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'overview' | 'employees' | 'leaves' | 'attendance' | 'payroll'>('overview');

  useEffect(() => {
    const loadDashboard = async () => {
      try {
        setLoading(true);
        setError(null);

        const [empRes, pendingLeavesRes, statsRes, payrollRes] = await Promise.all([
          employeeApi.getAll(1, 1),
          leaveApi.getPending(),
          attendanceApi.getStats(),
          payrollApi.getPending(),
        ]);

        setStats({
          totalEmployees: empRes.pagination?.total || 0,
          activeEmployees: empRes.data?.filter((e: any) => e.status === 'active').length || 0,
          pendingLeaves: pendingLeavesRes.pagination?.total || 0,
          todayPresent: (statsRes.data as any)?.present || 0,
          pendingPayrolls: payrollRes.pagination?.total || 0,
        });
      } catch (err) {
        const apiError = err as ApiError;
        setError(apiError.message || 'Failed to load dashboard');
        console.error('Dashboard error:', err);
      } finally {
        setLoading(false);
      }
    };

    loadDashboard();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-slate-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-slate-600">Loading dashboard...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-50">
      <header className="bg-white border-b border-slate-200 shadow-sm">
        <div className="max-w-7xl mx-auto px-6 py-6">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-3xl font-bold text-slate-900">HR Dashboard</h1>
              <p className="text-slate-600 mt-1">Manage employees, leaves, attendance & payroll</p>
            </div>
            <div className="bg-blue-100 text-blue-800 px-4 py-2 rounded-full text-sm font-semibold">
              HR-MFE v1.0
            </div>
          </div>
        </div>
      </header>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mx-6 mt-6">
          <p className="font-semibold">Error: {error}</p>
        </div>
      )}

      <div className="max-w-7xl mx-auto px-6 py-8">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4 mb-8">
          <div className="bg-white rounded-lg shadow p-6 border-l-4 border-blue-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-slate-600 text-sm font-medium">Total Employees</p>
                <p className="text-2xl font-bold text-slate-900 mt-2">{stats.totalEmployees}</p>
              </div>
              <div className="text-4xl text-blue-200">👥</div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6 border-l-4 border-green-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-slate-600 text-sm font-medium">Active</p>
                <p className="text-2xl font-bold text-slate-900 mt-2">{stats.activeEmployees}</p>
              </div>
              <div className="text-4xl text-green-200">✓</div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6 border-l-4 border-yellow-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-slate-600 text-sm font-medium">Pending Leaves</p>
                <p className="text-2xl font-bold text-slate-900 mt-2">{stats.pendingLeaves}</p>
              </div>
              <div className="text-4xl text-yellow-200">⏳</div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6 border-l-4 border-purple-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-slate-600 text-sm font-medium">Present Today</p>
                <p className="text-2xl font-bold text-slate-900 mt-2">{stats.todayPresent}</p>
              </div>
              <div className="text-4xl text-purple-200">📍</div>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6 border-l-4 border-orange-500">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-slate-600 text-sm font-medium">Pending Payroll</p>
                <p className="text-2xl font-bold text-slate-900 mt-2">{stats.pendingPayrolls}</p>
              </div>
              <div className="text-4xl text-orange-200">💰</div>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow mb-8">
          <div className="border-b border-slate-200 flex">
            {(['overview', 'employees', 'leaves', 'attendance', 'payroll'] as const).map(tab => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={`px-6 py-4 font-medium text-sm border-b-2 transition-colors ${
                  activeTab === tab
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-slate-600 hover:text-slate-900'
                }`}
              >
                {tab.charAt(0).toUpperCase() + tab.slice(1)}
              </button>
            ))}
          </div>

          <div className="p-6">
            {activeTab === 'overview' && (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="bg-blue-50 rounded-lg p-6 border border-blue-200">
                  <h3 className="font-semibold text-blue-900 mb-2">🎯 Quick Actions</h3>
                  <ul className="space-y-2 text-sm text-blue-800">
                    <li>• Add new employee</li>
                    <li>• Process leave requests</li>
                    <li>• Generate payroll</li>
                    <li>• View attendance reports</li>
                  </ul>
                </div>
                <div className="bg-green-50 rounded-lg p-6 border border-green-200">
                  <h3 className="font-semibold text-green-900 mb-2">📊 System Status</h3>
                  <ul className="space-y-2 text-sm text-green-800">
                    <li>✓ Database: Connected</li>
                    <li>✓ HR Service: Running</li>
                    <li>✓ All modules: Active</li>
                    <li>✓ No pending issues</li>
                  </ul>
                </div>
              </div>
            )}

            {activeTab === 'employees' && <EmployeesTab />}
            {activeTab === 'leaves' && <LeavesTab />}
            {activeTab === 'attendance' && <AttendanceTab />}
            {activeTab === 'payroll' && <PayrollTab />}
          </div>
        </div>
      </div>
    </div>
  );
}

function EmployeesTab() {
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadEmployees = async () => {
      try {
        const res = await employeeApi.getAll();
        setEmployees(res.data || []);
      } catch (err) {
        console.error('Failed to load employees:', err);
      } finally {
        setLoading(false);
      }
    };
    loadEmployees();
  }, []);

  if (loading) return <div className="text-center py-8">Loading employees...</div>;

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead className="bg-slate-100">
          <tr>
            <th className="px-4 py-3 text-left font-semibold text-slate-700">Employee ID</th>
            <th className="px-4 py-3 text-left font-semibold text-slate-700">Name</th>
            <th className="px-4 py-3 text-left font-semibold text-slate-700">Email</th>
            <th className="px-4 py-3 text-left font-semibold text-slate-700">Position</th>
            <th className="px-4 py-3 text-left font-semibold text-slate-700">Status</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-200">
          {employees.map((emp: any) => (
            <tr key={emp.id} className="hover:bg-slate-50">
              <td className="px-4 py-3 text-slate-600">{emp.employee_id}</td>
              <td className="px-4 py-3 font-medium text-slate-900">{emp.name}</td>
              <td className="px-4 py-3 text-slate-600">{emp.email}</td>
              <td className="px-4 py-3 text-slate-600">{emp.position}</td>
              <td className="px-4 py-3">
                <span className={`px-3 py-1 rounded-full text-xs font-semibold ${
                  emp.status === 'active'
                    ? 'bg-green-100 text-green-800'
                    : 'bg-slate-100 text-slate-800'
                }`}>
                  {emp.status}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function LeavesTab() {
  const [leaves, setLeaves] = useState<Leave[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadLeaves = async () => {
      try {
        const res = await leaveApi.getPending();
        setLeaves(res.data || []);
      } catch (err) {
        console.error('Failed to load leaves:', err);
      } finally {
        setLoading(false);
      }
    };
    loadLeaves();
  }, []);

  if (loading) return <div className="text-center py-8">Loading leave requests...</div>;

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead className="bg-slate-100">
          <tr>
            <th className="px-4 py-3 text-left font-semibold text-slate-700">Employee</th>
            <th className="px-4 py-3 text-left font-semibold text-slate-700">Leave Type</th>
            <th className="px-4 py-3 text-left font-semibold text-slate-700">Start Date</th>
            <th className="px-4 py-3 text-left font-semibold text-slate-700">End Date</th>
            <th className="px-4 py-3 text-left font-semibold text-slate-700">Status</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-200">
          {leaves.map((leave: any) => (
            <tr key={leave.id} className="hover:bg-slate-50">
              <td className="px-4 py-3 font-medium text-slate-900">{leave.employee_id}</td>
              <td className="px-4 py-3 text-slate-600">{leave.leave_type}</td>
              <td className="px-4 py-3 text-slate-600">{new Date(leave.start_date).toLocaleDateString()}</td>
              <td className="px-4 py-3 text-slate-600">{new Date(leave.end_date).toLocaleDateString()}</td>
              <td className="px-4 py-3">
                <span className={`px-3 py-1 rounded-full text-xs font-semibold ${
                  leave.status === 'approved'
                    ? 'bg-green-100 text-green-800'
                    : leave.status === 'pending'
                    ? 'bg-yellow-100 text-yellow-800'
                    : 'bg-red-100 text-red-800'
                }`}>
                  {leave.status}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function AttendanceTab() {
  const [stats, setStats] = useState<AttendanceStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadStats = async () => {
      try {
        const res = await attendanceApi.getStats();
        setStats(res.data as any);
      } catch (err) {
        console.error('Failed to load attendance stats:', err);
      } finally {
        setLoading(false);
      }
    };
    loadStats();
  }, []);

  if (loading) return <div className="text-center py-8">Loading attendance data...</div>;

  if (!stats) return <div className="text-center py-8 text-red-600">Failed to load attendance</div>;

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div className="bg-green-50 rounded-lg p-6 border border-green-200">
        <p className="text-slate-600 text-sm">Present</p>
        <p className="text-4xl font-bold text-green-600">{stats.present}</p>
      </div>
      <div className="bg-red-50 rounded-lg p-6 border border-red-200">
        <p className="text-slate-600 text-sm">Absent</p>
        <p className="text-4xl font-bold text-red-600">{stats.absent}</p>
      </div>
      <div className="bg-yellow-50 rounded-lg p-6 border border-yellow-200">
        <p className="text-slate-600 text-sm">On Leave</p>
        <p className="text-4xl font-bold text-yellow-600">{stats.on_leave}</p>
      </div>
      <div className="bg-blue-50 rounded-lg p-6 border border-blue-200">
        <p className="text-slate-600 text-sm">Attendance %</p>
        <p className="text-4xl font-bold text-blue-600">{stats.attendance_percentage.toFixed(1)}%</p>
      </div>
    </div>
  );
}

function PayrollTab() {
  const [payrolls, setPayrolls] = useState<Payroll[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadPayrolls = async () => {
      try {
        const res = await payrollApi.getPending();
        setPayrolls(res.data || []);
      } catch (err) {
        console.error('Failed to load payrolls:', err);
      } finally {
        setLoading(false);
      }
    };
    loadPayrolls();
  }, []);

  if (loading) return <div className="text-center py-8">Loading payroll data...</div>;

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead className="bg-slate-100">
          <tr>
            <th className="px-4 py-3 text-left font-semibold text-slate-700">Employee</th>
            <th className="px-4 py-3 text-left font-semibold text-slate-700">Month</th>
            <th className="px-4 py-3 text-right font-semibold text-slate-700">Gross</th>
            <th className="px-4 py-3 text-right font-semibold text-slate-700">Deductions</th>
            <th className="px-4 py-3 text-right font-semibold text-slate-700">Net</th>
            <th className="px-4 py-3 text-left font-semibold text-slate-700">Status</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-200">
          {payrolls.map((pr: any) => (
            <tr key={pr.id} className="hover:bg-slate-50">
              <td className="px-4 py-3 font-medium text-slate-900">{pr.employee_id}</td>
              <td className="px-4 py-3 text-slate-600">{pr.month}</td>
              <td className="px-4 py-3 text-right text-slate-900 font-semibold">${pr.gross_salary.toFixed(2)}</td>
              <td className="px-4 py-3 text-right text-red-600">-${pr.deductions.toFixed(2)}</td>
              <td className="px-4 py-3 text-right text-green-600 font-semibold">${pr.net_salary.toFixed(2)}</td>
              <td className="px-4 py-3">
                <span className={`px-3 py-1 rounded-full text-xs font-semibold ${
                  pr.status === 'approved'
                    ? 'bg-green-100 text-green-800'
                    : 'bg-yellow-100 text-yellow-800'
                }`}>
                  {pr.status}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
